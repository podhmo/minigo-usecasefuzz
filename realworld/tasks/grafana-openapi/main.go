// Command routes cross-checks grafana's HTTP route registrations
// (pkg/api/api.go) against `swagger:route` annotations in pkg/api.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type route struct {
	Method, Path, Handler string
	Pos                   string
}

func main() {
	root := os.Getenv("TARGET_DIR")
	dir := filepath.Join(root, "pkg", "api")
	fset := token.NewFileSet()

	routes := collectRoutes(fset, filepath.Join(dir, "api.go"))
	docs := collectSwagger(fset, dir)

	documented, undocumented := 0, 0
	var missing []string
	for _, r := range routes {
		if !strings.HasPrefix(r.Path, "/api/") {
			continue
		}
		key := r.Method + " " + normalize(strings.TrimPrefix(r.Path, "/api"))
		if docs[key] {
			documented++
			delete(docs, key)
		} else {
			undocumented++
			missing = append(missing, key+"  "+r.Handler)
		}
	}
	var stale []string
	for k := range docs {
		stale = append(stale, k)
	}
	sort.Strings(missing)
	sort.Strings(stale)
	fmt.Printf("routes=%d api=%d documented=%d undocumented=%d swagger-only=%d\n",
		len(routes), documented+undocumented, documented, undocumented, len(stale))
	for i, m := range missing {
		if i == 10 {
			fmt.Println("  ...")
			break
		}
		fmt.Println("  missing:", m)
	}
	for i, s := range stale {
		if i == 10 {
			fmt.Println("  ...")
			break
		}
		fmt.Println("  swagger-only:", s)
	}
}

// normalize rewrites macaron-style params (:uid, *) to swagger's {uid}.
func normalize(p string) string {
	parts := strings.Split(p, "/")
	for i, s := range parts {
		if strings.HasPrefix(s, ":") || strings.HasPrefix(s, "{") {
			parts[i] = "{}"
		}
	}
	return strings.TrimSuffix(strings.Join(parts, "/"), "/")
}

var methods = map[string]string{
	"Get": "GET", "Post": "POST", "Put": "PUT", "Delete": "DELETE", "Patch": "PATCH", "Any": "ANY",
}

func collectRoutes(fset *token.FileSet, file string) []route {
	f, err := parser.ParseFile(fset, file, nil, parser.SkipObjectResolution)
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse:", err)
		os.Exit(1)
	}
	var out []route
	var walk func(n ast.Node, prefix string)
	walk = func(n ast.Node, prefix string) {
		ast.Inspect(n, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || len(call.Args) == 0 {
				return true
			}
			lit, ok := call.Args[0].(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			p, _ := strconv.Unquote(lit.Value)
			if sel.Sel.Name == "Group" {
				for _, a := range call.Args[1:] {
					if fl, ok := a.(*ast.FuncLit); ok {
						walk(fl.Body, prefix+p)
					}
				}
				return false
			}
			m, ok := methods[sel.Sel.Name]
			if !ok {
				return true
			}
			h := ""
			if last, ok := call.Args[len(call.Args)-1].(*ast.SelectorExpr); ok {
				h = last.Sel.Name
			}
			out = append(out, route{Method: m, Path: prefix + p, Handler: h, Pos: fset.Position(call.Pos()).String()})
			return true
		})
	}
	for _, d := range f.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok && fd.Body != nil {
			walk(fd.Body, "")
		}
	}
	return out
}

func collectSwagger(fset *token.FileSet, dir string) map[string]bool {
	out := map[string]bool{}
	ents, _ := os.ReadDir(dir)
	for _, e := range ents {
		n := e.Name()
		if !strings.HasSuffix(n, ".go") || strings.HasSuffix(n, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, n), nil, parser.ParseComments|parser.SkipObjectResolution)
		if err != nil {
			continue
		}
		for _, cg := range f.Comments {
			for _, c := range cg.List {
				t := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
				if !strings.HasPrefix(t, "swagger:route ") {
					continue
				}
				fs := strings.Fields(t)
				if len(fs) >= 3 {
					out[fs[1]+" "+normalize(fs[2])] = true
				}
			}
		}
	}
	return out
}
