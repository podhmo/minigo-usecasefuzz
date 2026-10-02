package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Walk a fixture tree and print a relative-path listing with sizes —
// the "inventory a directory" script job.
func main() {
	rows := []string{}
	err := filepath.WalkDir("data", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel("data", path)
		if d.IsDir() {
			rows = append(rows, rel+"/")
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		rows = append(rows, fmt.Sprintf("%s %d", rel, info.Size()))
		return nil
	})
	if err != nil {
		fmt.Println("walk:", err)
		return
	}
	sort.Strings(rows)
	for _, r := range rows {
		fmt.Println(r)
	}

	// ReadDir flat listing too
	ents, err := os.ReadDir("data/sub")
	if err != nil {
		fmt.Println("readdir:", err)
		return
	}
	names := []string{}
	for _, e := range ents {
		names = append(names, e.Name())
	}
	fmt.Println("sub:", strings.Join(names, ","))
}
