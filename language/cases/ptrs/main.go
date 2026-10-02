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
	try("new", func() any { p := new(int); *p = 5; return *p })
	try("newStruct", func() any { return new(P).X })
	try("ampLit", func() any { x := 42; return *&x })
	try("ampStructLit", func() any { p := &P{X: 1}; return p.X })
	try("deref", func() any { x := 3; p := &x; return *p })
	try("writeThru", func() any { x := 3; p := &x; *p = 9; return x })
	try("nilDeref", func() any { var p *int; return *p })
	try("nilEq", func() any { var p *int; return p == nil })
	try("ptrEq", func() any { x := 1; a, b := &x, &x; return a == b })
	try("ptrNeq", func() any { a, b := new(int), new(int); return a == b })
	try("ptrToPtr", func() any { x := 1; p := &x; pp := &p; return **pp })
	try("ptrField", func() any {
		p := &P{X: 1}
		p.X = 42
		return p.X
	})
	try("ptrArithNil", func() any {
		var p *P
		return p == nil
	})
	try("swapThruPtr", func() any {
		a, b := 1, 2
		pa, pb := &a, &b
		*pa, *pb = *pb, *pa
		return fmt.Sprintf("%d %d", a, b)
	})
	try("sliceElemAddr", func() any {
		s := []int{1, 2, 3}
		p := &s[1]
		*p = 99
		return s
	})
	try("fieldAddr", func() any {
		p := P{X: 1}
		px := &p.X
		*px = 7
		return p.X
	})
	try("nilFieldRead", func() any { var p *P; return p.X }) // panic nil deref
	try("fnPtrParam", func() any {
		bump := func(p *int) { *p++ }
		v := 1
		bump(&v)
		return v
	})
}
