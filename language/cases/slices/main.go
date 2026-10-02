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
	try("literal", func() any { return []int{1, 2, 3} })
	try("makeLen", func() any { return make([]int, 3) })
	try("makeCap", func() any { s := make([]int, 2, 10); return cap(s) })
	try("appendGrow", func() any {
		a := []int{1}
		b := append(a, 2)
		return fmt.Sprintf("%v %v", a, b)
	})
	// aliasing: append within cap mutates shared array
	try("aliasAppend", func() any {
		a := make([]int, 1, 4)
		a[0] = 1
		b := append(a, 9)
		b[0] = 7
		return fmt.Sprintf("a=%v b=%v", a, b) // a[0] becomes 7
	})
	try("appendMulti", func() any { return append([]int{1}, 2, 3, 4) })
	try("appendSlice", func() any { return append([]int{1}, []int{2, 3}...) })
	try("appendNil", func() any { var s []int; return append(s, 1) })
	try("appendStr", func() any { return append([]byte("ab"), 'c') }) // []byte append
	try("copyBasic", func() any {
		d := []int{0, 0, 0}
		n := copy(d, []int{1, 2, 3, 4})
		return fmt.Sprintf("%d %v", n, d)
	})
	try("copyOverlap", func() any {
		d := []int{1, 2, 3, 4}
		n := copy(d[1:], d[:3])
		return fmt.Sprintf("%d %v", n, d)
	})
	try("copyShort", func() any {
		d := []int{0}
		n := copy(d, []int{1, 2, 3})
		return fmt.Sprintf("%d %v", n, d)
	})
	try("copyNil", func() any { var d []int; return copy(d, []int{1}) })
	try("reslice", func() any { s := []int{1, 2, 3}; return s[:2][1:] })
	try("resliceOOB", func() any { s := []int{1, 2, 3}; return s[1:][:5] })
	try("nilIdx", func() any { var s []int; return s[0] }) // panic
	try("nilSliceLen", func() any { var s []int; return len(s) })
	try("nilSliceRange", func() any {
		n := 0
		var s []int
		for _, v := range s {
			n += v
		}
		return n
	})
	try("idxNeg", func() any { s := []int{1}; i := -1; return s[i] })
	try("idxOOB", func() any { s := []int{1, 2}; return s[5] })
	try("nested", func() any { return [][]int{{1}, {2, 3}} })
	try("sliceEqNil", func() any { var s []int; return s == nil })
	try("sliceEqNil2", func() any { return []int{} == nil })
	try("arrToSlice", func() any { a := [3]int{1, 2, 3}; return a[1:] })
	try("lenArgNil", func() any { return len([]int(nil)) })
}
