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

type Iface interface{ M() int }
type A struct{ V int }

func (a A) M() int { return a.V }

type B struct{ V int }

func (b *B) M() int { return b.V } // pointer receiver: only *B satisfies

type S string

func (s S) M() int { return len(s) }

type Embedded interface {
	Iface // embedded interface
	N() int
}
type E2 struct{ A }

func (e E2) N() int { return 9 }

type Stringer interface{ String() string }
type SS struct{ N int }

func (s SS) String() string { return fmt.Sprintf("SS%d", s.N) }

type Niller interface{ M() }
type NP struct{}

func (n *NP) M() {}

func main() {
	try("dispatch", func() any {
		var i Iface = A{V: 7}
		return i.M()
	})
	try("dispatchPtr", func() any {
		var i Iface = &B{V: 8}
		return i.M()
	})
	try("dispatchNamedStr", func() any {
		var i Iface = S("hello")
		return i.M()
	})
	try("assertHit", func() any {
		var i any = A{V: 3}
		return i.(A).V
	})
	try("assertMiss", func() any {
		var i any = "s"
		return i.(int) // panic
	})
	try("assertOK", func() any {
		var i any = "s"
		v, ok := i.(string)
		return fmt.Sprintf("%q %v", v, ok)
	})
	try("assertOKMiss", func() any {
		var i any = "s"
		v, ok := i.(int)
		return fmt.Sprintf("%d %v", v, ok)
	})
	try("assertToIface", func() any {
		var i any = A{V: 2}
		return i.(Iface).M()
	})
	try("nilIface", func() any {
		var i any
		return i == nil
	})
	try("typedNil", func() any {
		var p *B = nil
		var i Iface = p
		return i == nil // false! typed nil
	})
	try("ifaceEq", func() any {
		var a, b any = A{V: 1}, A{V: 1}
		return a == b
	})
	try("ifaceNeq", func() any {
		var a, b any = A{V: 1}, A{V: 2}
		return a == b
	})
	try("ifaceDiffTypes", func() any {
		var a any = A{V: 1}
		var b any = "x"
		return a == b
	})
	try("ifacePanicCmp", func() any {
		var a, b any = []int{1}, []int{1}
		return a == b // panic: uncomparable
	})
	try("typeSwitch", func() any {
		var i any = "hi"
		switch v := i.(type) {
		case int:
			return fmt.Sprintf("int %d", v)
		case string:
			return fmt.Sprintf("str %s", v)
		default:
			return "?"
		}
	})
	try("typeSwitchMulti", func() any {
		var i any = 3.5
		switch i.(type) {
		case int, float64:
			return "num"
		default:
			return "other"
		}
	})
	try("typeSwitchNil", func() any {
		var i any
		switch i.(type) {
		case nil:
			return "nil"
		default:
			return "nonil"
		}
	})
	try("typeSwitchBind", func() any {
		var i any = 42
		switch v := i.(type) {
		case int:
			return v + 1 // narrowed to int
		default:
			return -1
		}
	})
	try("embedIface", func() any {
		var e Embedded = E2{A: A{V: 4}}
		return e.M()*10 + e.N() // 49
	})
	try("promotedMethod", func() any {
		var i Iface = E2{A: A{V: 5}} // M promoted from A
		return i.M()
	})
	try("stringerFmt", func() any {
		return fmt.Sprintf("%v", SS{N: 3}) // "SS3" via String()
	})
	try("stringerNoImpl", func() any {
		return fmt.Sprintf("%v", A{V: 1}) // "{1}" struct fmt
	})
	try("errIface", func() any {
		var e error = fmt.Errorf("oops %d", 42)
		return e.Error()
	})
	try("nilError", func() any {
		var e error = nil
		return e == nil
	})
	try("emptyAnyHold", func() any {
		var a any = func() {} // func in any
		return a != nil
	})
}
