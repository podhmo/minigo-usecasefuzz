package main

// Go rejects: invalid array index 5 (out of bounds for 3-element array).
func main() {
	_ = [3]int{1, 2, 3}[5]
}
