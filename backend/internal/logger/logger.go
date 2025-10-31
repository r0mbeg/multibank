package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Setup configures slog with rotation and custom text format
func Setup(level slog.Level) *slog.Logger {
	const logFile = "./logs/multibank.log"
	const maxSizeMB = 10
	const maxBackups = 5
	const maxAgeDays = 7
	const compress = true

	if err := os.MkdirAll(filepath.Dir(logFile), 0o755); err != nil {
		panic("cannot create log directory: " + err.Error())
	}

	rotator := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    maxSizeMB,
		MaxBackups: maxBackups,
		MaxAge:     maxAgeDays,
		Compress:   compress,
	}

	writer := io.MultiWriter(os.Stdout, rotator)

	handler := &prettyHandler{
		out:   writer,
		level: level,
	}

	l := slog.New(handler)
	slog.SetDefault(l)

	l.Info("logger initialized",
		"file", logFile,
		"max_size_mb", maxSizeMB,
		"max_backups", maxBackups,
		"max_age_days", maxAgeDays,
		"compress", compress,
	)

	return l
}

// prettyHandler — минимальная реализация slog.Handler
// с нужным форматированием: ISO8601, выравнивание уровня и key=value.
type prettyHandler struct {
	out   io.Writer
	level slog.Level
}

func (h *prettyHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	return lvl >= h.level
}

func (h *prettyHandler) Handle(_ context.Context, r slog.Record) error {
	timestamp := r.Time.Format("2006-01-02 15:04:05")
	level := levelString(r.Level)
	msg := r.Message

	// собираем key=value из атрибутов
	var attrs string
	r.Attrs(func(a slog.Attr) bool {
		attrs += fmt.Sprintf(" %s=%v", a.Key, a.Value)
		return true
	})

	line := fmt.Sprintf("%s %-5s %s%s\n", timestamp, level, msg, attrs)
	_, err := h.out.Write([]byte(line))
	return err
}

func (h *prettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *prettyHandler) WithGroup(name string) slog.Handler       { return h }

func levelString(l slog.Level) string {
	switch l {
	case slog.LevelDebug:
		return "DEBUG"
	case slog.LevelInfo:
		return "INFO"
	case slog.LevelWarn:
		return "WARN"
	case slog.LevelError:
		return "ERROR"
	default:
		return fmt.Sprintf("%d", l)
	}
}
