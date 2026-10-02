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

type V struct{ N int }

func (v V) Val() int   { return v.N }
func (v *V) Ptr() int  { return v.N }
func (v *V) Set(n int) { v.N = n }
func (v V) Both() int  { return v.N + 1 }

type E struct{ V } // embedded value
type EP struct{ *V }
type EN struct {
	Name string
	V
}

type IntList []int

func (l IntList) Sum() int {
	n := 0
	for _, x := range l {
		n += x
	}
	return n
}

type MyMap map[string]int

func (m MyMap) Get(k string) int { return m[k] }

type MyFunc func() int

func (f MyFunc) Call() int { return f() }

func main() {
	try("valRecv", func() any { return V{N: 3}.Val() })
	try("ptrRecvOnAddr", func() any {
		v := V{N: 3}
		return v.Ptr() // auto &v
	})
	try("ptrRecvSet", func() any {
		v := V{}
		v.Set(9)
		return v.N
	})
	try("ptrRecvOnPtr", func() any {
		v := &V{N: 4}
		return v.Val() // auto *v
	})
	try("ptrSetThruPtr", func() any {
		v := &V{}
		v.Set(5)
		return v.N
	})
	try("valRecvNoMutate", func() any {
		v := V{N: 1}
		v.Val()
		return v.N
	})
	try("embedPromote", func() any { return E{V: V{N: 5}}.Val() })
	try("embedPromoteSet", func() any {
		e := E{V: V{}}
		e.Set(7) // promoted *V method on addressable e
		return e.N
	})
	try("embedPtrPromote", func() any {
		e := EP{V: &V{N: 6}}
		return e.Val()
	})
	try("embedWithOwn", func() any {
		e := EN{Name: "x", V: V{N: 2}}
		return e.N + e.V.N
	})
	try("namedSliceMeth", func() any { return IntList{1, 2, 3}.Sum() })
	try("namedMapMeth", func() any { return MyMap{"a": 5}.Get("a") })
	try("namedFuncMeth", func() any {
		f := MyFunc(func() int { return 9 })
		return f.Call()
	})
	try("methodOnLiteral", func() any { return V{5}.Both() })
	try("chained", func() any {
		v := &V{N: 1}
		v.Set(v.Val() + v.Ptr())
		return v.N // 2
	})
}
