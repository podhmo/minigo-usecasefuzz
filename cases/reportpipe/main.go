package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// A realistic end-to-end "produce a report" pipeline: JSON events in,
// group/aggregate, format a table — the shape of an ad-hoc LLM-agent
// data job.
const events = `[
	{"ts":"2025-10-01T09:00:00Z","user":"ann","action":"login","ms":12},
	{"ts":"2025-10-01T09:01:00Z","user":"bob","action":"view","ms":88},
	{"ts":"2025-10-01T09:02:00Z","user":"ann","action":"buy","ms":230},
	{"ts":"2025-10-01T09:03:00Z","user":"ann","action":"view","ms":45},
	{"ts":"2025-10-01T09:04:00Z","user":"bob","action":"buy","ms":510},
	{"ts":"2025-10-01T09:05:00Z","user":"cid","action":"view","ms":30}
]`

type Ev struct {
	User   string `json:"user"`
	Action string `json:"action"`
	Ms     int    `json:"ms"`
}

func main() {
	var evs []Ev
	if err := json.Unmarshal([]byte(events), &evs); err != nil {
		fmt.Println(err)
		return
	}
	type Agg struct {
		Count, TotalMs int
	}
	agg := map[string]*Agg{}
	for _, e := range evs {
		a := agg[e.User]
		if a == nil {
			a = &Agg{}
			agg[e.User] = a
		}
		a.Count++
		a.TotalMs += e.Ms
	}
	users := make([]string, 0, len(agg))
	for u := range agg {
		users = append(users, u)
	}
	sort.Strings(users)

	var b strings.Builder
	fmt.Fprintf(&b, "%-6s %5s %7s %7s\n", "user", "n", "total", "avg")
	for _, u := range users {
		a := agg[u]
		fmt.Fprintf(&b, "%-6s %5d %7d %7.1f\n", u, a.Count, a.TotalMs, float64(a.TotalMs)/float64(a.Count))
	}
	fmt.Print(b.String())

	// top action by count
	acts := map[string]int{}
	for _, e := range evs {
		acts[e.Action]++
	}
	top, topN := "", -1
	for a, n := range acts {
		if n > topN || (n == topN && a < top) {
			top, topN = a, n
		}
	}
	fmt.Println("top action:", top, topN)
}
