package main

import (
	"fmt"
)

// fmt verb coverage over script values — the "debug-print a data
// structure" job.
type Point struct {
	X, Y int
}

type User struct {
	Name string
	Age  int
	Tags []string
}

func main() {
	u := User{Name: "ann", Age: 30, Tags: []string{"a", "b"}}
	p := Point{X: 1, Y: 2}

	fmt.Printf("v:   %v\n", u)
	fmt.Printf("+v:  %+v\n", u)
	fmt.Printf("#v:  %#v\n", u)
	fmt.Printf("T:   %T\n", u)
	fmt.Printf("ptr: %v\n", &p)

	m := map[string]int{"b": 2, "a": 1}
	fmt.Printf("map v:  %v\n", m)
	fmt.Printf("map +v: %+v\n", m)

	s := []int{1, 2, 3}
	fmt.Printf("slice: %v / %#v\n", s, s)

	fmt.Printf("q: %q %q\n", "hi\nthere", 'x')
	fmt.Printf("x: %x %X %04x %08b\n", 255, 255, 255, 5)
	fmt.Printf("f: %.2f %8.2f %g\n", 3.14159, 3.14159, 123456789.0)
	fmt.Printf("e: %e\n", 0.0001)
	fmt.Printf("pct: %.1f%%\n", 95.5)
	fmt.Printf("pad: |%6d|%-6d|%06d|\n", 42, 42, 42)
}
