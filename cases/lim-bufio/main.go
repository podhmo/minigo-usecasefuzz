package main

import (
	"bufio"
	"fmt"
	"strings"
)

// Line scanning via bufio.Scanner — the standard "read lines" idiom.
func main() {
	s := bufio.NewScanner(strings.NewReader("one\ntwo\nthree"))
	for s.Scan() {
		fmt.Println("line:", s.Text())
	}
	fmt.Println("err:", s.Err())
}
