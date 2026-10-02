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
	try("complexNew", func() any { return complex(1, 2) })
	try("realOf", func() any { return real(1 + 2i) })
	try("imagOf", func() any { return imag(1 + 2i) })
}
