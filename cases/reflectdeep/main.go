package main

import (
	"fmt"
	"reflect"
)

// Semantic equality + type-introspection checks — the "did the decoded
// value match the golden struct" job in tests.
type Point struct {
	X, Y int
}

func main() {
	a := Point{1, 2}
	b := Point{1, 2}
	c := &a

	fmt.Println(reflect.DeepEqual(a, b))             // equal structs
	fmt.Println(reflect.DeepEqual(&a, &b))           // equal pointed values
	fmt.Println(reflect.DeepEqual(c, &a))            // same pointer
	fmt.Println(reflect.DeepEqual([]int{1}, []int{})) // len differs
	fmt.Println(reflect.DeepEqual([]int(nil), []int{})) // nil vs empty

	// Chan Send/Recv: drive a script chan through the facade
	raw := make(chan int, 2)
	ch := reflect.ValueOf(raw)
	ch.Send(reflect.ValueOf(10))
	ch.Send(reflect.ValueOf(20))
	v, ok := ch.Recv()
	fmt.Println(v.Int(), ok)
	close(raw)
	_, ok = ch.Recv()
	fmt.Println(ok)

	// Swapper-style: Copy between slices of same type
	src := reflect.ValueOf([]int{1, 2, 3})
	dst := reflect.ValueOf([]int{0, 0, 0})
	fmt.Println(reflect.Copy(dst, src), dst.Interface())
}
