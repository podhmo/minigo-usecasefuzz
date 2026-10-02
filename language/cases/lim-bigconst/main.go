package main

import "fmt"

// DESIGNED LIMITATION: untyped constants exceeding int64 / arbitrary-precision
// const arithmetic — minigo constants evaluate to int64/float64 scalars, so
// huge compile-time values overflow silently (Go keeps big.Int precision).
const huge = 1 << 100

func main() {
	fmt.Println(float64(huge)) // 1.267...e+30
	fmt.Println(huge / huge)   // 1
}
