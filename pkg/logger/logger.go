package logger

import (
	"context"
	"log/slog"
)

const (
	LevelTrace = slog.Level(-8)
)

// Logger - логгер вокруг slog.Logger, задача которого - добавить уровень LevelTrace.
type Logger struct {
	*slog.Logger
}

func New(h slog.Handler) *Logger {
	l := Logger{
		Logger: slog.New(h),
	}

	return &l
}

func (l *Logger) Trace(msg string, args ...any) {
	l.TraceContext(context.Background(), msg, args...)
}

func (l *Logger) TraceContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, LevelTrace, msg, args...)
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
	}
}

// ReplaceAttrLevel использовать в HandlerOptions.ReplaceAttr,
// чтобы корректно отображать уровень LevelTrace.
func ReplaceAttrLevel(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)
		if level == LevelTrace {
			a.Value = slog.StringValue("TRACE")
		}
	}

	return a
}
