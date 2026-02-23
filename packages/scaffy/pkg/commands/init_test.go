package commands_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/25prabhu10/scaffy/internal/config"
	"github.com/25prabhu10/scaffy/internal/utils/test_utils"
	"github.com/spf13/viper"

	"github.com/25prabhu10/scaffy/pkg/build_info"
	"github.com/25prabhu10/scaffy/pkg/commands"
)

const (
	testAppName = "scaffy"
)

func init() { //nolint:gochecknoinits // needed to set test constants before tests run
	build_info.APP_NAME = testAppName
	config.DEFAULT_CONFIG_FILE_NAME = testAppName + ".json"
}

func TestInitCmd_Success(t *testing.T) { //nolint:paralleltest // t.Setenv used and global force flag
	tests := []struct {
		name           string
		args           []string
		expectedErrOut string
		verifyConfig   func(t *testing.T, configPath string)
	}{
		{
			name:           "success without args",
			args:           []string{},
			expectedErrOut: "initialized scaffy with default config at scaffy.json",
			verifyConfig: func(t *testing.T, configPath string) {
				t.Helper()

				if _, err := os.Stat(configPath); err != nil {
					t.Errorf("expected config file at %s, got error: %v", configPath, err)
				}
			},
		},
		{
			name:           "success with output dir arg",
			args:           []string{"custom_dir"},
			expectedErrOut: "initialized scaffy with default config at custom_dir/scaffy.json",
			verifyConfig: func(t *testing.T, configPath string) {
				t.Helper()

				expectedPath := filepath.Join(filepath.Dir(configPath), "custom_dir", "scaffy.json")
				if _, err := os.Stat(expectedPath); err != nil {
					t.Errorf("expected config file at %s, got error: %v", expectedPath, err)
				}
			},
		},
		{
			name:           "success with flags",
			args:           []string{"--author", "Test Author", "--log-level", "debug", "--output-dir", "out"},
			expectedErrOut: "initialized scaffy with default config at scaffy.json",
			verifyConfig: func(t *testing.T, configPath string) {
				t.Helper()

				cfg, err := config.LoadConfigFromFile(viper.New(), configPath)
				if err != nil {
					t.Fatalf("failed to load config: %v", err)
				}

				if cfg.Author != "Test Author" {
					t.Errorf("expected author 'Test Author', got '%s'", cfg.Author)
				}

				if cfg.LogLevel != "debug" {
					t.Errorf("expected log level 'debug', got '%s'", cfg.LogLevel)
				}

				if cfg.OutputDir != "out" {
					t.Errorf("expected output dir 'out', got '%s'", cfg.OutputDir)
				}
			},
		},
	}

	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := test_utils.SetupTestEnv(t)
			cmd := commands.GetInitCmd()

			uiStderr, _, err := test_utils.ExecuteTestCommandWithContext(t, cmd, tt.args, false, false)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			errOut := uiStderr.String()
			if tt.expectedErrOut != "" && !strings.Contains(errOut, tt.expectedErrOut) {
				t.Errorf("expected stderr to contain %q, got %q", tt.expectedErrOut, errOut)
			}

			if tt.verifyConfig != nil {
				tt.verifyConfig(t, filepath.Join(tmpDir, config.DEFAULT_CONFIG_FILE_NAME))
			}
		})
	}
}

func TestInitCmd_ExistingConfig(t *testing.T) { //nolint:paralleltest // t.Setenv used and global force flag
	tests := []struct {
		name           string
		args           []string
		expectError    bool
		expectedErrOut string
	}{
		{
			name:           "fails when config exists and not forced",
			args:           []string{},
			expectError:    true,
			expectedErrOut: "config file already exists",
		},
		{
			name:           "success when config exists and forced",
			args:           []string{"--force"},
			expectError:    false,
			expectedErrOut: "existing config will be overwritten (--force)",
		},
	}

	for _, tt := range tests { //nolint:paralleltest // t.Setenv used and global force flag
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := test_utils.SetupTestEnv(t)

			// Setup existing config
			configPath := filepath.Join(tmpDir, config.DEFAULT_CONFIG_FILE_NAME)
			if err := os.WriteFile(configPath, []byte("{}"), 0600); err != nil {
				t.Fatalf("failed to create existing config: %v", err)
			}

			cmd := commands.GetInitCmd()
			uiStderr, _, err := test_utils.ExecuteTestCommandWithContext(t, cmd, tt.args, false, false)

			if tt.expectError { //nolint:nestif // clearer to separate error vs no error cases here
				if err == nil {
					t.Error("expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.expectedErrOut) {
					t.Errorf("expected error to contain %q, got %q", tt.expectedErrOut, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				errOut := uiStderr.String()
				if !strings.Contains(errOut, tt.expectedErrOut) {
					t.Errorf("expected stderr to contain %q, got %q", tt.expectedErrOut, errOut)
				}
			}
		})
	}
}

func TestInitCmd_Errors(t *testing.T) { //nolint:paralleltest // t.Setenv used and global force flag
	tests := []struct {
		name        string
		args        []string
		setupFile   string
		expectedErr string
	}{
		{
			name:        "fails when output dir is a file",
			args:        []string{"file_as_dir"},
			setupFile:   "file_as_dir",
			expectedErr: "failed to create output directory",
		},
	}

	for _, tt := range tests { //nolint:paralleltest // t.Setenv used and global force flag
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := test_utils.SetupTestEnv(t)

			if tt.setupFile != "" {
				filePath := filepath.Join(tmpDir, tt.setupFile)
				if err := os.WriteFile(filePath, []byte("dummy"), 0600); err != nil {
					t.Fatalf("failed to create setup file: %v", err)
				}
			}

			cmd := commands.GetInitCmd()

			_, _, err := test_utils.ExecuteTestCommandWithContext(t, cmd, tt.args, false, false)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("expected error to contain %q, got %q", tt.expectedErr, err.Error())
			}
		})
	}
}
