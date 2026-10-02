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

type P struct{ X, Y int }
type W struct {
	P
	Name string
}
type W2 struct {
	A int
	P // promotes X,Y — P.X ambiguous? no: W2 has own A only
}
type Inner struct{ V int }
type Outer struct {
	Inner
	Other int
}
type Embed2 struct {
	Inner
	V int // shadows promoted Inner.V
}

type P3 struct{ X, Y, Z int }
type Pair struct{ A, B int }

func main() {
	try("lit", func() any { return P{1, 2} })
	try("litKeyed", func() any { return P{Y: 5, X: 1} })
	try("litPartial", func() any { return P{X: 3} }) // Y=0
	try("zero", func() any { var p P; return p })
	try("eq", func() any { return P{1, 2} == P{1, 2} })
	try("neq", func() any { return P{1, 2} == P{1, 3} })
	try("assign", func() any {
		a := P{1, 2}
		b := a
		b.X = 99
		return fmt.Sprintf("%v %v", a, b)
	})
	try("ptrField", func() any { p := &P{1, 2}; p.X = 9; return *p })
	try("embedField", func() any { w := W{P: P{1, 2}, Name: "n"}; return w.X + w.Y })
	try("embedSet", func() any { w := W{}; w.X = 7; return w.P.X })
	try("embedName", func() any { w := W{P: P{1, 2}}; return w.P.X })
	try("embedDeep", func() any {
		o := Outer{Inner: Inner{V: 5}}
		return o.V // promoted through one level
	})
	try("shadowField", func() any { e := Embed2{Inner: Inner{V: 1}, V: 2}; return e.V })
	try("anonLit", func() any { return struct{ A int }{A: 1}.A })
	try("anonCmp", func() any { return struct{ A int }{1} == struct{ A int }{1} })
	try("ptrLit", func() any { p := &P{1, 2}; return p.X })
	try("nestedLit", func() any {
		type O struct{ I P }
		return O{I: P{1, 2}}.I.Y
	})
	try("swap", func() any {
		p := Pair{1, 2}
		p.A, p.B = p.B, p.A
		return p
	})
	try("structSlice", func() any { return []P{{1, 2}, {3, 4}}[1].X })
	try("structMap", func() any { return map[string]P{"k": {1, 2}}["k"].Y })
	try("fieldOrder", func() any { return P3{Z: 3} }) // X,Y=0
	try("unexported", func() any { return P{}.X })    // field access (all fields exported here; name lowercase? X is exported)
}
