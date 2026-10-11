// Command routes cross-checks grafana's HTTP route registrations
// (pkg/api/api.go) against `swagger:route` annotations in pkg/api,
// using only the minigo.dev/inspect surface:
//
//   - route registrations live inside registerRoutes' function body —
//     read through inspect.Body/Nodes, the stmt/expr walk
//   - 117 of 121 annotations are method doc comments — inspect.Doc
//   - the remaining 4 are free comments detached from their decl
//     (quota.go:94, dataproxy.go:57, dataproxy.go:71, org_users.go:144)
//     — inspect.Comments reports them Free
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"minigo.dev/inspect"
)

type route struct {
	Method, Path, Handler string
	Pos                   string
}

func main() {
	root := os.Getenv("TARGET_DIR")
	dir := filepath.Join(root, "pkg", "api")
	pkg := inspect.DirOf(dir)
	fmt.Printf("package pkg/api: state=%s\n", inspect.State(pkg))

	// registerRoutes is a method on *HTTPServer.
	var routes []route
	for _, m := range inspect.Methods(inspect.Symbol(pkg, "HTTPServer")) {
		if inspect.Name(m) != "registerRoutes" {
			continue
		}
		body := inspect.Body(m)
		if body == nil {
			fmt.Println("registerRoutes has no body")
			return
		}
		walkRoutes(body, "", &routes)
	}

	docs, nDoc, nFree := collectSwagger(pkg)

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
	fmt.Printf("annotations: doc=%d free=%d\n", nDoc, nFree)
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

var methods = map[string]string{
	"Get": "GET", "Post": "POST", "Put": "PUT", "Delete": "DELETE", "Patch": "PATCH", "Any": "ANY",
}

// walkRoutes descends a function body, expanding Group calls'
// function literals under their path prefix — the script-side
// counterpart of the ast.Inspect walk, over inspect.Node views.
func walkRoutes(n *inspect.Node, prefix string, out *[]route) {
	if n.Kind == "CallExpr" {
		kids := inspect.Nodes(n)
		if len(kids) >= 2 && kids[1].Kind == "BasicLit" && kids[1].Detail == "STRING" {
			fun := kids[0]
			if fun.Kind == "SelectorExpr" {
				sel := inspect.Nodes(fun)
				p, _ := strconv.Unquote(kids[1].Text)
				if len(sel) == 2 && sel[1].Text == "Group" {
					for _, a := range kids[2:] {
						if a.Kind == "FuncLit" {
							for _, c := range inspect.Nodes(a) {
								if c.Role == "Body" {
									walkRoutes(c, prefix+p, out)
								}
							}
						}
					}
					return
				}
				if m, ok := methods[sel[1].Text]; ok {
					h := ""
					if last := kids[len(kids)-1]; last.Kind == "SelectorExpr" {
						h = inspect.Nodes(last)[1].Text
					}
					pos := ""
					if n.Pos != nil {
						pos = n.Pos.String()
					}
					*out = append(*out, route{Method: m, Path: prefix + p, Handler: h, Pos: pos})
					return
				}
			}
		}
	}
	for _, c := range inspect.Nodes(n) {
		walkRoutes(c, prefix, out)
	}
}

// collectSwagger reads `swagger:route` annotations through the
// inspect surface: method and decl doc comments (Doc), plus the free
// comments detached from any decl (Comments' Free reports).
func collectSwagger(pkg any) (map[string]bool, int, int) {
	out := map[string]bool{}
	nDoc, nFree := 0, 0
	scan := func(doc string, free bool) {
		for _, line := range strings.Split(doc, "\n") {
			// doc text comes marker-free; free comments keep their
			// // marker, exactly as written.
			t := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "//"))
			if !strings.HasPrefix(t, "swagger:route ") {
				continue
			}
			fs := strings.Fields(t)
			if len(fs) >= 3 {
				out[fs[1]+" "+normalize(fs[2])] = true
				if free {
					nFree++
				} else {
					nDoc++
				}
			}
		}
	}
	for _, d := range inspect.Decls(pkg) {
		scan(inspect.Doc(d), false)
	}
	for _, m := range inspect.Methods(pkg) {
		scan(inspect.Doc(m), false)
	}
	for _, f := range inspect.Files(pkg) {
		for _, c := range inspect.Comments(f) {
			if c.Free {
				scan(c.Text, true)
			}
		}
	}
	return out, nDoc, nFree
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
