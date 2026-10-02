package main

import (
	"fmt"
	"io"
	"strings"
)

// io.ReadAll over a strings.Reader.
func main() {
	b, err := io.ReadAll(strings.NewReader("payload"))
	fmt.Println(string(b), err)
}
