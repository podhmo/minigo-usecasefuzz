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
	try("lenStr", func() any { return len("abc") })
	try("lenArr", func() any { return len([4]int{}) })
	try("lenArrPtr", func() any { return len(&[4]int{}) }) // len of ptr-to-array works in Go!
	try("lenSlice", func() any { return len([]int{1, 2}) })
	try("lenMap", func() any { return len(map[string]int{"a": 1}) })
	try("lenNilMap", func() any { return len(map[string]int(nil)) })
	try("capSlice", func() any { return cap(make([]int, 2, 7)) })
	try("capArr", func() any { return cap([5]int{}) })
	try("capArrPtr", func() any { return cap(&[5]int{}) })
	try("capNilSlice", func() any { return cap([]int(nil)) })
	try("copyStrBytes", func() any {
		b := make([]byte, 3)
		n := copy(b, "hello")
		return fmt.Sprintf("%d %q", n, b)
	})
	try("deleteTyped", func() any {
		m := map[int]string{1: "a"}
		delete(m, 1)
		return len(m)
	})
	try("makeMap", func() any { return len(make(map[string]int, 10)) })
	try("makeChanNil", func() any { return cap(make(chan int)) })
	try("makeChanBuf", func() any { return cap(make(chan int, 3)) })
	try("newSlicePtr", func() any { p := new([]int); return *p == nil })
	try("minBuiltin", func() any { return min(1, 2) })       // Go 1.21 builtin
	try("maxBuiltin", func() any { return max(1.5, 2.5) })
	try("minStr", func() any { return min("a", "b") })
	try("minVar", func() any { return min(3, 1, 2) })        // variadic builtin
	try("clearMap", func() any {
		m := map[string]int{"a": 1}
		clear(m)
		return len(m)
	})
	try("clearSlice", func() any {
		s := []int{1, 2, 3}
		clear(s)
		return s
	})
	try("printRet", func() any {
		print("x")   // builtin print → stderr in Go, stdout in minigo?
		println("y")
		return "done"
	})
	try("panicAny", func() any {
		panic(struct{ X int }{X: 1}) // panic with struct value
	})
	try("panicNilIface", func() any {
		var e error
		panic(e) // panic(nil interface)
	})
}
