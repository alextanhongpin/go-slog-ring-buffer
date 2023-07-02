package main

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/rs/xid"
	"golang.org/x/exp/slog"
)

func main() {
	ctx := context.Background()
	Foo(ctx, NewLogger())
}

func Foo(ctx context.Context, l logger) {
	start := time.Now()
	defer l.DebugCtx(ctx, "called Foo", slog.Duration("elapsed", time.Since(start)))

	time.Sleep(250 * time.Millisecond)

	Bar(ctx, l)
}

func Bar(ctx context.Context, l logger) {
	start := time.Now()
	defer l.DebugCtx(ctx, "called Bar", slog.Duration("elapsed", time.Since(start)))

	isErr := rand.Intn(2) < 1
	if isErr {
		l.ErrorCtx(ctx, "failed to call Bar")
	} else {
		l.InfoCtx(ctx, "called Bar")
	}

	time.Sleep(250 * time.Millisecond)
}

type logger interface {
	ErrorCtx(ctx context.Context, msg string, attrs ...slog.Attr)
	DebugCtx(ctx context.Context, msg string, attrs ...slog.Attr)
	InfoCtx(ctx context.Context, msg string, attrs ...slog.Attr)
}

type Logger struct {
	*slog.LevelVar
	*slog.Logger
}

func NewLogger() *Logger {
	level := new(slog.LevelVar)
	logger := slog.
		New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})).
		With(slog.String("req_id", xid.New().String()))

	return &Logger{
		LevelVar: level,
		Logger:   logger,
	}
}

func (l *Logger) DebugCtx(ctx context.Context, msg string, attrs ...slog.Attr) {
	l.LogAttrs(ctx, slog.LevelDebug, msg, attrs...)
}

func (l *Logger) InfoCtx(ctx context.Context, msg string, attrs ...slog.Attr) {
	l.LogAttrs(ctx, slog.LevelInfo, msg, attrs...)
}

func (l *Logger) ErrorCtx(ctx context.Context, msg string, attrs ...slog.Attr) {
	l.LogAttrs(ctx, slog.LevelError, msg, attrs...)
}

func (l *Logger) LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	if level == slog.LevelError {
		// Enable up to DEBUG level.
		l.LevelVar.Set(slog.LevelDebug)
	}

	l.Logger.LogAttrs(ctx, level, msg, attrs...)
}
