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
	try("literal", func() any { return map[string]int{"a": 1, "b": 2}["a"] })
	try("missing", func() any { return map[string]int{}["x"] }) // zero value
	try("commaOkHit", func() any {
		v, ok := map[string]int{"a": 1}["a"]
		return fmt.Sprintf("%d %v", v, ok)
	})
	try("commaOkMiss", func() any {
		v, ok := map[string]int{}["a"]
		return fmt.Sprintf("%d %v", v, ok)
	})
	try("nilRead", func() any { var m map[string]int; return m["x"] }) // zero, no panic
	try("nilWrite", func() any { var m map[string]int; m["x"] = 1; return m }) // panic
	try("nilLen", func() any { var m map[string]int; return len(m) })
	try("nilDelete", func() any { var m map[string]int; delete(m, "x"); return 1 })
	try("nilRange", func() any {
		n := 0
		var m map[string]int
		for range m {
			n++
		}
		return n
	})
	try("deleteHit", func() any {
		m := map[string]int{"a": 1, "b": 2}
		delete(m, "a")
		_, ok := m["a"]
		return fmt.Sprintf("%d %v", len(m), ok)
	})
	try("deleteMiss", func() any { m := map[string]int{"a": 1}; delete(m, "z"); return len(m) })
	try("boolKey", func() any { m := map[bool]int{true: 1}; return m[true] })
	// deterministic small map — iterate to count (order irrelevant)
	try("rangeCount", func() any {
		m := map[string]int{"a": 1, "b": 2, "c": 3}
		keys, vals := 0, 0
		for k, v := range m {
			keys += len(k)
			vals += v
		}
		return fmt.Sprintf("%d %d", keys, vals)
	})
	try("keyOnly", func() any {
		n := 0
		for k := range map[string]int{"a": 1, "bb": 2} {
			n += len(k)
		}
		return n
	})
	try("assign", func() any { m := map[string]int{}; m["k"] = 5; return m["k"] })
	try("incr", func() any { m := map[string]int{}; m["k"]++; return m["k"] })      // read zero then write
	try("plusAssign", func() any { m := map[string]int{"k": 3}; m["k"] += 4; return m["k"] })
	try("mapOfSlice", func() any { m := map[string][]int{}; m["k"] = append(m["k"], 1); return m["k"] })
	try("nested", func() any { m := map[string]map[string]int{"a": {"b": 1}}; return m["a"]["b"] })
	try("elemAddr", func() any { m := map[string]int{"a": 1}; _ = m["a"]; return 1 }) // &m["a"] would reject; skip
}
