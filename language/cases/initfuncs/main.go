package main

import "fmt"

var order []int
var v = func() int { order = append(order, 0); return 1 }()

func init() { order = append(order, 1) }
func init() { order = append(order, 2) }

func main() {
	fmt.Printf("%v %d\n", order, v) // [0 1 2] 1
}
