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

type MyInt int

func main() {
	try("uintptrConv", func() any { return uintptr(0) })
	try("intToIface", func() any {
		var a any = MyInt(3)
		_, isMy := a.(MyInt)
		_, isInt := a.(int)
		return fmt.Sprintf("%v %v", isMy, isInt) // true false
	})
}
