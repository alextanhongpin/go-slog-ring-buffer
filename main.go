package main

import (
	"context"
	"io"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/rs/xid"
	"golang.org/x/exp/slog"
)

type ctxKey string

const (
	logKey   ctxKey = "log"
	reqIdKey ctxKey = "req_id"
)

func main() {
	ctx := context.Background()
	ctx = WithRequestID(ctx)

	// Logger is initialized per-request to control the level.
	logger := NewBufferLogger(ctx)
	defer logger.Flush()

	ctx = WithLogger(ctx, logger)
	Foo(ctx)
}

func Foo(ctx context.Context) {
	DebugCtx(ctx, "calling foo", slog.Group("user", slog.String("name", "John")))

	Bar(ctx)
}

func Bar(ctx context.Context) {
	if rand.Intn(2) < 1 {
		InfoCtx(ctx, "bar")
	} else {
		ErrorCtx(ctx, "bar")
	}
}

type BufferLogger struct {
	w       io.Writer
	reqID   string
	records []slog.Record
	level   slog.Level
}

func NewBufferLogger(ctx context.Context) *BufferLogger {
	reqID, ok := RequestID(ctx)
	if !ok {
		reqID = xid.New().String()
	}

	return &BufferLogger{
		reqID: reqID,
		w:     os.Stdout,
		level: slog.LevelInfo,
	}
}

func (l *BufferLogger) Flush() {
	logger := slog.New(slog.NewJSONHandler(l.w, &slog.HandlerOptions{
		AddSource: true,
		Level:     l.level, // Defaults to INFO.
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

func (l *BufferLogger) log(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip [Callers, log, and parent]

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.AddAttrs(attrs...)
	l.records = append(l.records, r)
}

func InfoCtx(ctx context.Context, msg string, attrs ...slog.Attr) {
	l, ok := Logger(ctx)
	if !ok {
		return
	}
	l.log(ctx, slog.LevelInfo, msg, attrs...)
}

func DebugCtx(ctx context.Context, msg string, attrs ...slog.Attr) {
	l, ok := Logger(ctx)
	if !ok {
		return
	}
	l.log(ctx, slog.LevelDebug, msg, attrs...)
}

func ErrorCtx(ctx context.Context, msg string, attrs ...slog.Attr) {
	l, ok := Logger(ctx)
	if !ok {
		return
	}
	l.level = slog.LevelDebug
	l.log(ctx, slog.LevelError, msg, attrs...)
}

func WithLogger(ctx context.Context, logger *BufferLogger) context.Context {
	return context.WithValue(ctx, logKey, logger)
}

func Logger(ctx context.Context) (*BufferLogger, bool) {
	l, ok := ctx.Value(logKey).(*BufferLogger)
	return l, ok
}

func WithRequestID(ctx context.Context) context.Context {
	return context.WithValue(ctx, reqIdKey, xid.New().String())
}

func RequestID(ctx context.Context) (string, bool) {
	key, ok := ctx.Value(reqIdKey).(string)
	return key, ok
}
