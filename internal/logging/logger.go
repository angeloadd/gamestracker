package logging

import (
	"fmt"
	"github.com/angeloadd/gamestracker/internal/config"
	"io"
	"log/slog"
	"strings"
)

var levelsStringToSlog = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func NewLogger(cfg config.Logs, stderr io.Writer) *slog.Logger {
	logLevel := cfg.Level

	if logLevel == "nope" {
		return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	}

	slogLevel, ok := levelsStringToSlog[strings.ToLower(logLevel)]

	if !ok {
		slogLevel = slog.LevelInfo
	}

	logger := makeLogger(stderr, slogLevel)

	if !ok {
		logger.Warn(fmt.Sprintf("LogLevel[%s] not valid, using 'Info' as default. Accepted: nope, debug, info, warn, error", logLevel))
	}

	return logger
}

func makeLogger(w io.Writer, level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}))
}
