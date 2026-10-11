package plug

import "github.com/podhmo/minigo/examples/extract-schema/extschema"

func install(c *extschema.Config) {
	c.Override("example.com/m.User", map[string]any{"type": "null"})
}
