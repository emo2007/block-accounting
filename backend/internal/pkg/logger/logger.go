package logger

import (
	"io"
	"log/slog"
)

// todo add levels Trace, Emergency
type Logger struct {
	*slog.Logger
}

func newLogger(lvl slog.Level, w io.Writer) *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(w, &slog.HandlerOptions{Level: lvl}),
	)
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func MapLevel(lvl string) slog.Level {
	switch lvl {
	case "dev", "local", "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
