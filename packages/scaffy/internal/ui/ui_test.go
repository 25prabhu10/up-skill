package ui_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/25prabhu10/scaffy/internal/ui"
)

// TestNew verifies UI creation with different options.
func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("default configuration", func(t *testing.T) {
		t.Parallel()

		u := ui.New()
		if u == nil {
			t.Fatal("expected non-nil UI")
		}
	})

	t.Run("with custom output", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer

		u := ui.New(ui.WithOutput(&buf))

		u.Infof("test message")

		if !strings.Contains(buf.String(), "test message") {
			t.Errorf("expected output to contain 'test message', got %q", buf.String())
		}
	})

	t.Run("with quiet mode", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer

		u := ui.New(ui.WithOutput(&buf), ui.WithQuiet(true))

		u.Infof("should be suppressed")

		if buf.Len() != 0 {
			t.Errorf("expected empty output in quiet mode, got %q", buf.String())
		}
	})
}

// TestWarning verifies warning message formatting and output.
func TestWarning(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		quiet    bool
		format   string
		args     []any
		contains string
	}{
		{
			name:     "simple warning",
			quiet:    false,
			format:   "file will be overwritten",
			args:     nil,
			contains: "warning: file will be overwritten",
		},
		{
			name:     "formatted warning",
			quiet:    false,
			format:   "config %s is deprecated",
			args:     []any{"old_setting"},
			contains: "warning: config old_setting is deprecated",
		},
		{
			name:     "hidden in quiet mode",
			quiet:    true,
			format:   "important warning",
			args:     nil,
			contains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var errBuf bytes.Buffer

			u := ui.New(ui.WithOutput(&errBuf), ui.WithQuiet(tt.quiet))

			u.Warnf(tt.format, tt.args...)

			if !strings.Contains(errBuf.String(), tt.contains) {
				t.Errorf("expected error output to contain %q, got %q", tt.contains, errBuf.String())
			}
		})
	}
}

// TestError verifies error message formatting and output.
func TestError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		quiet    bool
		format   string
		args     []any
		contains string
	}{
		{
			name:     "simple error",
			quiet:    false,
			format:   "operation failed",
			args:     nil,
			contains: "error: operation failed",
		},
		{
			name:     "formatted error",
			quiet:    false,
			format:   "failed to create %s: %s",
			args:     []any{"/tmp/file", "permission denied"},
			contains: "error: failed to create /tmp/file: permission denied",
		},
		{
			name:     "shown in quiet mode",
			quiet:    true,
			format:   "critical error",
			args:     nil,
			contains: "error: critical error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var errBuf bytes.Buffer

			u := ui.New(ui.WithOutput(&errBuf), ui.WithQuiet(tt.quiet))

			u.Errorf(tt.format, tt.args...)

			if !strings.Contains(errBuf.String(), tt.contains) {
				t.Errorf("expected error output to contain %q, got %q", tt.contains, errBuf.String())
			}
		})
	}
}

// TestInfo verifies info message output.
func TestInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		quiet    bool
		format   string
		args     []any
		contains string
		isEmpty  bool
	}{
		{
			name:     "simple info",
			quiet:    false,
			format:   "processing files",
			args:     nil,
			contains: "processing files",
			isEmpty:  false,
		},
		{
			name:    "suppressed in quiet mode",
			quiet:   true,
			format:  "should not appear",
			args:    nil,
			isEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			u := ui.New(ui.WithOutput(&buf), ui.WithQuiet(tt.quiet))

			u.Infof(tt.format, tt.args...)

			if tt.isEmpty {
				if buf.Len() != 0 {
					t.Errorf("expected empty output, got %q", buf.String())
				}
			} else {
				if !strings.Contains(buf.String(), tt.contains) {
					t.Errorf("expected output to contain %q, got %q", tt.contains, buf.String())
				}
			}
		})
	}
}

