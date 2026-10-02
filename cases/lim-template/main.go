package main

import (
	"fmt"
	"strings"
	"text/template"
)

// Render through text/template — the canonical templating answer.
func main() {
	t := template.Must(template.New("mail").Parse(
		"Hello {{.Name}}, you have {{.Count}} items{{if .Vip}} (VIP){{end}}."))
	var b strings.Builder
	err := t.Execute(&b, map[string]any{"Name": "ann", "Count": 3, "Vip": true})
	fmt.Println(b.String(), err)
}
