package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JSONL pipeline: decode each line, enrich, re-marshal — the standard
// "transform a log file of JSON objects" job.
const input = `{"event":"click","user":"alice","n":1}
{"event":"view","user":"bob","n":3}
{"event":"click","user":"alice","n":2}
not-json
{"event":"view","user":"carol","n":5}`

func main() {
	bad := 0
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var rec map[string]any
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			bad++
			continue
		}
		rec["processed"] = true
		if user, ok := rec["user"].(string); ok {
			rec["user"] = strings.ToUpper(user)
		}
		b, err := json.Marshal(rec)
		if err != nil {
			fmt.Println("marshal:", err)
			continue
		}
		fmt.Println(string(b))
	}
	fmt.Println("bad lines:", bad)

	// json.Valid on a fragment
	fmt.Println("valid:", json.Valid([]byte(`{"a":[1,2]}`)), json.Valid([]byte(`{"a":`)))
}
