// Package logger provides structured logging functionality for the up-skill CLI.
//
// It uses the standard library's log/slog package for structured logging with
// support for different log levels and output formats. Technical logs are
// written to stderr to separate them from user-facing output on stdout.
package logger

import (
	"io"
	"log/slog"
	"os"
)

// New creates a new structured logger with the specified configuration.
// The log level is determined by the following priority:
//  1. If quiet is true, only errors are logged
//  2. If verbose is true, debug and above are logged
//  3. Otherwise, the level parameter determines the log level
//
// Valid level values are: "debug", "info", "warn", "error".
// If an invalid level is provided, "info" is used as the default.
// Logs are written to stderr to separate technical output from user messages.
func New(level string, verbose, quiet bool) *slog.Logger {
	return NewWithWriter(level, verbose, quiet, os.Stderr)
}

// NewWithWriter creates a new structured logger that writes to the specified writer.
// This is useful for testing or redirecting log output.
func NewWithWriter(level string, verbose, quiet bool, w io.Writer) *slog.Logger {
	var logLevel slog.Level

	if quiet { //nolint:gocritic // explicit is better than implicit
		logLevel = slog.LevelError
	} else if verbose {
		logLevel = slog.LevelDebug
	} else {
		switch level {
		case "debug":
			logLevel = slog.LevelDebug
		case "info":
			logLevel = slog.LevelInfo
		case "warn":
			logLevel = slog.LevelWarn
		case "error":
			logLevel = slog.LevelError
		default:
			logLevel = slog.LevelError
		}
	}

	opts := &slog.HandlerOptions{Level: logLevel}
	handler := slog.NewTextHandler(w, opts)

	return slog.New(handler)
}
