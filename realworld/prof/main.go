// Command prof runs one realworld task under the minigo engine with
// runtime/pprof enabled (the minigo CLI has no profiling flags).
// run.sh builds it against $MINIGO_DIR via a generated go.mod.
//
//	prof -dir <task dir> -out <prefix>
//
// writes <prefix>.cpu.pprof and <prefix>.allocs.pprof, plus
// <prefix>.trace when -trace is set.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"time"

	"github.com/podhmo/minigo"
)

func main() {
	dir := flag.String("dir", ".", "task dir (engine start dir)")
	out := flag.String("out", "prof", "output path prefix")
	withTrace := flag.Bool("trace", false, "also write an execution trace")
	flag.Parse()

	cpu, err := os.Create(*out + ".cpu.pprof")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cpu.Close()
	if *withTrace {
		tf, err := os.Create(*out + ".trace")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer tf.Close()
		if err := trace.Start(tf); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	if err := pprof.StartCPUProfile(cpu); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	t0 := time.Now()
	e := minigo.NewEngine(*dir, minigo.WithOutput(io.Discard))
	_, runErr := e.Run(context.Background(), *dir, "main")
	elapsed := time.Since(t0)

	pprof.StopCPUProfile()
	if *withTrace {
		trace.Stop()
	}
	if f, err := os.Create(*out + ".allocs.pprof"); err == nil {
		pprof.Lookup("allocs").WriteTo(f, 0)
		f.Close()
	}
	status := "ok"
	if runErr != nil {
		fmt.Fprintln(os.Stderr, runErr)
		status = "error (profile covers the run up to it)"
	}
	fmt.Printf("elapsed=%.2fs %s\n", elapsed.Seconds(), status)
}
