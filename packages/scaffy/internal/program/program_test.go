package program_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/25prabhu10/scaffy/internal/config"
	"github.com/25prabhu10/scaffy/internal/program"
	"github.com/25prabhu10/scaffy/internal/utils"
	"github.com/25prabhu10/scaffy/internal/utils/test_utils"
	"github.com/25prabhu10/scaffy/pkg/build_info"
	"github.com/spf13/cobra"
)

const (
	testAppName = "scaffy"
)

func init() { //nolint:gochecknoinits // needed to set test constants before tests run
	build_info.APP_NAME = testAppName
	config.DEFAULT_CONFIG_FILE_NAME = testAppName + ".json"
}

func TestNew(t *testing.T) { //nolint:paralleltest // t.Setenv used
	tests := []struct {
		name        string
		setupFunc   func(t *testing.T, tmpDir string) string
		logLevel    string
		verbose     bool
		quiet       bool
		expectError bool
	}{
		{
			name: "success with default config creation",
			setupFunc: func(t *testing.T, tmpDir string) string {
				t.Helper()
				return ""
			},
			logLevel:    config.GetDefaultLogLevel(),
			verbose:     false,
			quiet:       false,
			expectError: false,
		},
		{
			name: "success with existing config file",
			setupFunc: func(t *testing.T, tmpDir string) string {
				t.Helper()

				cfgPath := filepath.Join(tmpDir, "test-config.json")
				cfg := config.GetDefaultConfig()

				if err := cfg.Save(cfgPath, false, utils.NewFileSystem()); err != nil {
					t.Fatalf("failed to save config: %v", err)
				}

				return cfgPath
			},
			logLevel:    config.GetDefaultLogLevel(),
			verbose:     true,
			quiet:       false,
			expectError: false,
		},
		{
			name: "fails with invalid config file",
			setupFunc: func(t *testing.T, tmpDir string) string {
				t.Helper()

				cfgPath := filepath.Join(tmpDir, "invalid-config.json")

				if err := os.WriteFile(cfgPath, []byte("{invalid json}"), 0600); err != nil {
					t.Fatalf("failed to write invalid config: %v", err)
				}

				return cfgPath
			},
			logLevel:    config.GetDefaultLogLevel(),
			verbose:     false,
			quiet:       true,
			expectError: true,
		},
		{
			name: "fails with non-existent explicit config file",
			setupFunc: func(t *testing.T, tmpDir string) string {
				t.Helper()
				return filepath.Join(tmpDir, "non-existent.json")
			},
			logLevel:    config.GetDefaultLogLevel(),
			verbose:     false,
			quiet:       false,
			expectError: true,
		},
	}

	for _, tt := range tests { //nolint:paralleltest // t.Setenv used
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := test_utils.SetupTestEnv(t)

			configFile := tt.setupFunc(t, tmpDir)

			p, err := program.New(configFile, tt.verbose, tt.quiet)

			if tt.expectError { //nolint:nestif // clearer to read this way in tests
				if err == nil {
					t.Errorf("expected error, got nil")
				}

				if p != nil {
					t.Errorf("expected nil program on error, got %v", p)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if p == nil {
					t.Errorf("expected program, got nil")
				}
			}
		})
	}
}

func setupExistingConfig(t *testing.T, tmpDir string, outputDir string) {
	t.Helper()

	cfg := config.GetDefaultConfig()
	dir := outputDir

	if dir == "" {
		dir = tmpDir
	} else {
		dir = filepath.Join(tmpDir, dir)
		if err := os.MkdirAll(dir, 0750); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
	}

	configFilePath := filepath.Join(dir, config.DEFAULT_CONFIG_FILE_NAME)
	if err := cfg.Save(configFilePath, false, utils.NewFileSystem()); err != nil {
		t.Fatalf("failed to write existing config: %v", err)
	}
}

func TestProgram_InitializeConfig(t *testing.T) { //nolint:paralleltest // t.Setenv used
	tests := []struct {
		name          string
		outputDir     string
		force         bool
		setupExisting bool
		expectError   bool
	}{
		{
			name:          "success without output dir",
			outputDir:     "",
			force:         false,
			setupExisting: false,
			expectError:   false,
		},
		{
			name:          "success with output dir",
			outputDir:     "testdir",
			force:         false,
			setupExisting: false,
			expectError:   false,
		},
		{
			name:          "fails when config exists and not forced",
			outputDir:     "",
			force:         false,
			setupExisting: true,
			expectError:   true,
		},
		{
			name:          "success when config exists and forced",
			outputDir:     "",
			force:         true,
			setupExisting: true,
			expectError:   false,
		},
	}

	for _, tt := range tests { //nolint:paralleltest // t.Setenv used
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := test_utils.SetupTestEnv(t)

			cfg := config.GetDefaultConfig()

			p, err := program.New("", false, false)
			if err != nil {
				t.Fatalf("failed to create program: %v", err)
			}

			if tt.setupExisting {
				setupExistingConfig(t, tmpDir, tt.outputDir)
			}

			_, err = p.InitializeConfig(cfg, tt.outputDir, tt.force)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestProgram_GenerateLLMDocs(t *testing.T) { //nolint:paralleltest // t.Setenv used
	tests := []struct {
		name        string
		format      string
		frontMatter bool
		expectError bool
	}{
		{
			name:        "markdown without frontmatter",
			format:      program.Markdown,
			frontMatter: false,
			expectError: false,
		},
		{
			name:        "markdown with frontmatter",
			format:      program.Markdown,
			frontMatter: true,
			expectError: false,
		},
		{
			name:        "man pages",
			format:      program.Man,
			frontMatter: false,
			expectError: false,
		},
		{
			name:        "rest docs",
			format:      program.Rest,
			frontMatter: false,
			expectError: false,
		},
		{
			name:        "invalid format",
			format:      "invalid",
			frontMatter: false,
			expectError: true,
		},
	}

	for _, tt := range tests { //nolint:paralleltest // t.Setenv used
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := test_utils.SetupTestEnv(t)

			p, err := program.New("", false, false)
			if err != nil {
				t.Fatalf("failed to create program: %v", err)
			}

			cmd := &cobra.Command{Use: "testcmd"}

			err = p.GenerateLLMDocs(cmd, tmpDir, tt.format, tt.frontMatter)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestContext(t *testing.T) { //nolint:paralleltest // t.Setenv used
	t.Run("WithProgram and FromContext", func(t *testing.T) { //nolint:paralleltest // t.Setenv used
		test_utils.SetupTestEnv(t)

		p, err := program.New("", false, false)
		if err != nil {
			t.Fatalf("failed to create program: %v", err)
		}

		ctx := context.Background()
		ctx = program.WithProgram(ctx, p)

		retrieved := program.FromContext(ctx)
		if retrieved != p {
			t.Errorf("expected retrieved program to be the same instance")
		}
	})

	t.Run("FromContext without WithProgram", func(t *testing.T) { //nolint:paralleltest // t.Setenv used
		test_utils.SetupTestEnv(t)

		ctx := context.Background()
		retrieved := program.FromContext(ctx)

		if retrieved == nil {
			t.Errorf("expected default program, got nil")
		}
	})
}
