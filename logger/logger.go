package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"golang.org/x/exp/slog"
)

type ctxKey string

const (
	pkgPath        = "github.com/alextanhongpin/go-slog-ring-buffer"
	logKey  ctxKey = "log"
)

type Buffer struct {
	w       io.Writer
	reqID   string
	records []slog.Record
	level   slog.Level
}

func NewBuffer(ctx context.Context, reqID string) *Buffer {

	return &Buffer{
		reqID: reqID,
		w:     os.Stdout,
		level: slog.LevelInfo,
	}
}

func (l *Buffer) Flush() {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		// Remove source file, and show the condensed information.
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)

			// Show file path relative to project path.
			//
			// Transforms from
			// in: /Users/alextanhongpin/Documents/golang/src/github.com/alextanhongpin/go-slog-ring-buffer/main.go
			// out: main.go
			fileParts := strings.Split(source.File, pkgPath)
			file := fileParts[len(fileParts)-1]
			file = strings.TrimPrefix(file, "/")

			// Transforms from
			// in: github.com/alextanhongpin/go-slog-ring-buffer/bar.Bar"
			// out: bar.Bar
			funcParts := strings.Split(source.Function, pkgPath)
			fn := funcParts[len(funcParts)-1]
			fn = strings.TrimPrefix(fn, "/")

			return slog.String("src", fmt.Sprintf("%s:%d %s", file, source.Line, fn))
		}
		return a
	}

	logger := slog.New(slog.NewJSONHandler(l.w, &slog.HandlerOptions{
		AddSource:   true,
		Level:       l.level, // Defaults to INFO.
		ReplaceAttr: replace,
	}))

	ctx := context.Background()

	h := logger.Handler()
	for _, r := range l.records {
		if !h.Enabled(ctx, r.Level) {
			continue
		}
		// To prevent out-of-order of logs, set the time it is recorded as an
		// attribute, and update the logging time.
		r.AddAttrs(
			slog.Time("event_time", r.Time),
			slog.String("req_id", l.reqID),
		)
		r.Time = time.Now()
		h.Handle(ctx, r)
	}
}

func (l *Buffer) log(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip [Callers, log, and parent]

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.AddAttrs(attrs...)
	l.records = append(l.records, r)
}

func InfoCtx(ctx context.Context, msg string, attrs ...slog.Attr) {
	l, ok := ContextValue(ctx)
	if !ok {
		return
	}
	l.log(ctx, slog.LevelInfo, msg, attrs...)
}

func DebugCtx(ctx context.Context, msg string, attrs ...slog.Attr) {
	l, ok := ContextValue(ctx)
	if !ok {
		return
	}
	l.log(ctx, slog.LevelDebug, msg, attrs...)
}

func ErrorCtx(ctx context.Context, msg string, attrs ...slog.Attr) {
	l, ok := ContextValue(ctx)
	if !ok {
		return
	}
	l.level = slog.LevelDebug
	l.log(ctx, slog.LevelError, msg, attrs...)
}

func ContextWithValue(ctx context.Context, logger *Buffer) context.Context {
	return context.WithValue(ctx, logKey, logger)
}

func ContextValue(ctx context.Context) (*Buffer, bool) {
	l, ok := ctx.Value(logKey).(*Buffer)
	return l, ok
}
