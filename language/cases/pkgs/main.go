package main

import (
	"fmt"
	"langfuzz/cases/pkgs/sub"
)

var x = sub.V() + 1

func main() {
	fmt.Printf("%d %d %s\n", x, sub.W, sub.Name())
}
