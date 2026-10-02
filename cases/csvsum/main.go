package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Parse an embedded TSV "spreadsheet", group and sum, emit an aligned
// report — a small ETL-style text job.
const data = `date	item	qty	price
2025-10-01	pen	3	120
2025-10-01	pad	2	300
2025-10-02	pen	1	120
2025-10-03	pad	1	300
2025-10-03	book	4	1500
bad	rows	shorter`

func main() {
	totals := map[string]int{}
	items := []string{}
	seen := map[string]bool{}
	bad := 0
	totalRev := 0
	for i, line := range strings.Split(data, "\n") {
		if i == 0 {
			continue
		}
		f := strings.Split(line, "\t")
		if len(f) != 4 {
			bad++
			continue
		}
		qty, err1 := strconv.Atoi(f[2])
		price, err2 := strconv.Atoi(f[3])
		if err1 != nil || err2 != nil {
			bad++
			continue
		}
		totals[f[1]] += qty
		totalRev += qty * price
		if !seen[f[1]] {
			seen[f[1]] = true
			items = append(items, f[1])
		}
	}
	sort.Strings(items)
	for _, it := range items {
		fmt.Printf("%-8s qty=%d\n", it, totals[it])
	}
	fmt.Println("revenue:", totalRev, "bad rows:", bad)
}
