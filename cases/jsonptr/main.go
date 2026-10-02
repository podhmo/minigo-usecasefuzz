package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// JSON-pointer-ish navigation over unmarshalled maps — safe nested
// access with type-switch errors, like reading a nested config.
const doc = `{
	"services": {
		"web": {"replicas": 3, "ports": [80, 443]},
		"db":  {"replicas": 1, "ports": [5432]}
	},
	"meta": {"version": 2}
}`

func get(root any, path string) (any, error) {
	cur := root
	for _, seg := range strings.Split(strings.TrimPrefix(path, "/"), "/") {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("%s: not an object", seg)
		}
		v, ok := m[seg]
		if !ok {
			return nil, fmt.Errorf("%s: missing", seg)
		}
		cur = v
	}
	return cur, nil
}

func main() {
	var root any
	if err := json.Unmarshal([]byte(doc), &root); err != nil {
		fmt.Println(err)
		return
	}
	for _, p := range []string{
		"/services/web/replicas",
		"/services/db/ports",
		"/meta/version",
		"/services/nope",
		"/meta/version/x",
	} {
		v, err := get(root, p)
		fmt.Printf("%-24s -> %v (err=%v)\n", p, v, err)
	}

	// numbers come back float64 — exercise the common strconv dance
	v, _ := get(root, "/services/web/replicas")
	n := int(v.(float64))
	fmt.Println("replicas+", n+1, strconv.Itoa(n))
}
