package bar

import (
	"context"
	"math/rand"

	"github.com/alextanhongpin/go-slog-ring-buffer/logger"
	"golang.org/x/exp/slog"
)

func Bar(ctx context.Context) {
	logger.InfoCtx(ctx, "calling Bar")
	logger.DebugCtx(ctx, "SELECT 1 + $1", slog.Group("args", slog.Int64("$1", 42)))
	if rand.Intn(2) < 1 {
		logger.InfoCtx(ctx, "Bar called")
	} else {
		logger.ErrorCtx(ctx, "failed to call Bar")
	}
}
