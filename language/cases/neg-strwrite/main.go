package main

// Go rejects: cannot assign to s[0] (strings are immutable).
func main() {
	s := "ab"
	s[0] = 'X'
	_ = s
}
