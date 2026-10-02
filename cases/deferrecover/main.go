package main

import (
	"errors"
	"fmt"
	"strings"
)

// panic/recover + deferred cleanup — the "sandbox a risky operation"
// control-flow pattern.
func parse(s string) (n int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parse panic: %v", r)
		}
	}()
	if s == "" {
		panic(errors.New("empty input"))
	}
	if !strings.Contains(s, ":") {
		panic("missing colon")
	}
	return len(s), nil
}

func cleanup() (err error) {
	defer fmt.Println("cleanup start")
	defer func() {
		fmt.Println("cleanup end")
	}()
	return errors.New("inner")
}

func main() {
	for _, in := range []string{"ok:1", "nocolon", ""} {
		n, err := parse(in)
		fmt.Printf("%q -> n=%d err=%v\n", in, n, err)
	}
	fmt.Println("cleanup err:", cleanup())
	fmt.Println("after")
}
