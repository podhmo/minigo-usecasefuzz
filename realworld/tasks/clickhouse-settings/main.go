// Command clickhouse-settings cross-checks the clickhouse datasource's Go Settings
// struct (json tags) against the frontend CHConfig TypeScript interface.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"minigo.dev/inspect"
)

func main() {
	root := os.Getenv("TARGET_DIR")
	pkg := inspect.DirOf(filepath.Join(root, "pkg/plugin"))
	goKeys := map[string]string{}
	for _, f := range inspect.Fields(inspect.Symbol(pkg, "Settings")) {
		if len(f.Names) == 0 {
			continue
		}
		key := strings.Split(tagGet(f.Tag, "json"), ",")[0]
		if key == "-" {
			continue
		}
		if key == "" {
			key = f.Names[0] // encoding/json default
		}
		goKeys[key] = f.Names[0]
	}
	fmt.Println("package state:", inspect.State(pkg))

	src, err := os.ReadFile(filepath.Join(root, "src/types/config.ts"))
	if err != nil {
		fmt.Println(err)
		return
	}
	body := string(src)
	i := strings.Index(body, "interface CHConfig")
	body = body[i:]
	body = body[:strings.Index(body, "\n}")]
	re := regexp.MustCompile(`(?m)^  (\w+)\??:`)
	tsKeys := map[string]bool{}
	for _, m := range re.FindAllStringSubmatch(body, -1) {
		tsKeys[m[1]] = true
	}

	var goOnly, tsOnly []string
	for k, n := range goKeys {
		if !tsKeys[k] {
			goOnly = append(goOnly, k+" ("+n+")")
		}
	}
	for k := range tsKeys {
		if _, ok := goKeys[k]; !ok {
			tsOnly = append(tsOnly, k)
		}
	}
	sort.Strings(goOnly)
	sort.Strings(tsOnly)
	fmt.Printf("go=%d ts=%d go-only=%d ts-only=%d\n", len(goKeys), len(tsKeys), len(goOnly), len(tsOnly))
	fmt.Println("go-only:", goOnly)
	fmt.Println("ts-only:", tsOnly)
}

// tagGet is a minimal reflect.StructTag.Get (bound reflect lacks it).
func tagGet(tag, key string) string {
	i := strings.Index(tag, key+":\"")
	if i < 0 {
		return ""
	}
	rest := tag[i+len(key)+2:]
	return rest[:strings.Index(rest, "\"")]
}
