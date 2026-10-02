package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Read fixture files, merge into one output file, verify — the
// "concatenate inputs" script job.
func main() {
	parts := []string{}
	for _, name := range []string{"a.txt", "b.txt"} {
		b, err := os.ReadFile(filepath.Join("data", name))
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		parts = append(parts, strings.TrimSpace(string(b)))
	}
	merged := strings.Join(parts, "\n---\n")

	tmp, err := os.MkdirTemp("", "merge")
	if err != nil {
		fmt.Println("mktemp:", err)
		return
	}
	defer os.RemoveAll(tmp)

	out := filepath.Join(tmp, "merged.txt")
	if err := os.WriteFile(out, []byte(merged), 0o644); err != nil {
		fmt.Println("write:", err)
		return
	}
	back, err := os.ReadFile(out)
	if err != nil {
		fmt.Println("reread:", err)
		return
	}
	fmt.Printf("wrote %d bytes, roundtrip=%v\n", len(back), string(back) == merged)
	fmt.Println("first bytes:", hex.EncodeToString(back[:8]))

	info, err := os.Stat(out)
	if err != nil {
		fmt.Println("stat:", err)
		return
	}
	fmt.Println("size via stat:", info.Size(), "isdir:", info.IsDir())
}
