//go:build codegen

package gen

import (
	"fmt" // imported and not used — compile error in real Go

	"example.com/m/destination"
	"example.com/m/source"
	"github.com/podhmo/minigo/examples/convert-define/define"
)

func main() {
	define.Convert(func(c *define.Config, dst *destination.DstUser, src *source.SrcUser) {
		c.Map(dst.Name, src.Name)
		ghost()               // undefined identifier — never evaluated, only walked
		var x int = "no"      // type error — walked, ignored
		_ = x
	})
}
