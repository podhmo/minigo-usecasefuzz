package main

import "fmt"

func try(name string, f func() any) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s: panic %v\n", name, r)
		}
	}()
	fmt.Printf("%s: %v\n", name, f())
}

type P struct{ X int }

func main() {
	try("nilFieldWrite", func() any { var p *P; p.X = 1; return 0 })
	try("reassign", func() any {
		x, y := 1, 2
		p := &x
		p = &y
		return *p
	})
}
