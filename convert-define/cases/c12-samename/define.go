//go:build codegen

package gen

import (
	"example.com/m/a"
	"example.com/m/b"
	"example.com/m/c"
	"github.com/podhmo/minigo/examples/convert-define/define"
)

func main() {
	define.Convert(func(cc *define.Config, dst *b.User, src *a.User) {})
	define.Convert(func(cc *define.Config, dst *c.User, src *a.User) {
		cc.Map(dst.UID, src.ID)
	})
}
