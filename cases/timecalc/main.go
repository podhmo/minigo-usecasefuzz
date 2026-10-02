package main

import (
	"fmt"
	"time"
)

// Date/duration arithmetic — scheduling-style calculations.
func main() {
	t0, err := time.Parse("2006-01-02", "2025-10-01")
	if err != nil {
		fmt.Println("parse:", err)
		return
	}
	fmt.Println("parsed:", t0.Format("2006-01-02 Mon"))
	fmt.Println("+90d:", t0.Add(90*24*time.Hour).Format("2006-01-02"))
	fmt.Println("unix:", t0.Unix())

	t1, _ := time.Parse("2006-01-02 15:04", "2025-12-31 23:59")
	fmt.Println("delta hours:", int(t1.Sub(t0).Hours()))

	d := 3*time.Hour + 25*time.Minute
	fmt.Println("dur:", d.String(), "mins:", int(d.Minutes()))

	// duration parsing (bound?)
	pd, err := time.ParseDuration("1h30m")
	fmt.Println("parsed dur:", pd, err)
}
