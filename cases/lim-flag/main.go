package main

import (
	"flag"
	"fmt"
)

// Standard CLI flag parsing.
func main() {
	name := flag.String("name", "world", "who to greet")
	verbose := flag.Bool("v", false, "verbose")
	flag.Parse()
	fmt.Println("hello", *name, "verbose:", *verbose, "rest:", flag.Args())
}
