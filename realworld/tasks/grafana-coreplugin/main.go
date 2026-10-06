// Command plugins cross-checks grafana's core backend plugin registry
// (Go constants in pkg/services/pluginsintegration/coreplugin) against
// the plugin.json manifests under public/app/plugins/datasource.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

type manifest struct {
	ID         string `json:"id"`
	Backend    bool   `json:"backend"`
	Executable string `json:"executable"`
}

func main() {
	root := os.Getenv("TARGET_DIR")
	goIDs := constStrings(filepath.Join(root, "pkg/services/pluginsintegration/coreplugin/coreplugins.go"))

	files, _ := filepath.Glob(filepath.Join(root, "public/app/plugins/datasource/*/plugin.json"))
	jsonIDs := map[string]manifest{}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			fmt.Println("read:", err)
			continue
		}
		var m manifest
		if err := json.Unmarshal(b, &m); err != nil {
			fmt.Println("json:", f, err)
			continue
		}
		jsonIDs[m.ID] = m
	}

	var out []string
	for id, m := range jsonIDs {
		_, inGo := goIDs[id]
		switch {
		case m.Backend && !inGo:
			out = append(out, "backend plugin without Go registration: "+id)
		case !m.Backend && inGo:
			out = append(out, "Go registration for non-backend plugin: "+id)
		}
	}
	for id, name := range goIDs {
		if _, ok := jsonIDs[id]; !ok {
			out = append(out, "Go constant without plugin.json: "+name+"="+id)
		}
	}
	sort.Strings(out)
	fmt.Printf("go=%d json=%d issues=%d\n", len(goIDs), len(jsonIDs), len(out))
	for _, s := range out {
		fmt.Println(" ", s)
	}
}

// constStrings reads string constants from a file's syntax without
// evaluating the package (no init, no imports).
func constStrings(file string) map[string]string {
	out := map[string]string{}
	f, err := parser.ParseFile(token.NewFileSet(), file, nil, parser.SkipObjectResolution)
	if err != nil {
		fmt.Println("parse:", err)
		return out
	}
	for _, d := range f.Decls {
		g, ok := d.(*ast.GenDecl)
		if !ok || g.Tok != token.CONST {
			continue
		}
		for _, sp := range g.Specs {
			vs := sp.(*ast.ValueSpec)
			for i, n := range vs.Names {
				if i < len(vs.Values) {
					if lit, ok := vs.Values[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
						v, _ := strconv.Unquote(lit.Value)
						out[v] = n.Name
					}
				}
			}
		}
	}
	return out
}
