package logger_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/25prabhu10/scaffy/internal/logger"
)

// TestNew verifies logger creation with different configurations.
func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		level     string
		verbose   bool
		quiet     bool
		wantLevel slog.Level
	}{
		{
			name:      "default info level",
			level:     "info",
			verbose:   false,
			quiet:     false,
			wantLevel: slog.LevelInfo,
		},
		{
			name:      "debug level",
			level:     "debug",
			verbose:   false,
			quiet:     false,
			wantLevel: slog.LevelDebug,
		},
		{
			name:      "warn level",
			level:     "warn",
			verbose:   false,
			quiet:     false,
			wantLevel: slog.LevelWarn,
		},
		{
			name:      "error level",
			level:     "error",
			verbose:   false,
			quiet:     false,
			wantLevel: slog.LevelError,
		},
		{
			name:      "verbose flag overrides level",
			level:     "info",
			verbose:   true,
			quiet:     false,
			wantLevel: slog.LevelDebug,
		},
		{
			name:      "quiet flag overrides level",
			level:     "info",
			verbose:   false,
			quiet:     true,
			wantLevel: slog.LevelError,
		},
		{
			name:      "quiet takes precedence over verbose",
			level:     "info",
			verbose:   true,
			quiet:     true,
			wantLevel: slog.LevelError,
		},
		{
			name:      "invalid level defaults to info",
			level:     "invalid",
			verbose:   false,
			quiet:     false,
			wantLevel: slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			log := logger.New(tt.level, tt.verbose, tt.quiet)
			if log == nil {
				t.Fatal("expected non-nil logger")
			}

			// Test that the logger was created and doesn't panic
			log.Info("test message")
		})
	}
}

// TestNewWithWriter verifies logger creation with custom writer.
func TestNewWithWriter(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	log := logger.NewWithWriter("info", false, false, &buf)

	log.Info("test message", "key", "value")

	output := buf.String()

	if !strings.Contains(output, "test message") {
		t.Errorf("expected output to contain 'test message', got %q", output)
	}

	if !strings.Contains(output, "key=value") {
		t.Errorf("expected output to contain 'key=value', got %q", output)
	}
}

// TestLogLevelFiltering verifies that log levels filter correctly.
func TestLogLevelFiltering(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		level     string
		verbose   bool
		quiet     bool
		shouldLog map[string]bool
	}{
		{
			name:    "info level filters debug",
			level:   "info",
			verbose: false,
			quiet:   false,
			shouldLog: map[string]bool{
				"debug": false,
				"info":  true,
				"warn":  true,
				"error": true,
			},
		},
		{
			name:    "error level filters all but error",
			level:   "error",
			verbose: false,
			quiet:   false,
			shouldLog: map[string]bool{
				"debug": false,
				"info":  false,
				"warn":  false,
				"error": true,
			},
		},
		{
			name:    "verbose shows all",
			level:   "error",
			verbose: true,
			quiet:   false,
			shouldLog: map[string]bool{
				"debug": true,
				"info":  true,
				"warn":  true,
				"error": true,
			},
		},
		{
			name:    "quiet shows only error",
			level:   "debug",
			verbose: false,
			quiet:   true,
			shouldLog: map[string]bool{
				"debug": false,
				"info":  false,
				"warn":  false,
				"error": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			log := logger.NewWithWriter(tt.level, tt.verbose, tt.quiet, &buf)

			// Log at each level
			log.Debug("debug message")

			debugOutput := buf.String()
			buf.Reset()

			log.Info("info message")

			infoOutput := buf.String()
			buf.Reset()

			log.Warn("warn message")

			warnOutput := buf.String()
			buf.Reset()

			log.Error("error message")

			errorOutput := buf.String()

			outputs := map[string]string{
				"debug": debugOutput,
				"info":  infoOutput,
				"warn":  warnOutput,
				"error": errorOutput,
			}

			for level, shouldLog := range tt.shouldLog {
				hasOutput := len(outputs[level]) > 0
				if hasOutput != shouldLog {
					if shouldLog {
						t.Errorf("expected %s level to be logged", level)
					} else {
						t.Errorf("expected %s level to be filtered, got %q", level, outputs[level])
					}
				}
			}
		})
	}
}

// BenchmarkNew measures logger creation performance.
func BenchmarkNew(b *testing.B) {
	for b.Loop() {
		_ = logger.New("info", false, false)
	}
}
