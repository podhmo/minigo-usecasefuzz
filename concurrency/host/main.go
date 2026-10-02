package main

import (
	"context"
	"fmt"
	goruntime "runtime"
	"time"

	minigo "github.com/podhmo/minigo"
	"github.com/podhmo/minigo/runtime"
)

func main() {
	ctx := context.Background()
	// T1: init panic -> package must keep failing on subsequent access
	{
		e := minigo.NewEngine(".")
		v1, err1 := e.Run(ctx, "./testdata/initfail", "Use")
		fmt.Printf("T1a first-run  v=%v err=%v\n", v1, err1)
		v2, err2 := e.Run(ctx, "./testdata/initfail", "Use")
		fmt.Printf("T1b second-run v=%v err=%v (want err non-nil)\n", v2, err2)
		v3, err3 := e.Run(ctx, "./testdata/initfail", "Use")
		fmt.Printf("T1c third-run  v=%v err=%v\n", v3, err3)
	}

	// T2: ImportRef.Materialize panic -> later calls must not return (nil,nil)
	{
		r := &runtime.ImportRef{Path: "x", Load: func(string) (*runtime.Package, error) {
			panic("boom")
		}}
		func() {
			defer func() { fmt.Printf("T2a recovered=%v\n", recover() != nil) }()
			r.Materialize()
		}()
		p, err := r.Materialize()
		fmt.Printf("T2b second Materialize: pkg-nil=%v err=%v (want err non-nil)\n", p == nil, err)
	}

	// T3: two concurrent e.Run on the same engine (shared e.vmm)
	{
		e := minigo.NewEngine(".")
		done := make(chan string, 2)
		for i := 0; i < 2; i++ {
			go func() {
				_, err := e.Run(ctx, "./testdata/concurrency", "MutexCounter")
				done <- fmt.Sprint(err)
			}()
		}
		fmt.Printf("T3 concurrent Run errs: %q %q\n", <-done, <-done)
	}

	// T4: goroutine leak — a spawned script goroutine parked after Call returns
	{
		before := goruntime.NumGoroutine()
		e := minigo.NewEngine(".")
		_, err := e.Run(ctx, "./testdata/concurrency", "DetachedLeak")
		time.Sleep(50 * time.Millisecond)
		fmt.Printf("T4 leak err=%v goroutines before=%d after=%d\n", err, before, goruntime.NumGoroutine())
	}

	// T5: a spawned goroutine parked in a HOST call (WaitGroup.Wait, not a
	// select) is not released by proc kill — it leaks.
	{
		before := goruntime.NumGoroutine()
		e := minigo.NewEngine(".")
		_, err := e.Run(ctx, "./testdata/concurrency", "DetachedWait")
		time.Sleep(50 * time.Millisecond)
		fmt.Printf("T5 wait-leak err=%v goroutines before=%d after=%d\n", err, before, goruntime.NumGoroutine())
	}
}
