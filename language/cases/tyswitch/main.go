package main

import "fmt"

type I interface{ M() }
type A struct{ X int }
func (a A) M() {}
type B struct{ S string }
func (b B) M() {}

func main() {
	var i any
	i = 42
	switch v := i.(type) {
	case int:
		fmt.Println("int", v+1)
	case string:
		fmt.Println("str", v)
	default:
		fmt.Println("other")
	}
	i = A{X: 1}
	switch i.(type) {
	case A:
		fmt.Println("A")
	case B:
		fmt.Println("B")
	default:
		fmt.Println("none")
	}
	// interface case in type switch
	i = B{S: "x"}
	switch i.(type) {
	case I:
		fmt.Println("is I")
	default:
		fmt.Println("not I")
	}
	// fallthrough in value switch
	n := 2
	switch n {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
		fallthrough
	case 3:
		fmt.Println("three-ish")
	}
	// tagless switch
	switch {
	case n < 0:
		fmt.Println("neg")
	case n == 2:
		fmt.Println("eq2")
	default:
		fmt.Println("pos")
	}
	// switch on typed nil interface value is fine
	var ni any = nil
	switch ni.(type) {
	case nil:
		fmt.Println("nil case")
	default:
		fmt.Println("non-nil")
	}
	// comma-ok assert in if
	if v, ok := i.(B); ok {
		fmt.Println("okB", v.S)
	}
}
