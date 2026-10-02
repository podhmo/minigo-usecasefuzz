package main

import (
	"bytes"
	"fmt"
	"strings"
)

// Efficient string assembly — strings.Builder plus fmt.Fprintf into
// a buffer, the idiomatic "build a big string" job.
func main() {
	var b strings.Builder
	for i := 0; i < 4; i++ {
		fmt.Fprintf(&b, "[%d]", i)
	}
	b.WriteString("|")
	b.WriteByte('x')
	fmt.Println(b.String())
	fmt.Println("len:", b.Len())

	// bytes.Buffer equivalent
	buf := bytes.NewBufferString("start")
	buf.WriteString("+more")
	fmt.Fprintf(buf, "=%d", 42)
	fmt.Println(buf.String())
	fmt.Println("buf len:", buf.Len())
}
