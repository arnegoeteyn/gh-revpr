package cmd

import (
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
