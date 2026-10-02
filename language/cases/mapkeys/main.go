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

type K struct{ A, B int }

func main() {
type K2 struct{ A, B int }


	try("structKey", func() any {
		m := map[K]string{{1, 2}: "x"}
		return m[K{1, 2}]
	})
	try("arrayKey", func() any {
		m := map[[2]int]int{[2]int{1, 2}: 9}
		return m[[2]int{1, 2}]
	})
	try("ifaceKey", func() any {
		m := map[any]int{"s": 1, 42: 2, true: 3}
		return m["s"] + m[42] + m[true]
	})
	try("ifaceKeyNil", func() any {
		m := map[any]int{nil: 7}
		return m[nil]
	})
	try("nanKey", func() any {
		z := 0.0
		nan := z / z // NaN at runtime
		m := map[float64]int{nan: 1}
		m[nan] = 2                       // NaN keys are unfindable: adds second entry
		_, ok := m[nan]
		return fmt.Sprintf("%d %v", len(m), ok)
	})
}
