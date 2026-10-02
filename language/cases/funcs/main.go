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

func vaSum(xs ...int) int {
	n := 0
	for _, x := range xs {
		n += x
	}
	return n
}
func vaFirst(a int, rest ...int) int { return a*100 + len(rest) }
func two() (int, int)              { return 1, 2 }
func named() (a, b int) {
	a, b = 3, 4
	return
}
func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
func add(a, b int) int { return a + b }

type T struct{ V int }

func (t T) Get() int    { return t.V }
func (t *T) Bump()      { t.V++ }
func (t T) Plus(n int) int { return t.V + n }

func main() {
	try("vaNone", func() any { return vaSum() })
	try("vaMany", func() any { return vaSum(1, 2, 3) })
	try("vaSpread", func() any { return vaSum([]int{1, 2, 3}...) })
	try("vaMixed", func() any { return vaFirst(1, 2, 3, 4) })  // 103
	try("vaSliceParam", func() any {
		f := func(xs ...int) int { return len(xs) }
		return f(9, 9)
	})
	try("multi", func() any {
		a, b := two()
		return a*10 + b
	})
	try("multiDiscard", func() any {
		_, b := two()
		return b
	})
	try("namedRet", func() any {
		a, b := named()
		return a*10 + b
	})
	try("recursion", func() any { return fib(10) }) // 55
	try("fnValue", func() any {
		f := add
		return f(1, 2)
	})
	try("fnAsArg", func() any {
		apply := func(g func(int, int) int, x, y int) int { return g(x, y) }
		return apply(add, 3, 4)
	})
	try("fnNilCmp", func() any {
		var f func()
		return f == nil
	})
	try("callNilFn", func() any { var f func(); f(); return 1 }) // panic
	try("closureShare", func() any {
		fns := make([]func() int, 0)
		x := 0
		for i := 0; i < 3; i++ {
			x = i
			fns = append(fns, func() int { return x })
		}
		return fns[0]() + fns[1]() + fns[2]() // shared x=2 → 6
	})
	try("methodValue", func() any {
		t := T{V: 5}
		f := t.Get // bound method value, copies t
		t.V = 9
		return f() // 5
	})
	try("methodValuePtr", func() any {
		t := &T{V: 5}
		f := t.Bump // bound to *T
		f()
		return t.V // 6
	})
	try("methodExpr", func() any {
		f := T.Get // method expression: f(t)
		return f(T{V: 8})
	})
	try("methodExprArgs", func() any {
		f := T.Plus
		return f(T{V: 1}, 10) // 11
	})
	try("ptrMethodExpr", func() any {
		f := (*T).Bump
		t := &T{V: 1}
		f(t)
		return t.V // 2
	})
	try("closureRec", func() any {
		var f func(int) int
		f = func(n int) int {
			if n < 2 {
				return n
			}
			return f(n-1) + n
		}
		return f(5) // 4+3+2+1+... wait: f(5)=f(4)+5 = f(3)+4+5 = ... = 0+1+2+3+4+5=15
	})
	try("shadowBuiltin", func() any {
		len := 5 // shadows builtin in this scope
		_ = len
		return 5
	})
	try("anonDirect", func() any {
		return func() int { return 42 }()
	})
	try("multiAsArgs", func() any { return fmt.Sprint(two()) }) // multi-ret spread into variadic → "1 2"
}
