package main

import (
	"fmt"
	"strings"
)

// Render an aligned ASCII table with fmt's width verbs — the
// "print a report table to stdout" job.
func main() {
	type Row struct {
		Name  string
		Count int
		Price float64
	}
	rows := []Row{
		{"apple", 12, 0.5},
		{"banana", 3, 1.25},
		{"cherry-pie", 140, 0.05},
	}
	w := 0
	for _, r := range rows {
		if len(r.Name) > w {
			w = len(r.Name)
		}
	}
	sep := strings.Repeat("-", w+22)
	fmt.Println(sep)
	fmt.Printf("%-*s %5s %8s\n", w, "name", "count", "price")
	fmt.Println(sep)
	for _, r := range rows {
		fmt.Printf("%-*s %5d %8.2f\n", w, r.Name, r.Count, r.Price)
	}
	fmt.Println(sep)
}
