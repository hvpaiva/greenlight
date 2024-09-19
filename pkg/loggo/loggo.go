package loggo

import (
	"context"
	"io"
	"log/slog"
)

type LogHandler struct {
	*slog.TextHandler
}

func NewLogHandler(w io.Writer, opts *slog.HandlerOptions) *LogHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &LogHandler{
		slog.NewTextHandler(w, opts),
	}
}

func (h LogHandler) Enabled(c context.Context, level slog.Level) bool {
	return h.Enabled(c, level)
}

func (h LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.WithAttrs(attrs)
}

func (h LogHandler) WithGroup(name string) slog.Handler {
	return h.WithGroup(name)
}

func (h LogHandler) Handle(c context.Context, r slog.Record) error {
	return h.Handle(c, r)
}
