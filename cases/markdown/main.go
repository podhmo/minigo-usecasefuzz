package main

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

// Mini Markdown -> HTML: headings, unordered lists, paragraphs, and
// **bold** / `code` inline markup. A representative "render a doc
// snippet" text transform.
const doc = `# Title

Some *intro* text with **bold** and ` + "`code`" + `.

## List
- first
- second **item**

A closing paragraph.`

var (
	boldRe = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	codeRe = regexp.MustCompile("`([^`]+)`")
	emRe   = regexp.MustCompile(`\*([^*]+)\*`)
)

func inline(s string) string {
	s = html.EscapeString(s)
	s = boldRe.ReplaceAllString(s, "<strong>$1</strong>")
	s = codeRe.ReplaceAllString(s, "<code>$1</code>")
	s = emRe.ReplaceAllString(s, "<em>$1</em>")
	return s
}

func main() {
	var out strings.Builder
	inList := false
	for _, line := range strings.Split(doc, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case line == "":
			if inList {
				out.WriteString("</ul>\n")
				inList = false
			}
		case strings.HasPrefix(line, "## "):
			out.WriteString("<h2>" + inline(line[3:]) + "</h2>\n")
		case strings.HasPrefix(line, "# "):
			out.WriteString("<h1>" + inline(line[2:]) + "</h1>\n")
		case strings.HasPrefix(line, "- "):
			if !inList {
				out.WriteString("<ul>\n")
				inList = true
			}
			out.WriteString("  <li>" + inline(line[2:]) + "</li>\n")
		default:
			out.WriteString("<p>" + inline(line) + "</p>\n")
		}
	}
	if inList {
		out.WriteString("</ul>\n")
	}
	fmt.Print(out.String())
}
