package slg

import (
	"log/slog"
)

// Err(error) slog.Attr is for wrapping errors in slog logger
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
