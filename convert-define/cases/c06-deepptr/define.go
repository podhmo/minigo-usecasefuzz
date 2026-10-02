//go:build codegen

package gen

import (
	"example.com/m/destination"
	"example.com/m/source"
	"github.com/podhmo/minigo/examples/convert-define/define"
)

func main() {
	define.Convert(func(c *define.Config, dst *destination.DstWrap, src *source.SrcWrap) {})
}
