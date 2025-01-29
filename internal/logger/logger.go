package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	timeFormat = "2006-01-02 15:04:05.000"
)

type CustomHandler struct {
	handler slog.Handler
	writer  io.Writer
}

type PrettyJSONHandler struct {
	handler slog.Handler
	writer  io.Writer
}

func Setup(env string) *slog.Logger {
	switch env {
	case "local":
		return slog.New(NewLocalHandler())
	case "dev", "prod":
		return slog.New(NewPrettyJSONHandler(os.Stdout, getLevel(env)))
	default:
		panic("unknown environment: " + env)
	}
}

func getLevel(env string) slog.Level {
	switch env {
	case "local", "dev":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}

func NewLocalHandler() *CustomHandler {
	return &CustomHandler{
		writer:  os.Stdout,
		handler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	}
}

func NewPrettyJSONHandler(w io.Writer, level slog.Level) *PrettyJSONHandler {
	return &PrettyJSONHandler{
		writer: w,
		handler: slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: level,
		}),
	}
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	var sb strings.Builder

	// Timestamp
	sb.WriteString(color.WhiteString(r.Time.Format(timeFormat)) + " ")

	// Level
	sb.WriteString(h.colorLevel(r.Level) + " ")

	// Message
	sb.WriteString(color.WhiteString(r.Message) + " ")

	// Attributes
	h.addAttributes(&sb, r)

	// Source
	h.addSource(&sb, r)

	sb.WriteString("\n")
	_, err := h.writer.Write([]byte(sb.String()))
	return err
}

func (h *CustomHandler) addAttributes(sb *strings.Builder, r slog.Record) {
	attrs := make(map[string]interface{})
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})

	if len(attrs) > 0 {
		jsonData, _ := json.MarshalIndent(attrs, "", "  ")
		sb.WriteString(color.CyanString("\n") + color.CyanString(string(jsonData)))
	}
}

func (h *CustomHandler) addSource(sb *strings.Builder, r slog.Record) {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()
	source := fmt.Sprintf("\n%s:%d", f.File, f.Line)
	sb.WriteString(color.MagentaString(source))
}

func (h *CustomHandler) colorLevel(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return color.HiBlackString("DEBUG")
	case slog.LevelInfo:
		return color.GreenString("INFO ")
	case slog.LevelWarn:
		return color.YellowString("WARN ")
	case slog.LevelError:
		return color.RedString("ERROR")
	default:
		return color.WhiteString(level.String())
	}
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomHandler{
		handler: h.handler.WithAttrs(attrs),
		writer:  h.writer,
	}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{
		handler: h.handler.WithGroup(name),
		writer:  h.writer,
	}
}

func (h *PrettyJSONHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *PrettyJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	logEntry := map[string]interface{}{
		"time":    r.Time.Format(time.RFC3339Nano),
		"level":   r.Level.String(),
		"message": r.Message,
	}

	r.Attrs(func(a slog.Attr) bool {
		logEntry[a.Key] = a.Value.Any()
		return true
	})

	fs := runtime.CallersFrames([]uintptr{r.PC})
	if f, _ := fs.Next(); f.File != "" {
		logEntry["source"] = fmt.Sprintf("%s:%d", f.File, f.Line)
	}

	jsonData, err := json.MarshalIndent(logEntry, "", "  ")
	if err != nil {
		return err
	}

	_, err = h.writer.Write(append(jsonData, '\n'))
	return err
}

func (h *PrettyJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyJSONHandler{
		handler: h.handler.WithAttrs(attrs),
		writer:  h.writer,
	}
}

func (h *PrettyJSONHandler) WithGroup(name string) slog.Handler {
	return &PrettyJSONHandler{
		handler: h.handler.WithGroup(name),
		writer:  h.writer,
	}
}
