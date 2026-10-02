package main

import (
	"fmt"
	"slices"
	"strings"
)

// Script-defined generics over slices/maps — the "write a small
// collection library" job.
func Map[T any, R any](xs []T, f func(T) R) []R {
	out := make([]R, 0, len(xs))
	for _, x := range xs {
		out = append(out, f(x))
	}
	return out
}

func Filter[T any](xs []T, f func(T) bool) []T {
	out := []T{}
	for _, x := range xs {
		if f(x) {
			out = append(out, x)
		}
	}
	return out
}

func Sum[T ~int | ~int64 | ~float64](xs []T) T {
	var acc T
	for _, x := range xs {
		acc += x
	}
	return acc
}

func KeysOf[K comparable, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func main() {
	xs := []int{1, 2, 3, 4, 5}
	fmt.Println("sum:", Sum(xs))
	fmt.Println("map:", Map(xs, func(x int) string { return fmt.Sprintf("#%d", x) }))
	fmt.Println("filter:", Filter(xs, func(x int) bool { return x%2 == 0 }))

	m := map[string]int{"b": 2, "a": 1}
	ks := KeysOf(m)
	slices.Sort(ks)
	fmt.Println("keys:", ks)

	fmt.Println("strs:", Map([]string{"a", "b"}, strings.ToUpper))
	i := slices.IndexFunc(xs, func(x int) bool { return x == 3 })
	fmt.Println("index of 3:", i)
}
