package main

import (
	"context"

	_ "embed"

	"github.com/alextanhongpin/go-slog-ring-buffer/bar"
	"github.com/alextanhongpin/go-slog-ring-buffer/logger"
	"github.com/rs/xid"
	"golang.org/x/exp/slog"
	"golang.org/x/mod/modfile"
)

func main() {
	ctx := context.Background()
	reqID := xid.New().String()

	// Logger is initialized per-request to control the level.
	l := logger.NewBuffer(ctx, reqID)
	defer l.Flush()

	ctx = logger.ContextWithValue(ctx, l)
	Foo(ctx)
}

func Foo(ctx context.Context) {
	logger.DebugCtx(ctx, "calling Foo", slog.Group("user", slog.String("name", "John")))
	logger.DebugCtx(ctx, "Foo called")

	bar.Bar(ctx)
}

// NOTE: This is one way of obtaining the pkg path ...
// The easier way? Just hardcode.
//
//go:embed go.mod
var gomod []byte

func Pkg() string {
	m, err := modfile.Parse("go.mod", gomod, nil)
	if err != nil {
		panic(err)
	}

	return m.Module.Mod.Path
}
