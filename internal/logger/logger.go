package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	// Default to error level
	level := slog.LevelError
	
	// Check for debug flag
	if os.Getenv("ALIGN_DEBUG") == "true" {
		level = slog.LevelDebug
	}
	
	opts := &slog.HandlerOptions{
		Level: level,
	}
	
	handler := slog.NewTextHandler(os.Stderr, opts)
	Logger = slog.New(handler)
}

// Debug returns the logger for debug output
func Debug() *slog.Logger {
	return Logger
}