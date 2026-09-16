// Package logx configures the structured logger flai uses on stderr,
// following design/conventions/logging.md: five levels including fatal,
// key-value text on a terminal and JSON otherwise, controlled by --verbose,
// LOG_LEVEL, and LOG_FORMAT.
package logx

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

// LevelFatal sits above slog.LevelError: the process cannot continue.
const LevelFatal = slog.Level(12)

// Options control New.
type Options struct {
	Verbose    bool   // --verbose: debug level
	Level      string // LOG_LEVEL: debug, info, warn, error, fatal
	Format     string // LOG_FORMAT: text or json
	IsTerminal bool   // default format is text on a terminal, json otherwise
}

// FromEnv reads LOG_LEVEL and LOG_FORMAT.
func FromEnv() Options {
	return Options{Level: os.Getenv("LOG_LEVEL"), Format: os.Getenv("LOG_FORMAT")}
}

// New builds the logger. Unknown level or format values fall back to the
// defaults and are reported on the returned logger at warn.
func New(w io.Writer, opt Options) *slog.Logger {
	level, badLevel := parseLevel(opt.Level)
	if opt.Verbose {
		level = slog.LevelDebug
	}
	format, badFormat := parseFormat(opt.Format, opt.IsTerminal)
	hopts := &slog.HandlerOptions{Level: level, ReplaceAttr: replaceAttr}
	var h slog.Handler
	if format == "json" {
		h = slog.NewJSONHandler(w, hopts)
	} else {
		h = slog.NewTextHandler(w, hopts)
	}
	l := slog.New(h)
	if badLevel {
		l.Warn("unknown LOG_LEVEL, using info", "component", "logx", "value", opt.Level)
	}
	if badFormat {
		l.Warn("unknown LOG_FORMAT, using default", "component", "logx", "value", opt.Format)
	}
	return l
}

// Fatal logs at the fatal level. The caller exits afterwards.
func Fatal(l *slog.Logger, msg string, args ...any) {
	l.Log(context.Background(), LevelFatal, msg, args...)
}

// replaceAttr renames the standard keys to the convention's (ts, level,
// msg) and renders the custom fatal level.
func replaceAttr(_ []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.TimeKey:
		a.Key = "ts"
		if t, ok := a.Value.Any().(time.Time); ok {
			a.Value = slog.StringValue(t.UTC().Format("2006-01-02T15:04:05.000Z"))
		}
	case slog.LevelKey:
		if lv, ok := a.Value.Any().(slog.Level); ok {
			a.Value = slog.StringValue(levelName(lv))
		}
	}
	return a
}

func levelName(lv slog.Level) string {
	switch {
	case lv >= LevelFatal:
		return "FATAL"
	case lv >= slog.LevelError:
		return "ERROR"
	case lv >= slog.LevelWarn:
		return "WARN"
	case lv >= slog.LevelInfo:
		return "INFO"
	default:
		return "DEBUG"
	}
}

func parseLevel(s string) (slog.Level, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "info":
		return slog.LevelInfo, false
	case "debug":
		return slog.LevelDebug, false
	case "warn", "warning":
		return slog.LevelWarn, false
	case "error":
		return slog.LevelError, false
	case "fatal":
		return LevelFatal, false
	}
	return slog.LevelInfo, true
}

func parseFormat(s string, terminal bool) (string, bool) {
	def := "json"
	if terminal {
		def = "text"
	}
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "":
		return def, false
	case "text", "json":
		return strings.ToLower(strings.TrimSpace(s)), false
	}
	return def, true
}

// String renders a level for messages and tests.
func String(lv slog.Level) string { return levelName(lv) }
