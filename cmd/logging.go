package cmd

import (
	"context"
	"log/slog"
	"os"
)

var (
	debug    bool
	logLevel = slog.LevelInfo
	logger   *slog.Logger
)

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)
}

func SetupLogging(debugFlag bool) {
	debug = debugFlag
	if debug {
		logLevel = slog.LevelDebug
	}
	logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "logger", logger)
}
