package plug

import "github.com/podhmo/minigo/examples/extract-schema/extschema"

func Install(c *extschema.Config) {
	var m map[string]any
	m["boom"] = 1 // write to nil map
}
