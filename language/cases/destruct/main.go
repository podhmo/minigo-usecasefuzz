package main

import "fmt"

func pair() (int, string) { return 1, "one" }
func triple() (int, int, error) { return 2, 3, nil }
func maybe() (int, bool) { return 0, false }

func main() {
	a, b := pair()
	fmt.Println(a, b)
	x, y, err := triple()
	fmt.Println(x, y, err)
	if v, ok := maybe(); !ok {
		fmt.Println("notok", v)
	}
	// destructure from map comma-ok
	m := map[string]int{"k": 9}
	v, ok := m["k"]
	fmt.Println(v, ok)
	v2, ok2 := m["no"]
	fmt.Println(v2, ok2)
	// destructure range
	for i, s := range []string{"p", "q"} {
		fmt.Println(i, s)
	}
	for k, vv := range m {
		fmt.Println(k, vv)
	}
	// blank in destructure
	_, b2 := pair()
	fmt.Println(b2)
}
