package plug

import "github.com/podhmo/minigo/examples/extract-schema/extschema"

// wrong shapes: none of these is func(*extschema.Config)
func NoArgs()                               {}
func TwoArgs(a, b *extschema.Config)        {}
func Returns(c *extschema.Config) error     { return nil }
func Value(c extschema.Config)              {}
func unexported(c *extschema.Config)        {}
func Method()                               {}
