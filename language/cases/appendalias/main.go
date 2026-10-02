package main

import "fmt"

func main() {
	// append may share backing array: writes alias
	a := make([]int, 0, 4)
	a = append(a, 1, 2, 3)
	b := a
	b = append(b, 99) // fits in cap → shares backing
	fmt.Println(len(a), cap(b))
	// append beyond cap copies
	c := a[:2]
	d := append(c, 7, 8, 9)
	d[0] = 42
	fmt.Println(a[0], d[0]) // 42 only on d? backing shared if cap allowed
	// append to nil slice
	var e []int
	e = append(e, 5)
	fmt.Println(e, len(e), cap(e))
	// append with spread
	f := append([]int{1}, []int{2, 3}...)
	fmt.Println(f)
	// append string-ish: []byte append
	bb := append([]byte("he"), []byte("llo")...)
	fmt.Println(string(bb))
}
