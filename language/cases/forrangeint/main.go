package main

import "fmt"

func main() {
	// Go 1.22: range over int
	for i := range 3 {
		fmt.Print(i)
	}
	fmt.Println()
	n := 2
	for range n {
		fmt.Print("x")
	}
	fmt.Println()
}
