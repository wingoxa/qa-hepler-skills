package logging

import (
	"log/slog"
	"os"
	"sort"
	"strings"
	"time"
)

func NewFromEnv() *slog.Logger {
	level := new(slog.LevelVar)
	level.Set(parseLevel(os.Getenv("CML_LOG_LEVEL")))
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))
}

func Duration(start time.Time) string {
	return time.Since(start).Round(time.Millisecond).String()
}

func Keys(data map[string]any) string {
	if len(data) == 0 {
		return "-"
	}
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return strings.Join(keys, ",")
}

func parseLevel(value string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return slog.LevelDebug
	case "", "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "off", "none", "false", "0":
		return slog.Level(100)
	default:
		return slog.LevelInfo
	}
}
