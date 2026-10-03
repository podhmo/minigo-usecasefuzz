package main

import (
	"fmt"
	"reflect"
)

func main() { var a any = reflect.Int; var b any = uint(reflect.Int); fmt.Println(a == b) }
