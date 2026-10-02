package main

import (
	"encoding/csv"
	"fmt"
	"strings"
)

// CSV reading via encoding/csv.
func main() {
	r := csv.NewReader(strings.NewReader("a,b\n1,2\n"))
	recs, err := r.ReadAll()
	fmt.Println(recs, err)
}
