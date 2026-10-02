package convutil

import (
	"context"
	"fmt"

	"github.com/podhmo/minigo/examples/convert-define/model"
)

func Int64ToString(ctx context.Context, ec *model.ErrorCollector, v int64) string {
	return fmt.Sprintf("%d", v)
}
