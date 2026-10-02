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

func Id[T any](x T) T             { return x }
func Max[T int | float64](a, b T) T {
	if a > b {
		return a
	}
	return b
}
func Sum[T ~int](xs []T) T {
	var n T
	for _, x := range xs {
		n += x
	}
	return n
}

type Box[T any] struct{ V T }

func (b Box[T]) Get() T      { return b.V }
func (b *Box[T]) Set(v T)    { b.V = v }
func (b Box[T]) Map(f func(T) T) Box[T] { return Box[T]{V: f(b.V)} }

type Pair[K comparable, V any] struct {
	M map[K]V
}

func main() {
	try("idInt", func() any { return Id(42) })
	try("idStr", func() any { return Id("x") })
	try("idExplicit", func() any { return Id[int](7) })
	try("maxInt", func() any { return Max(1, 2) })
	try("maxFloat", func() any { return Max(1.5, 2.5) })
	try("sumInt", func() any { return Sum([]int{1, 2, 3}) })
	try("sumNamed", func() any {
		type MyInt int
		return Sum([]MyInt{1, 2, 3}) // ~int covers named
	})
	try("boxLit", func() any { return Box[int]{V: 5}.Get() })
	try("boxMethod", func() any {
		b := &Box[string]{V: "a"}
		b.Set("b")
		return b.Get()
	})
	try("boxMap", func() any { return Box[int]{V: 2}.Map(func(x int) int { return x * 3 }).V })
	try("pairGeneric", func() any {
		p := Pair[string, int]{M: map[string]int{"a": 1}}
		return p.M["a"]
	})
	try("genericZero", func() any {
		var b Box[int]
		return b.V // zero int
	})
	try("inferRet", func() any {
		f := func() int { return Id(3) }
		return f()
	})
	try("multiArgInfer", func() any {
		return Fst(1, 2)
	})
}

func Fst[T any](a, b T) T { return a }
