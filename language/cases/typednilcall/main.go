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

type Niller interface{ M() }
type NP struct{}

func (n *NP) M() {}

func main() {
	try("typedNilCall", func() any {
		var p *NP = nil
		var i Niller = p
		i.M() // calls method on nil receiver — works if M doesn't deref
		return "ok"
	})
}
