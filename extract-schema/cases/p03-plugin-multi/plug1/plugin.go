package plug1

import "github.com/podhmo/minigo/examples/extract-schema/extschema"

func Install(c *extschema.Config) {
	c.Override("example.com/m.Money", map[string]any{
		"type": "string",
		"note": "from plug1",
	})
	c.DescField(func(f any, s string) string {
		return "[p1] " + s
	})
}
