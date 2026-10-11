package plug

import (
	"strings"

	"github.com/podhmo/minigo/examples/extract-schema/extschema"
)

func Install(c *extschema.Config) {
	c.DescField(func(f any, s string) string {
		return "field: " + s
	})
	c.DescDecl(func(d any, s string) string {
		return strings.ToUpper(s)
	})
}