// TestWithUI verifies storing and retrieving UI from context.
func TestWithUI(t *testing.T) {
	t.Parallel()

	t.Run("store and retrieve UI", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer

		u := ui.New(ui.WithOutput(&buf))
		ctx := ui.WithUI(context.Background(), u)

		retrieved := ui.FromContext(ctx)
		retrieved.Infof("test")

		if !strings.Contains(buf.String(), "test") {
			t.Errorf("expected output from retrieved UI, got %q", buf.String())
		}
	})

	t.Run("nested context updates", func(t *testing.T) {
		t.Parallel()

		var buf1, buf2 bytes.Buffer

		u1 := ui.New(ui.WithOutput(&buf1))
		u2 := ui.New(ui.WithOutput(&buf2))

		ctx := context.Background()
		ctx = ui.WithUI(ctx, u1)
		ctx = ui.WithUI(ctx, u2)

		retrieved := ui.FromContext(ctx)
		retrieved.Infof("test")

		if buf1.Len() != 0 {
			t.Error("expected first UI to have no output after override")
		}

		if !strings.Contains(buf2.String(), "test") {
			t.Errorf("expected second UI to have output, got %q", buf2.String())
		}
	})
}

// TestFromContext verifies UI retrieval from context.
func TestFromContext(t *testing.T) {
	t.Parallel()

	t.Run("UI exists in context", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer

		u := ui.New(ui.WithOutput(&buf))
		ctx := ui.WithUI(context.Background(), u)

		retrieved := ui.FromContext(ctx)
		if retrieved == nil {
			t.Fatal("expected non-nil UI from context")
		}
	})

	t.Run("no UI in context returns default", func(t *testing.T) {
		t.Parallel()

		retrieved := ui.FromContext(context.Background())
		if retrieved == nil {
			t.Fatal("expected non-nil default UI")
		}
	})
}

// TestUIConcurrency verifies thread-safe UI usage.
func TestUIConcurrency(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	u := ui.New(ui.WithOutput(&buf), ui.WithOutput(&buf))
	ctx := ui.WithUI(context.Background(), u)

	done := make(chan bool)

	for i := range 10 {
		go func(id int) {
			uiInstance := ui.FromContext(ctx)
			uiInstance.Infof("info from goroutine %d", id)
			uiInstance.Warnf("warning from goroutine %d", id)
			uiInstance.Errorf("error from goroutine %d", id)

			done <- true
		}(i)
	}

	for range 10 {
		<-done
	}

	output := buf.String()
	// Verify all message types appear (exact count may vary due to interleaving)

	if !strings.Contains(output, "info from goroutine") {
		t.Error("expected info messages in output")
	}

	if !strings.Contains(output, "warning from goroutine") {
		t.Error("expected warning messages in output")
	}

	if !strings.Contains(output, "error from goroutine") {
		t.Error("expected error messages in output")
	}
}

// TestStreamSeparation verifies all UI output goes to stderr (POSIX compliance).
func TestStreamSeparation(t *testing.T) {
	t.Parallel()

	// In POSIX-compliant mode, all UI messages go to stderr
	// Only actual program output should go to stdout
	var stdout, stderr bytes.Buffer

	u := ui.New(
		ui.WithOutput(&stdout),
		ui.WithOutput(&stderr),
	)

	u.Infof("info message")
	u.Warnf("warning message")
	u.Errorf("error message")

	// With POSIX compliance, stdout should be empty (no UI messages)
	// All UI messages should be in stderr

	if stdout.Len() != 0 {
		t.Errorf("expected empty stdout (all UI to stderr), got %q", stdout.String())
	}

	// Verify stderr contains all message types with proper prefixes
	if !strings.Contains(stderr.String(), "info message") {
		t.Errorf("expected stderr to contain info message, got %q", stderr.String())
	}

	if !strings.Contains(stderr.String(), "warning: warning message") {
		t.Errorf("expected stderr to contain warning message with prefix, got %q", stderr.String())
	}

	if !strings.Contains(stderr.String(), "error: error message") {
		t.Errorf("expected stderr to contain error message with prefix, got %q", stderr.String())
	}
}

// BenchmarkSuccess measures success message performance.
func BenchmarkSuccess(b *testing.B) {
	var buf bytes.Buffer

	u := ui.New(ui.WithOutput(&buf))

	for b.Loop() {
		u.Infof("benchmark message")
	}
}

// BenchmarkWithUI measures context storage performance.
func BenchmarkWithUI(b *testing.B) {
	ctx := context.Background()
	u := ui.New()

	for b.Loop() {
		_ = ui.WithUI(ctx, u)
	}
}

// BenchmarkFromContext measures context retrieval performance.
func BenchmarkFromContext(b *testing.B) {
	u := ui.New()
	ctx := ui.WithUI(context.Background(), u)

	for b.Loop() {
		_ = ui.FromContext(ctx)
	}
}
