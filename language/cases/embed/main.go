package main

import "fmt"

type Base struct{ N int }
func (b Base) Hello() string { return fmt.Sprintf("base%d", b.N) }
func (b *Base) Set(n int)    { b.N = n }

type Mid struct{ Base }
type Outer struct {
	Mid
	Tag string
}
func (o Outer) Who() string { return "outer" }

type Iface interface{ Hello() string }

func main() {
	o := Outer{Mid: Mid{Base: Base{N: 3}}, Tag: "t"}
	fmt.Println(o.Hello()) // promoted through two levels
	fmt.Println(o.Who())
	p := &o
	p.Set(9) // pointer-receiver promoted via addressable outer
	fmt.Println(o.N)
	// interface satisfaction through embedding
	var ifc Iface = o
	fmt.Println(ifc.Hello())
	// embedded field name access
	fmt.Println(o.Mid.Base.N)
	// method shadowing: Outer.Who vs promoted
}
