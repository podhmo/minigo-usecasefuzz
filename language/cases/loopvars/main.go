package main

import "fmt"

func main() {
	// Go 1.22: loop var is per-iteration — closures capture distinct cells
	fs := make([]func() int, 0, 3)
	for i := 0; i < 3; i++ {
		fs = append(fs, func() int { return i })
	}
	for _, f := range fs {
		fmt.Print(f(), " ")
	}
	fmt.Println()
	// range var same rule
	gs := make([]func() string, 0, 3)
	for _, v := range []string{"a", "b", "c"} {
		v := v // explicit copy also allowed
		gs = append(gs, func() string { return v })
	}
	fmt.Println(gs[0](), gs[1](), gs[2]())
}
