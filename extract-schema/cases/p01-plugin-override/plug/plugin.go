package plug

import "github.com/podhmo/minigo/examples/extract-schema/extschema"

func Install(c *extschema.Config) {
	c.Override("example.com/m.Money", map[string]any{
		"type":        "string",
		"pattern":     "^[0-9]+(USD|EUR)$",
		"description": "money as '123USD'",
	})
	c.Override("time.Time", map[string]any{
		"type":   "string",
		"format": "date-time",
	})
}
