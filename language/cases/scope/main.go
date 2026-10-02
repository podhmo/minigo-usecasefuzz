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

var gx = 10

func main() {
	try("shadowBlock", func() any {
		x := 1
		{
			x := 2 // new var
			x++
			_ = x
		}
		return x // 1
	})
	try("shadowAssign", func() any {
		x := 1
		{
			x = 2 // same var
		}
		return x // 2
	})
	try("ifInitScope", func() any {
		x := 1
		if x := 5; x > 0 {
			return x
		}
		return x
	})
	try("forInitScope", func() any {
		x := 99
		for x := 0; x < 2; x++ {
		}
		return x // 99
	})
	try("switchInitScope", func() any {
		v := 0
		switch x := 5; x {
		case 5:
			v = x
		}
		return v
	})
	try("redeclare", func() any {
		x := 1
		x, y := 2, 3 // redeclare with one new var — legal
		return x + y
	})
	try("blankIdent", func() any {
		_ = 5
		_, b := 1, 2
		return b
	})
	try("globalRef", func() any { return gx })
	try("globalWrite", func() any { gx = 11; return gx })   // mutates package var
	try("globalShadow", func() any { gx := 99; return gx }) // local shadows
	try("closureCapLocal", func() any {
		x := 1
		f := func() { x = 2 }
		f()
		return x
	})
	try("namedResultScope", func() any {
		f := func() (n int) {
			n = 3
			{
				n := 9 // shadows result param — must not touch it
				_ = n
			}
			return
		}
		return f() // 3
	})
	try("constScope", func() any {
		const c = 1
		{
			const c = 2
			_ = c
		}
		return c
	})
	try("typeScope", func() any {
		type T int
		var v T = 5
		return int(v)
	})
}
