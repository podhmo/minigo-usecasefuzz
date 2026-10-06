// Command grafana-swagger-spec cross-checks grafana pkg/api (+dtos) swagger
// annotations against the bundled generated spec public/api-merged.json,
// using only the minigo.dev/inspect surface for the Go side.
//
// Direction is annotation -> spec only: the spec covers every package's
// routes, and 4 of pkg/api's 121 swagger:route comments are free comments
// detached from their decl (quota.go:94, dataproxy.go:57, dataproxy.go:71,
// org_users.go:144) which inspect.Doc cannot reach — a reverse check would
// false-positive on both. inspect.Decls also lists only top-level decls:
// all 117 visible route docs are method docs, reachable via type + Methods.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"minigo.dev/inspect"
)

var (
	routeRe = regexp.MustCompile(`(?m)swagger:route\s+([A-Z]+)\s+(\S+)((?:\s+\S+)*?)\s*$`)
	respRe  = regexp.MustCompile(`(?m)^(\d{3}):\s*([A-Za-z0-9_]+)\s*$`)
	respAnn = regexp.MustCompile(`(?m)swagger:response\s+([A-Za-z0-9_]+)`)
	modAnn  = regexp.MustCompile(`(?m)swagger:model(?:\s+([A-Za-z0-9_]+))?\s*$`)
)

type ann struct {
	method, path, opID string
	responses          map[string]string // status -> ref name
	pos                string
}

func posOf(d *inspect.Decl) string {
	p := inspect.Pos(d)
	if p == nil {
		return ""
	}
	return fmt.Sprintf("%s:%d", filepath.Base(p.File), p.Line)
}

func scanDoc(doc, pos, declName string, routes *[]ann, respTypes, modelNames map[string]bool) {
	if doc == "" {
		return
	}
	for _, m := range routeRe.FindAllStringSubmatch(doc, -1) {
		fields := strings.Fields(m[3])
		opID := ""
		if len(fields) > 0 {
			opID = fields[len(fields)-1]
		}
		a := ann{method: m[1], path: m[2], opID: opID, responses: map[string]string{}, pos: pos}
		for _, r := range respRe.FindAllStringSubmatch(doc, -1) {
			a.responses[r[1]] = r[2]
		}
		*routes = append(*routes, a)
	}
	for _, m := range respAnn.FindAllStringSubmatch(doc, -1) {
		respTypes[m[1]] = true
	}
	for _, m := range modAnn.FindAllStringSubmatch(doc, -1) {
		name := m[1]
		if name == "" {
			name = declName // bare "// swagger:model": the type name is the model name
		}
		modelNames[name] = true
	}
}

func main() {
	root := os.Getenv("TARGET_DIR")

	var routes []ann
	respTypes := map[string]bool{}
	modelNames := map[string]bool{}
	for _, sub := range []string{"pkg/api", "pkg/api/dtos"} {
		pkg := inspect.DirOf(filepath.Join(root, filepath.FromSlash(sub)))
		fmt.Printf("package %s: state=%s\n", sub, inspect.State(pkg))
		for _, d := range inspect.Decls(pkg) {
			scanDoc(inspect.Doc(d), posOf(d), inspect.Name(d), &routes, respTypes, modelNames)
			if inspect.Kind(d) == "type" {
				for _, m := range inspect.Methods(d) {
					scanDoc(inspect.Doc(m), posOf(m), "", &routes, respTypes, modelNames)
				}
			}
		}
	}
	fmt.Printf("annotations: routes=%d responses=%d models=%d\n", len(routes), len(respTypes), len(modelNames))

	// ---- spec side ----
	raw, err := os.ReadFile(filepath.Join(root, "public", "api-merged.json"))
	if err != nil {
		fmt.Println(err)
		return
	}
	var spec map[string]any
	if err := json.Unmarshal(raw, &spec); err != nil {
		fmt.Println(err)
		return
	}
	paths := spec["paths"].(map[string]any)
	responses := spec["responses"].(map[string]any)
	definitions := spec["definitions"].(map[string]any)
	fmt.Printf("spec: paths=%d responses=%d definitions=%d\n", len(paths), len(responses), len(definitions))

	var missingPath, wrongOp, respDrift, respMissing, modelMissing []string
	for _, a := range routes {
		entry, ok := paths[a.path].(map[string]any)
		if !ok {
			missingPath = append(missingPath, fmt.Sprintf("%s %s (%s @%s)", a.method, a.path, a.opID, a.pos))
			continue
		}
		op, ok := entry[strings.ToLower(a.method)].(map[string]any)
		if !ok {
			missingPath = append(missingPath, fmt.Sprintf("%s %s (%s @%s)", a.method, a.path, a.opID, a.pos))
			continue
		}
		if specOp, _ := op["operationId"].(string); specOp != a.opID {
			wrongOp = append(wrongOp, fmt.Sprintf("%s %s: annotation=%s spec=%s", a.method, a.path, a.opID, specOp))
		}
		specResps, _ := op["responses"].(map[string]any)
		for status, ref := range a.responses {
			sr, ok := specResps[status].(map[string]any)
			if !ok {
				respDrift = append(respDrift, fmt.Sprintf("%s %s: status %s (%s) absent in spec", a.method, a.path, status, ref))
				continue
			}
			got, _ := sr["$ref"].(string)
			if got == "" {
				if schema, ok := sr["schema"].(map[string]any); ok {
					got, _ = schema["$ref"].(string)
				}
			}
			if got != "#/responses/"+ref && got != "#/definitions/"+ref {
				respDrift = append(respDrift, fmt.Sprintf("%s %s: status %s annotation=%s spec=%s", a.method, a.path, status, ref, got))
			}
		}
	}
	for name := range respTypes {
		if _, ok := responses[name]; !ok {
			respMissing = append(respMissing, name)
		}
	}
	for name := range modelNames {
		if _, ok := definitions[name]; !ok {
			modelMissing = append(modelMissing, name)
		}
	}
	sort.Strings(missingPath)
	sort.Strings(wrongOp)
	sort.Strings(respDrift)
	sort.Strings(respMissing)
	sort.Strings(modelMissing)

	fmt.Printf("drift: missing-path=%d wrong-opid=%d resp-drift=%d resp-missing=%d model-missing=%d\n",
		len(missingPath), len(wrongOp), len(respDrift), len(respMissing), len(modelMissing))
	for _, s := range missingPath {
		fmt.Println("missing-path:", s)
	}
	for _, s := range wrongOp {
		fmt.Println("wrong-opid:", s)
	}
	for _, s := range respDrift {
		fmt.Println("resp-drift:", s)
	}
	for _, s := range respMissing {
		fmt.Println("resp-missing:", s)
	}
	for _, s := range modelMissing {
		fmt.Println("model-missing:", s)
	}
}
