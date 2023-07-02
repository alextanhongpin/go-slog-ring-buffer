package main

import (
	"math/rand"
	"os"
	"time"

	"golang.org/x/exp/slog"
)

// TODO: Pass this per request and improve way of toggling.
var requestLevel = new(slog.LevelVar)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: requestLevel,
	}))
	Foo(logger)
}

func Foo(l *slog.Logger) {
	start := time.Now()
	defer l.Debug("called Foo", slog.Duration("elapsed", time.Since(start)))

	time.Sleep(250 * time.Millisecond)

	Bar(l)
}

func Bar(l *slog.Logger) {
	start := time.Now()
	defer l.Debug("called Bar", slog.Duration("elapsed", time.Since(start)))

	isErr := rand.Intn(2) < 1
	if isErr {
		l.Error("failed to call Bar")
		requestLevel.Set(slog.LevelDebug)
	} else {
		l.Info("called Bar")
	}

	time.Sleep(250 * time.Millisecond)
}
