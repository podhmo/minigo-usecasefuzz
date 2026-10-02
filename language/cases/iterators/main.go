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

func gen(yield func(int) bool) {
	for i := 0; i < 5; i++ {
		if !yield(i) {
			return
		}
	}
}
func gen2(yield func(int, string) bool) {
	for i := 0; i < 3; i++ {
		if !yield(i, "x") {
			return
		}
	}
}
func genEarly(yield func(int) bool) {
	for i := 0; i < 100; i++ {
		if !yield(i) {
			return
		}
	}
}
func genNoCall(yield func(int) bool) {}

func main() {
	try("basic", func() any {
		n := 0
		for x := range gen {
			n += x
		}
		return n // 10
	})
	try("earlyBreak", func() any {
		n := 0
		for x := range genEarly {
			if x == 2 {
				break
			}
			n += x
		}
		return n // 1 — break stops iteration; yield returns false
	})
	try("twoVar", func() any {
		sum, cnt := 0, 0
		for i, s := range gen2 {
			sum += i
			cnt += len(s)
		}
		return sum*10 + cnt // 33
	})
	try("empty", func() any {
		n := 0
		for range genNoCall {
			n++
		}
		return n
	})
	try("rangeIntZero", func() any {
		n := 0
		for i := range 0 {
			n += i
		}
		return n
	})
	try("rangeIntNamed", func() any {
		type N int
		var k N = 3
		n := 0
		for i := range k { // range over named int — Go 1.22 ok
			n += int(i)
		}
		return n
	})
	try("nested", func() any {
		total := 0
		for x := range gen {
			for y := range gen {
				total += x * y
			}
		}
		return total // (0+1+2+3+4)^2 = 100
	})
}
