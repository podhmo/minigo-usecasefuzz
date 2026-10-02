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
	try("swap", func() any { a, b := 1, 2; a, b = b, a; return fmt.Sprintf("%d%d", a, b) })
	try("swapEval", func() any {
		// RHS fully evaluated before any assignment (Go spec)
		x := []int{1, 2}
		i := 0
		i, x[i] = 1, 99 // i becomes 1 first; x[i] uses OLD i=0 → x[0]=99
		return fmt.Sprintf("%d %v", i, x)
	})
	try("multiAssign", func() any { a, b, c := 1, 2, 3; return a + b + c })
	try("compoundAll", func() any {
		x := 10
		x += 1
		x -= 2
		x *= 3
		x /= 4
		x %= 5
		x <<= 2
		x >>= 1
		x |= 3
		x &= 6
		x ^= 1
		x &^= 1
		return x
	})
	try("compoundStr", func() any { s := "a"; s += "b"; return s })
	try("incdec", func() any {
		i := 5
		i++
		i--
		i--
		return i
	})
	try("fieldInc", func() any {
		type T struct{ N int }
		t := T{N: 1}
		t.N++
		return t.N
	})
	try("indexInc", func() any {
		s := []int{1, 2}
		s[1]++
		return s
	})
	try("mapInc", func() any {
		m := map[string]int{}
		m["k"]++
		m["k"] += 10
		return m["k"]
	})
	try("starInc", func() any {
		x := 1
		p := &x
		*p++
		return x
	})
	try("parenStar", func() any {
		x := 1
		p := &x
		(*p) += 10
		return x
	})
	try("nestedIdx", func() any {
		m := map[string][]int{"k": {0, 0}}
		m["k"][1] = 9
		return m["k"]
	})
	try("chanAssign", func() any {
		c := make(chan int, 1)
		c <- 7
		v := <-c
		return v
	})
	try("blankField", func() any {
		type T struct{ A int }
		t := T{}
		_ = t.A
		return 1
	})
}
