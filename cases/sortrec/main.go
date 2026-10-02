package main

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

// Sort structured records several ways — the report-ordering job.
type Rec struct {
	Name  string
	Team  string
	Score int
}

func main() {
	recs := []Rec{
		{"nina", "b", 70},
		{"amy", "a", 90},
		{"zed", "a", 90},
		{"bob", "b", 85},
		{"max", "a", 70},
	}

	// stable multi-key: team asc, score desc, name asc
	sort.SliceStable(recs, func(i, j int) bool {
		a, b := recs[i], recs[j]
		if a.Team != b.Team {
			return a.Team < b.Team
		}
		if a.Score != b.Score {
			return a.Score > b.Score
		}
		return a.Name < b.Name
	})
	for _, r := range recs {
		fmt.Printf("%s/%s=%d\n", r.Team, r.Name, r.Score)
	}

	// slices.SortFunc on names
	names := []string{"delta", "alpha", "charlie"}
	slices.SortFunc(names, strings.Compare)
	fmt.Println(names)
	fmt.Println("sorted?", slices.IsSorted(names))

	// binary search
	pos, found := slices.BinarySearch(names, "charlie")
	fmt.Println("binary search:", pos, found)

	// sort.Search
	xs := []int{10, 20, 30, 40}
	i := sort.Search(len(xs), func(i int) bool { return xs[i] >= 25 })
	fmt.Println("search idx:", i)
}
