package plug2

import "github.com/podhmo/minigo/examples/extract-schema/extschema"

func Install(c *extschema.Config) {
	c.Override("example.com/m.Money", map[string]any{
		"type": "integer",
		"note": "from plug2",
	})
	c.DescField(func(f any, s string) string {
		return s + " [p2]"
	})
}
