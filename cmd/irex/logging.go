package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

/* -------------------- Elapsed Handler -------------------- */

type elapsedHandler struct {
	start   time.Time
	handler slog.Handler
}

func newElapsedHandler(h slog.Handler) slog.Handler {
	return &elapsedHandler{
		start:   time.Now(),
		handler: h,
	}
}

func (h *elapsedHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *elapsedHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(
		slog.Duration("elapsed", time.Since(h.start).Round(time.Millisecond)),
	)
	return h.handler.Handle(ctx, r)
}

func (h *elapsedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &elapsedHandler{
		start:   h.start,
		handler: h.handler.WithAttrs(attrs),
	}
}

func (h *elapsedHandler) WithGroup(name string) slog.Handler {
	return &elapsedHandler{
		start:   h.start,
		handler: h.handler.WithGroup(name),
	}
}

/* -------------------- Base Handlers -------------------- */

func NewJSONHandler(w io.Writer, level slog.Level) slog.Handler {
	return slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level})
}

func NewTextHandler(w io.Writer, level slog.Level) slog.Handler {
	return slog.NewTextHandler(w, &slog.HandlerOptions{Level: level})
}

/* -------------------- Color Handler -------------------- */

type colorHandler struct {
	out   io.Writer
	level slog.Level
	attrs []slog.Attr
	mu    sync.Mutex
}

func newColorHandler(w io.Writer, level slog.Level) slog.Handler {
	return &colorHandler{
		out:   w,
		level: level,
	}
}

func (c *colorHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= c.level
}

func (c *colorHandler) Handle(ctx context.Context, r slog.Record) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	color := levelColor(r.Level)
	ts := r.Time.Format(time.RFC3339)

	var elapsed string
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "elapsed" {
			elapsed = a.Value.String()
			return false
		}
		return true
	})

	_, err := fmt.Fprintf(
		c.out,
		"%s[%s] [%s] %-5s %s\x1b[0m\n",
		color,
		ts,
		elapsed,
		r.Level.String(),
		r.Message,
	)

	return err
}

func (c *colorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// avoid copying sync.Mutex by allocating a new handler
	newAttrs := make([]slog.Attr, 0, len(c.attrs)+len(attrs))
	newAttrs = append(newAttrs, c.attrs...)
	newAttrs = append(newAttrs, attrs...)
	return &colorHandler{
		out:   c.out,
		level: c.level,
		attrs: newAttrs,
	}
}

func (c *colorHandler) WithGroup(name string) slog.Handler {
	// return a new handler to avoid sharing the mutex
	return &colorHandler{
		out:   c.out,
		level: c.level,
		attrs: append([]slog.Attr{}, c.attrs...),
	}
}

func levelColor(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return "\x1b[36m"
	case slog.LevelWarn:
		return "\x1b[33m"
	case slog.LevelError:
		return "\x1b[31m"
	default:
		return "\x1b[37m"
	}
}

/* -------------------- Setup -------------------- */

// Env:
//
//	IREX_LOG_LEVEL  = debug|info|warn|error
//	IREX_LOG_FORMAT = json|text|color
func SetupLogging() {
	level := parseLevel(os.Getenv("IREX_LOG_LEVEL"))
	format := strings.ToLower(os.Getenv("IREX_LOG_FORMAT"))

	var base slog.Handler
	switch format {
	case "text":
		base = NewTextHandler(os.Stdout, level)
	case "color":
		base = newColorHandler(os.Stdout, level)
	default:
		base = NewJSONHandler(os.Stdout, level)
	}

	logger := slog.New(newElapsedHandler(base))
	slog.SetDefault(logger)
}

func parseLevel(v string) slog.Level {
	switch strings.ToLower(v) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
