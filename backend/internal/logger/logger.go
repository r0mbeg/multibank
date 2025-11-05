// internal/logger/logger.go

package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Setup configures slog with rotation and pretty text output.
// - Console: colored
// - File: plain text (no ANSI)
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

	// Один handler, который пишет ДВЕ строки: в консоль (с цветами) и в файл (без цветов)
	handler := &prettyHandler{
		level:        level,
		outFile:      rotator,
		outConsole:   os.Stdout,
		colorConsole: true,
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

// prettyHandler — slog.Handler с плоским форматированием:
// "2006-01-02 15:04:05 LEVEL msg key=value ..."
// Поддерживает WithAttrs/WithGroup и расцвечивает только консольный вывод.
type prettyHandler struct {
	level        slog.Level
	outFile      io.Writer
	outConsole   io.Writer
	colorConsole bool

	// Накопленные через WithAttrs() атрибуты (контекстные)
	attrs  []slog.Attr
	groups []string
}

func (h *prettyHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	return lvl >= h.level
}

func (h *prettyHandler) Handle(_ context.Context, r slog.Record) error {
	// timestamp/level/msg
	ts := r.Time
	if ts.IsZero() {
		ts = time.Now()
	}
	timestamp := ts.Format("2006-01-02 15:04:05")

	levelText := levelString(r.Level)
	msg := r.Message

	// Собираем key=value
	var sbFile strings.Builder
	var sbConsole strings.Builder

	// 1) Контекстные attrs из WithAttrs(...)
	for _, a := range h.attrs {
		appendAttr(&sbFile, h.groups, a, false)
		appendAttr(&sbConsole, h.groups, a, false)
	}

	// 2) Атрибуты самого рекорда
	r.Attrs(func(a slog.Attr) bool {
		appendAttr(&sbFile, h.groups, a, false)
		appendAttr(&sbConsole, h.groups, a, false)
		return true
	})

	// Формируем строки
	lineFile := fmt.Sprintf("%s %-5s %s%s\n", timestamp, levelText, msg, sbFile.String())
	lineConsole := lineFile

	// Расцвечиваем консольный вариант, если надо
	if h.colorConsole && h.outConsole != nil {
		levelColored := colorizeLevel(levelText)
		msgColored := color.CyanString(msg)

		lineConsole = fmt.Sprintf("%s %-5s %s%s\n",
			timestamp, levelColored, msgColored, sbConsole.String())
	}

	// Пишем в консоль (с цветами)
	if h.outConsole != nil {
		if _, err := io.WriteString(h.outConsole, lineConsole); err != nil {
			return err
		}
	}
	// Пишем в файл (без цветов)
	if h.outFile != nil {
		if _, err := io.WriteString(h.outFile, lineFile); err != nil {
			return err
		}
	}

	return nil
}

// WithAttrs должен вернуть НОВЫЙ handler c добавленными атрибутами.
// Обязательно копируем срезы, чтобы не портить старый handler.
func (h *prettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	nh := *h
	nh.attrs = append(append([]slog.Attr{}, h.attrs...), attrs...)
	return &nh
}

// WithGroup добавляет префикс к ключам (group.key=value)
func (h *prettyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	nh := *h
	nh.groups = append(append([]string{}, h.groups...), name)
	return &nh
}

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

// Расцветка уровня только для консоли
func colorizeLevel(level string) string {
	switch level {
	case "DEBUG":
		return color.MagentaString(level)
	case "INFO":
		return color.BlueString(level)
	case "WARN":
		return color.YellowString(level)
	case "ERROR":
		return color.RedString(level)
	default:
		return level
	}
}

// appendAttr добавляет " key=value" в буфер, учитывая группы и вложенность.
func appendAttr(sb *strings.Builder, groups []string, a slog.Attr, forceQuote bool) {
	//a = a.Resolve() // раскрываем lazy-значения

	// Группы внутри атрибутов
	if a.Value.Kind() == slog.KindGroup {
		// Префикс для вложенных значений: текущие группы + имя этого group
		prefix := append(groups, a.Key)
		for _, ga := range a.Value.Group() {
			appendAttr(sb, prefix, ga, forceQuote)
		}
		return
	}

	key := a.Key
	if key == "" {
		// Безымянный — пропускаем
		return
	}

	// Полный ключ с группами: g1.g2.key
	fullKey := key
	if len(groups) > 0 {
		fullKey = strings.Join(groups, ".") + "." + key
	}

	// Значение приводим к строке «по-людски»
	val := attrValueString(a.Value, forceQuote)

	sb.WriteByte(' ')
	sb.WriteString(fullKey)
	sb.WriteByte('=')
	sb.WriteString(val)
}

func attrValueString(v slog.Value, forceQuote bool) string {
	switch v.Kind() {
	case slog.KindString:
		s := v.String()
		// Если пробелы — в кавычки
		if forceQuote || strings.ContainsAny(s, " \t") {
			return fmt.Sprintf("%q", s)
		}
		return s
	case slog.KindTime:
		return v.Time().Format(time.RFC3339)
	case slog.KindInt64:
		return fmt.Sprintf("%d", v.Int64())
	case slog.KindUint64:
		return fmt.Sprintf("%d", v.Uint64())
	case slog.KindFloat64:
		return fmt.Sprintf("%g", v.Float64())
	case slog.KindBool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case slog.KindDuration:
		return v.Duration().String()
	case slog.KindAny:
		// Попробуем fmt
		return fmt.Sprintf("%v", v.Any())
	default:
		return v.String()
	}
}

// Err — удобный атрибут для ошибок
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
