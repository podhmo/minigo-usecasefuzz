package main

// Go rejects: cannot take the address of m["a"].
func main() {
	m := map[string]int{"a": 1}
	_ = &m["a"]
}
