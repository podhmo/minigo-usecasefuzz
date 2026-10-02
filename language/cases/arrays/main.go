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
	try("literal", func() any { return [3]int{1, 2, 3} })
	try("literalPartial", func() any { return [5]int{1, 2} }) // zeros tail
	try("literalKeyed", func() any { return [5]int{3: 9, 0: 1} })
	try("dots", func() any { return len([...]int{1, 2, 3, 4}) })
	try("len", func() any { return len([3]int{}) })
	try("cap", func() any { return cap([3]int{}) })
	try("idx", func() any { return [3]int{7, 8, 9}[2] })
	try("idxOOB", func() any { a := [3]int{1}; i := 9; return a[i] })
	// value semantics: assignment copies
	try("copySem", func() any {
		a := [2]int{1, 2}
		b := a
		b[0] = 99
		return fmt.Sprintf("%v %v", a, b)
	})
	try("argCopy", func() any {
		f := func(a [2]int) { a[0] = 99 }
		x := [2]int{1, 2}
		f(x)
		return x
	})
	try("eq", func() any { return [2]int{1, 2} == [2]int{1, 2} })
	try("neq", func() any { return [2]int{1, 2} != [2]int{1, 3} })
	try("sliceOf", func() any { a := [4]int{1, 2, 3, 4}; s := a[1:3]; return fmt.Sprintf("%v %d", s, cap(s)) })
	try("sliceShare", func() any {
		a := [3]int{1, 2, 3}
		s := a[:]
		s[0] = 99
		return a // shares backing → [99 2 3]
	})
	try("range", func() any {
		n := 0
		for _, v := range [3]int{1, 2, 3} {
			n += v
		}
		return n
	})
	try("rangePtr", func() any {
		a := [2]int{5, 6}
		n := 0
		for _, v := range &a {
			n += v
		}
		return n
	})
	try("ptrIdx", func() any { a := &[3]int{1, 2, 3}; return a[1] })     // auto-deref
	try("ptrSlice", func() any { a := &[3]int{1, 2, 3}; return a[1:] }) // (*a)[1:]
	try("nested", func() any { return [2][2]int{{1, 2}, {3, 4}}[1][0] })
	try("ptrEq", func() any {
		a := [2]int{1, 2}
		pa, pb := &a, &a
		return pa == pb
	})
}
