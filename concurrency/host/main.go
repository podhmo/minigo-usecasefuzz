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

	// T4: a spawned script goroutine parked in a channel op dies with the
	// process — its deferred parkprobe.Gone reports the unwind, so the
	// check needs no goroutine counting (a net NumGoroutine delta is
	// maskable by unrelated goroutine churn).
	{
		e := minigo.NewEngine(".")
		gone := make(chan struct{})
		e.Bind("parkprobe", map[string]runtime.Value{
			"Gone": &runtime.BuiltinFunc{
				Name: "parkprobe.Gone",
				Fn: func(_ runtime.VMCaller, _ []runtime.Value) (runtime.Value, error) {
					close(gone)
					return nil, nil
				},
			},
		})
		before := goruntime.NumGoroutine()
		_, err := e.Run(ctx, "./testdata/hostpark", "DetachedLeakProbe")
		var fired bool
		select {
		case <-gone:
			fired = true
		case <-time.After(30 * time.Second): // anti-hang bound, not a timing check
		}
		fmt.Printf("T4 detach-dies err=%v gone=%v (want true) goroutines before=%d after=%d\n",
			err, fired, before, goruntime.NumGoroutine())
	}

	// T5: a spawned goroutine parked in a HOST call (WaitGroup.Wait, not a
	// select) is not released by proc kill — it leaks. parkprobe.Parked
	// reports the goroutine reached the park point; the reflect method call
	// into inner.Wait() then provably parks it in a real WaitGroup.
	{
		e := minigo.NewEngine(".")
		parked := make(chan struct{})
		e.Bind("parkprobe", map[string]runtime.Value{
			"Parked": &runtime.BuiltinFunc{
				Name: "parkprobe.Parked",
				Fn: func(_ runtime.VMCaller, _ []runtime.Value) (runtime.Value, error) {
					close(parked)
					return nil, nil
				},
			},
		})
		before := goruntime.NumGoroutine()
		_, err := e.Run(ctx, "./testdata/hostpark", "DetachedWaitProbe")
		var fired bool
		select {
		case <-parked:
			fired = true
		case <-time.After(30 * time.Second): // anti-hang bound, not a timing check
		}
		fmt.Printf("T5 wait-leak err=%v parked=%v (want true) goroutines before=%d after=%d\n",
			err, fired, before, goruntime.NumGoroutine())
	}
}
