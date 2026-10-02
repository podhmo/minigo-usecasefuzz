package main

import "fmt"

var calls []string

func track(name string, v int) int {
	calls = append(calls, name)
	return v
}

var a = track("a", 1) + b // a depends on b → b first
var b = track("b", 2)
var c = track("c", 3)

func main() {
	fmt.Printf("order: %v\n", calls)
	fmt.Printf("vals: %d %d %d\n", a, b, c)
}
