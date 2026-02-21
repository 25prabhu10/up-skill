package program_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/25prabhu10/scaffy/internal/config"
	"github.com/25prabhu10/scaffy/internal/program"
	"github.com/25prabhu10/scaffy/pkg/build_info"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	goosWindows = "windows"
	goosDarwin  = "darwin"
	goosIOS     = "ios"
	testAppName = "scaffy"
)

func init() { //nolint:gochecknoinits // needed to set test constants before tests run
	build_info.APP_NAME = testAppName
	config.DEFAULT_CONFIG_FILE_NAME = testAppName + ".json"
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()

			tmpHome := t.TempDir()

			switch runtime.GOOS {
			case goosWindows:
				t.Setenv("AppData", tmpHome)
				t.Setenv("USERPROFILE", tmpHome)
			case goosDarwin, goosIOS:
				t.Setenv("HOME", tmpHome)
			default:
				t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmpHome, ".config"))
			}

			tmpDir := t.TempDir()
			t.Chdir(tmpDir)

			cfg := config.GetDefaultConfig()

			p, err := program.New("", false, false)
			if err != nil {
				t.Fatalf("failed to create program: %v", err)
			}

			if tt.setupExisting {
				dir := tt.outputDir
				if dir == "" {
					dir = "."
				} else {
					if err := os.MkdirAll(dir, 0750); err != nil {
						t.Fatalf("failed to create dir: %v", err)
					}
				}

				configFilePath := filepath.Join(dir, config.DEFAULT_CONFIG_FILE_NAME)
				if err := os.WriteFile(configFilePath, []byte("{}"), 0600); err != nil {
					t.Fatalf("failed to write existing config: %v", err)
				}
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

func TestProgram_GenerateLLMDocs(t *testing.T) {
	tests := []struct {
		name        string
		format      string
		frontMatter bool
		expectError bool
	}{
		{
			name:        "markdown without frontmatter",
			format:      "markdown",
			frontMatter: false,
			expectError: false,
		},
		{
			name:        "markdown with frontmatter",
			format:      "markdown",
			frontMatter: true,
			expectError: false,
		},
		{
			name:        "man pages",
			format:      "man",
			frontMatter: false,
			expectError: false,
		},
		{
			name:        "rest docs",
			format:      "rest",
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()

			tmpHome := t.TempDir()

			switch runtime.GOOS {
			case goosWindows:
				t.Setenv("AppData", tmpHome)
				t.Setenv("USERPROFILE", tmpHome)
			case goosDarwin, goosIOS:
				t.Setenv("HOME", tmpHome)
			default:
				t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmpHome, ".config"))
			}

			tmpDir := t.TempDir()

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
