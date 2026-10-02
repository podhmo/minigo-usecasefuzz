package main

import "sync"

func main() {
	var o sync.Once
	n := 0
	f := func() { n++ }
	o.Do(f)
	o.Do(f)
	println("n =", n)
}
