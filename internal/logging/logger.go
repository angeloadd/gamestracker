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

	logger := slog.New(slog.NewJSONHandler(stderr, &slog.HandlerOptions{
		Level:       slogLevel,
		AddSource:   true,
		ReplaceAttr: redactSensitiveFields,
	}))

	if !ok {
		logger.Warn(fmt.Sprintf("LogLevel[%s] not valid, using 'Info' as default. Accepted: nope, debug, info, warn, error", logLevel))
	}

	return logger
}

func redactSensitiveFields(_ []string, a slog.Attr) slog.Attr {
	sensitiveKeys := []string{"password", "token", "secret", "authorization"}

	for _, key := range sensitiveKeys {
		if a.Key == key {
			return slog.String(a.Key, "[REDACTED]")
		}
	}

	return a
}
