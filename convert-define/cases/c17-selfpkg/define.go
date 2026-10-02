//go:build codegen

package gen

import (
	m "example.com/m"

	"github.com/podhmo/minigo/examples/convert-define/define"
)

func main() {
	define.Convert(func(c *define.Config, dst *m.DstSelf, src *m.SrcSelf) {})
}
