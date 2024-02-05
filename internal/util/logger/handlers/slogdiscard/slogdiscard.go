// DiscardLogger is for testing purposes. It's messages will be omitted when logger functionality is called in tests

package slogdiscard

import (
	"context"
	"log/slog"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

// Implements slog handler. Then it can be passed to Constructor and be treated as new slog.Logger (without messages)
type DiscardHandler struct {
}

// Empty logger handler to skip log in tests
func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	// Just ignoring log record
	return nil
}

func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	// Returns the same handler as there are not attrs for saving
	return h
}

func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	// Returns the same handler as there is no group for saving
	return h
}

func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	// Returns false as record is always ignored
	return false
}
