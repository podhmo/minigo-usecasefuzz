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

func main() {
	try("lifo", func() any {
		s := ""
		f := func() (r string) {
			defer func() { s += "A" }()
			defer func() { s += "B" }()
			defer func() { s += "C" }()
			return "x"
		}
		r := f()
		return r + s // "xCBA"
	})
	try("argsEvalAtDefer", func() any {
		s := ""
		f := func() {
			i := 1
			defer func(x int) { s += fmt.Sprint(x) }(i) // captures i=1 NOW
			i = 99
		}
		f()
		return s // "1"
	})
	try("namedRet", func() any {
		f := func() (r int) {
			defer func() { r += 10 }()
			return 5
		}
		return f() // 15
	})
	try("recoverRet", func() any {
		f := func() (r int) {
			defer func() {
				if recover() != nil {
					r = -1
				}
			}()
			panic("boom")
		}
		return f() // -1
	})
	try("recoverVal", func() any {
		var got any
		func() {
			defer func() { got = recover() }()
			panic("msg")
		}()
		return got
	})
	try("recoverInt", func() any {
		var got any
		func() {
			defer func() { got = recover() }()
			panic(42)
		}()
		return got
	})
	try("recoverOutside", func() any { return recover() }) // nil — not in defer
	try("panicNil", func() any {
		var got any
		func() {
			defer func() { got = recover() }()
			panic(nil)
		}()
		if got == nil {
			return "nil"
		}
		return fmt.Sprintf("%T", got) // *runtime.PanicNilError in Go ≥1.21
	})
	try("panicInDefer", func() any {
		var got any
		func() {
			defer func() { got = fmt.Sprintf("%v", recover()) }()
			defer func() { panic("second") }() // replaces first panic
			panic("first")
		}()
		return got // "second"
	})
	try("recoverStopsUnwind", func() any {
		order := ""
		func() {
			defer func() { order += "outer" }()
			func() {
				defer func() {
					recover()
					order += "inner"
				}()
				panic("x")
			}()
		}()
		return order // "innerouter" — outer still runs
	})
	try("deferInLoop", func() any {
		n := 0
		func() {
			for i := 0; i < 3; i++ {
				defer func() { n++ }()
			}
		}()
		return n // 3
	})
	try("deferMethod", func() any {
		s := ""
		func() {
			i := 0
			defer func() { i = 7 }()
			defer func() { s = fmt.Sprint(i) }() // reads i AFTER first defer ran? LIFO: s-defer runs first
		}()
		return s // "0" — s read before i=7 defer
	})
	try("deferReturnPanic", func() any {
		f := func() (r any) {
			defer func() { r = recover() }()
			panic("deep")
		}
		return f()
	})
	try("multiPanicSeq", func() any {
		var got []string
		func() {
			defer func() { recover() }()
			defer func() { got = append(got, "d1") }()
			defer func() { got = append(got, "d2"); panic("p2") }()
			panic("p1")
		}()
		return got // [d2 d1]
	})
	try("deferZero", func() any {
		f := func() int { return 1 }
		defer f() // defer discards return
		return 0
	})
	try("deferClosureCap", func() any {
		s := ""
		i := 1
		defer func() { s = fmt.Sprint(i) }() // closure reads CURRENT i at call time → 99
		i = 99
		return s
	})
}
