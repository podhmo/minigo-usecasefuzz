package main

import (
	"fmt"
	"reflect"
)

type S struct{}

func (S) f() {}
func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic")
		}
	}()
	fmt.Println(reflect.TypeOf(S{}).NumMethod())
}
