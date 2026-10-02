//go:build codegen

package gen

import (
	"example.com/m/convutil"
	"example.com/m/destination"
	"example.com/m/source"
	"github.com/podhmo/minigo/examples/convert-define/define"
)

func main() {
	define.Convert(func(c *define.Config, dst *destination.DstOrder, src *source.SrcOrder) {
		c.Convert(dst.OrderID, src.ID, convutil.Int64ToString)
		c.Convert(dst.TotalText, src.Total, convutil.Int64ToString)
	})
}
