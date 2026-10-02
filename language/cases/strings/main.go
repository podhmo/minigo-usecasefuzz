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
	try("len", func() any { return len("héllo") }) // 6 bytes
	try("index", func() any { return "abc"[1] })   // byte 98
	try("indexOOB", func() any { s := "abc"; i := 5; return s[i] })
	try("indexNeg", func() any { s := "abc"; i := -1; return s[i] })
	try("slice", func() any { return "hello"[1:3] })
	try("sliceFull", func() any { s := "hello"; return s[:] })
	try("sliceOOB", func() any { s := "abc"; lo, hi := 2, 10; return s[lo:hi] })
	try("sliceBadOrder", func() any { s := "abc"; lo, hi := 2, 1; return s[lo:hi] })
	try("concat", func() any { return "a" + "b" + "c" })
	try("concatAssign", func() any { s := "a"; s += "b"; return s })
	try("compare", func() any { return "abc" < "abd" })
	try("equal", func() any { return "x" == "x" })
	try("empty", func() any { return len("") })
	var rangecount int
	var runes []rune
	for _, r := range "héllo" {
		rangecount++
		runes = append(runes, r)
	}
	fmt.Printf("range-runes: %d %c\n", rangecount, runes)
	var idxs []int
	for i := range "abc" {
		idxs = append(idxs, i)
	}
	fmt.Printf("range-idx: %v\n", idxs)
	var bad []rune
	for _, r := range "a\xffb" { // invalid utf8 → RuneError
		bad = append(bad, r)
	}
	fmt.Printf("range-bad: %v\n", bad)
	try("strIterBytes", func() any {
		n := 0
		s := "héllo"
		for i := 0; i < len(s); i++ {
			n += int(s[i])
		}
		return n
	})
	try("multibyteSlice", func() any { return "héllo"[1:3] })
	try("charlit", func() any { return 'あ' })
	try("escape", func() any { return "a\nb\t\\\"'\\x41A" })
	try("backtick", func() any { return `a\nb` }) // raw literal keeps backslash
}
