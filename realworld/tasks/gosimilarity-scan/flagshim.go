package gosimilarity

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

// Local stand-in for github.com/Songmu/flagx so the package runs under
// minigo (flagx calls runtime.AddCleanup, which minigo does not bind).

func flagUsage(fs *flag.FlagSet, w io.Writer, usage, desc string) func() {
	return func() {
		fmt.Fprintln(w, usage)
		if desc != "" {
			fmt.Fprintln(w)
			fmt.Fprintln(w, desc)
		}
		fs.PrintDefaults()
	}
}

func eachName(fs *flag.FlagSet, names string, fn func(name string)) {
	for _, n := range strings.Split(names, "|") {
		fn(n)
	}
}

func float64Var(fs *flag.FlagSet, p *float64, names string, v float64, u string) {
	eachName(fs, names, func(n string) { fs.Float64Var(p, n, v, u) })
}

func intVar(fs *flag.FlagSet, p *int, names string, v int, u string) {
	eachName(fs, names, func(n string) { fs.IntVar(p, n, v, u) })
}

func boolVar(fs *flag.FlagSet, p *bool, names string, v bool, u string) {
	eachName(fs, names, func(n string) { fs.BoolVar(p, n, v, u) })
}

func stringVar(fs *flag.FlagSet, p *string, names string, v, u string) {
	eachName(fs, names, func(n string) { fs.StringVar(p, n, v, u) })
}
