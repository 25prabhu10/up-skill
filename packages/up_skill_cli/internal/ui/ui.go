// Package ui provides POSIX-compliant user-facing output for the up-skill CLI.
//
// All diagnostic and user feedback messages are written to stderr, leaving
// stdout free for program output that may be piped to other tools.
// This follows POSIX utility conventions where:
//   - stdout: program output meant for piping/processing
//   - stderr: diagnostic messages, errors, warnings, and user feedback
//
// The UI respects quiet mode (-q, --quiet) to suppress non-error output.
package ui

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
)

// contextKey is a private type for context keys to avoid collisions.
type contextKey struct{}

// uiKey is the context key for storing the UI instance.
var uiKey = contextKey{}

// UI provides POSIX-compliant user-facing output.
// All messages write to stderr for proper stream separation.
// It is safe for concurrent use.
type UI struct {
	errW  io.Writer
	quiet bool
	mu    sync.Mutex
}

// Option is a functional option for configuring UI.
type Option func(*UI)

// New creates a new UI instance with the given options.
func New(opts ...Option) *UI {
	u := &UI{
		errW:  os.Stderr,
		quiet: false,
	}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

// WithOutput sets the output writer for all UI messages (primarily for testing).
func WithOutput(w io.Writer) Option {
	return func(u *UI) {
		u.errW = w
	}
}

// WithQuiet enables quiet mode, suppressing non-error output.
func WithQuiet(quiet bool) Option {
	return func(u *UI) {
		u.quiet = quiet
	}
}

// Warningf prints a warning message to stderr.
// Warnings are shown even in quiet mode as they require user attention.
func (u *UI) Warningf(format string, args ...any) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.quiet {
		return
	}

	msg := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(u.errW, "warning: %s\n", msg)
}

// Errorf prints an error message to stderr.
// Errors are always shown regardless of quiet mode.
func (u *UI) Errorf(format string, args ...any) {
	u.mu.Lock()
	defer u.mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(u.errW, "error: %s\n", msg)
}

// Infof prints an informational message to stderr without prefix.
// Use this for general progress or status updates.
// In quiet mode, info messages are suppressed.
func (u *UI) Infof(format string, args ...any) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.quiet {
		return
	}

	msg := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintln(u.errW, msg)
}

// WithUI returns a new context with the UI attached.
func WithUI(ctx context.Context, u *UI) context.Context {
	return context.WithValue(ctx, uiKey, u)
}

// FromContext extracts the UI from the context.
// If no UI is found, it returns a default UI instance.
func FromContext(ctx context.Context) *UI {
	if u, ok := ctx.Value(uiKey).(*UI); ok {
		return u
	}

	fmt.Fprintf(os.Stderr, "warning: no UI found in context, using default UI\n")

	// Return default UI if not found in context
	return New()
}
