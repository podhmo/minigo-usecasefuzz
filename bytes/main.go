// Command byteprobe measures what it costs minigo to move []byte
// across the script<->host boundary and through slice conversions.
//
// Every op runs inside a real Engine through the public API, so the
// numbers include marshaling exactly where script code triggers it.
// Sizes stay small (KiB..1MiB, a handful of reps): the point is the
// per-byte shape of each path, not throughput under load.
//
// Usage:
//
//	go run . [-n 65536] [-reps 5]
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	goruntime "runtime"
	"testing"
	"time"

	"github.com/podhmo/minigo"
	"github.com/podhmo/minigo/runtime"
)

// devNull is a host io.Writer bound as a GoValue: `w.Write(buf)` in a
// script dispatches through reflection (callReflectFunc), exercising
// the script->host element-wise marshal.
type devNull struct{}

func (devNull) Write(p []byte) (int, error) { return len(p), nil }

// patternReader fills every buffer with 0xAB: its writes must land on
// the script's slice for Read callers to see them.
type patternReader struct{}

func (patternReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0xAB
	}
	return len(p), nil
}

type probeDef struct {
	name string
	fn   string // script function to call
	arg  bool   // whether to pass n
}

var probes = []probeDef{
	{"Noop (call overhead)", "Noop", false},
	{"make([]byte,n)", "MakeBuf", true},
	{"index loop b[i]", "PIndex", false},
	{"range loop", "PRange", false},
	{"copy(dst,buf)", "PCopy", false},
	{"append x1", "PAppend", true},
	{"string(buf)", "PToString", false},
	{"[]byte(str)", "PFromString", false},
	{"fmt %x", "PFmtHex", false},
	{"bytes.NewReader", "PNewReader", false},
	{"r.Read(buf)", "PRead", false},
	{"w.Write(buf)", "PWrite", false},
	{"io.ReadAll(n)", "PReadAll", true},
	{"io.ReadFull", "PReadFull", false},
	{"io.Copy(Discard,r)", "PIoCopy", false},
}

func main() {
	var (
		sizes = flag.String("sizes", "4096,65536,1048576", "comma-separated byte sizes")
		reps  = flag.Int("reps", 5, "timing repetitions per probe")
	)
	flag.Parse()

	dir, err := filepath.Abs("script")
	if err != nil {
		fatal(err)
	}
	ctx := context.Background()

	e := minigo.NewEngine(dir, minigo.WithOutput(os.Stderr))
	e.Bind("probe", map[string]runtime.Value{
		"DevNull":    &runtime.GoValue{V: devNull{}},
		"FillReader": &runtime.GoValue{V: patternReader{}},
		"NewReader": minigo.WrapFunc("probe.NewReader", func(n int) any {
			return &runtime.GoValue{V: bytes.NewReader(make([]byte, n))}
		}),
	})
	pkg, err := e.LoadFile(ctx, filepath.Join(dir, "probe.go"))
	if err != nil {
		fatal(err)
	}
	call := func(name string, args ...runtime.Value) {
		if _, err := e.Call(ctx, pkg, name, args...); err != nil {
			fatalf("%s: %v", name, err)
		}
	}

	// functional check first: does io.ReadFull write into the script
	// slice at all? (byteSlice() copies — expect a missing write-back.)
	call("Setup", int64(64))
	if r, err := e.Call(ctx, pkg, "PReadCheck"); err != nil {
		fatalf("PReadCheck: %v", err)
	} else {
		fmt.Printf("# check: buf[0] after io.ReadFull = %v (want 171)\n\n", runtime.Unwrap(r))
	}

	fmt.Printf("%-22s %10s %12s %12s %10s %12s\n", "op", "n", "ns/op", "ns/byte", "allocs/op", "allocs/byte")
	var ns []int
	for _, s := range splitSizes(*sizes) {
		ns = append(ns, s)
	}
	for _, n := range ns {
		call("Setup", int64(n))
		for _, p := range probes {
			var args []runtime.Value
			if p.arg {
				args = []runtime.Value{int64(n)}
			}
			// warm up once so lazy materialization isn't measured
			call(p.fn, args...)
			allocs := testing.AllocsPerRun(3, func() { call(p.fn, args...) })
			// timing: min over reps of a single call
			var best time.Duration
			for r := 0; r < *reps; r++ {
				start := time.Now()
				call(p.fn, args...)
				if d := time.Since(start); best == 0 || d < best {
					best = d
				}
			}
			fmt.Printf("%-22s %10d %12d %12.3f %10.0f %12.4f\n",
				p.name, n, best.Nanoseconds(), float64(best.Nanoseconds())/float64(n),
				allocs, allocs/float64(n))
		}
		fmt.Println()
	}
	goruntime.KeepAlive(e)
}

func splitSizes(s string) []int {
	var out []int
	cur := 0
	for _, c := range s {
		if c == ',' {
			out = append(out, cur)
			cur = 0
			continue
		}
		cur = cur*10 + int(c-'0')
	}
	return append(out, cur)
}

func fatal(err error) { fatalf("%v", err) }
func fatalf(f string, a ...any) {
	fmt.Fprintf(os.Stderr, "byteprobe: "+f+"\n", a...)
	os.Exit(1)
}
