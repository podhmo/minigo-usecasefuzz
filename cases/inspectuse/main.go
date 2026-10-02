package main

import (
	"fmt"
	"sort"

	"minigo.dev/inspect"
)

// minigo-only: script introspection over its own package — the
// "enumerate what this module provides" job. No go-run oracle:
// minigo.dev/inspect does not exist for the toolchain, so this case
// is verified for sensible output rather than diffed.

// Greeting is a documented func.
func Greeting(name string) string { return "hi " + name }

type Pair struct{ A, B int }

var Table = map[string]int{"x": 1}

func main() {
	p := inspect.Current()
	fmt.Println("name:", inspect.Name(p))
	fmt.Println("path:", inspect.Path(p))
	fmt.Println("standard:", inspect.Standard(p))

	ds := inspect.Decls(p)
	names := []string{}
	for _, d := range ds {
		names = append(names, inspect.Kind(d)+":"+d.Name)
	}
	sort.Strings(names)
	for _, n := range names {
		fmt.Println("decl", n)
	}

	d := inspect.Symbol(p, "Greeting")
	fmt.Println("doc:", inspect.Doc(d))
	sig := inspect.Signature(d)
	fmt.Println("params:", len(sig.Params), "results:", len(sig.Results))
}
