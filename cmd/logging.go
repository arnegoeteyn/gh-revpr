package cmd

import (
	"log/slog"
	"os"
	"strings"

	"github.com/arnegoeteyn/gh-revpr/ui"
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

func setupLogging(debugFlag bool) {
	debug = debugFlag
	if debug {
		logLevel = slog.LevelDebug
	}
	logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)
}

func handleErr(err error, msg string) {
	if err != nil {
		slog.Debug(msg, "error", err)
		ui.Error("%s", strings.ToUpper(msg))
		os.Exit(1)
	}
}
