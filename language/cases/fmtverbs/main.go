package main

import "fmt"

type P struct{ X, Y int }

func main() {
	fmt.Printf("%b %o %x %X\n", 5, 8, 255, 255)
	fmt.Printf("%e %g\n", 12345.678, 0.00001234)
	fmt.Printf("%q %s\n", "a\nb", "hi")
	fmt.Printf("%c %c\n", 'A', 233)
	fmt.Printf("%5d|%-5d|%05d\n", 42, 42, 42)
	fmt.Printf("%.2f %8.2f\n", 3.14159, 3.14159)
	fmt.Printf("%v %+v %#v\n", P{1, 2}, P{1, 2}, P{1, 2})
	fmt.Printf("%v %v\n", []int{1, 2}, map[string]int{"a": 1})
	fmt.Printf("%T %T %T\n", 42, "s", P{})
	fmt.Printf("%T %T\n", []int{1}, map[string]bool{})
	fmt.Printf("%[1]d %[1]d\n", 7)
	fmt.Printf("%d%%\n", 50)
}
