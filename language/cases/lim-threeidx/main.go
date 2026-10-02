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
	try("threeIdx", func() any {
		s := []int{1, 2, 3, 4, 5}[1:3:4]
		return fmt.Sprintf("%v %d", s, cap(s))
	})
	try("threeIdxBad", func() any {
		s := []int{1, 2, 3}
		lo, hi, mx := 1, 2, 1
		return s[lo:hi:mx] // runtime panic: max < high
	})
}
