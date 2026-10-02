//go:build codegen

package gen

import (
	"github.com/podhmo/minigo/examples/convert-define/define"
)

func main() {
	// bare type names: SrcSelf/DstSelf live in this directory's
	// package — like c17-selfpkg but without the self-import.
	define.Convert(func(c *define.Config, dst *DstSelf, src *SrcSelf) {})
}
