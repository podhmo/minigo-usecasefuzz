package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Hand-rolled {{name}} templating — the "render a message" job when
// text/template is unavailable.
var varRe = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_.]+)\s*\}\}`)

func render(tmpl string, ctx map[string]string) string {
	missing := []string{}
	out := varRe.ReplaceAllStringFunc(tmpl, func(m string) string {
		name := varRe.FindStringSubmatch(m)[1]
		if v, ok := ctx[name]; ok {
			return v
		}
		missing = append(missing, name)
		return "<" + name + "?>"
	})
	if len(missing) > 0 {
		out += "\n(missing: " + strings.Join(missing, ", ") + ")"
	}
	return out
}

func main() {
	ctx := map[string]string{
		"user.name": "alice",
		"count":     "3",
		"site":      "example.com",
	}
	tmpl := `Hello {{user.name}},

You have {{ count }} new messages on {{site}}.
Unknown slot: {{nope}}`
	fmt.Println(render(tmpl, ctx))
}
