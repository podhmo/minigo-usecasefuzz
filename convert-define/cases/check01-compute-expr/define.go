//go:build codegen

package gen

import (
	"example.com/m/destination"
	"example.com/m/source"
	"github.com/podhmo/minigo/examples/convert-define/define"
)

func main() {
	define.Convert(func(c *define.Config, dst *destination.B, src *source.A) {
		c.Compute(dst.V, src.V+"!")
	})
}
