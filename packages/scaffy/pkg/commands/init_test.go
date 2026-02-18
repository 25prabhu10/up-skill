package commands_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/25prabhu10/scaffy/cmd/cli"
	"github.com/25prabhu10/scaffy/internal/config"
	"github.com/25prabhu10/scaffy/internal/utils"
	"github.com/25prabhu10/scaffy/pkg/build_info"
	"github.com/25prabhu10/scaffy/pkg/commands"
)

const testAppName = "scaffy"

func init() {
	build_info.AppName = testAppName
	config.DEFAULT_CONFIG_FILE_NAME = testAppName + ".json"
}

func TestNewInitCmd(t *testing.T) {
	t.Parallel()

	cmd := commands.NewInitCmd()

	if cmd == nil {
		t.Fatal("expected non-nil command")
	}

	if cmd.Use != "init" {
		t.Errorf("expected Use 'init', got %q", cmd.Use)
	}

	if cmd.RunE == nil {
		t.Error("expected non-nil RunE function")
	}
}

func TestInitCommand_CreateNew(t *testing.T) { //nolint:paralleltest // `t.Chdir` does not support parallel tests
	tests := []struct {
		name       string
		args       []string
		outputFile string
	}{
		{
			name:       "current dir",
			args:       []string{},
			outputFile: config.DEFAULT_CONFIG_FILE_NAME,
		},
		{
			name:       "custom dir",
			args:       []string{"docs"},
			outputFile: filepath.Join("docs", config.DEFAULT_CONFIG_FILE_NAME),
		},
	}

	for _, tt := range tests { //nolint:paralleltest // `t.Chdir` does not support parallel tests
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			t.Chdir(tmpDir)

			baseConfigPath := filepath.Join(tmpDir, "base.json")
			if err := config.GetDefaultConfig().Save(baseConfigPath, false); err != nil {
				t.Fatalf("failed to create base config: %v", err)
			}

			args := append([]string{"--config", baseConfigPath, "init"}, tt.args...)

			_, _, err := utils.ExecuteTestCommandWithContext(t, cli.GetRootCmd(), args, false, false)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if _, err := os.Stat(tt.outputFile); os.IsNotExist(err) {
				t.Errorf("expected config file at %s", tt.outputFile)
			}
		})
	}
}

func TestInitCommand_AlreadyExists(t *testing.T) { //nolint:paralleltest // `t.Chdir` does not support parallel tests
	tests := []struct {
		name          string
		args          []string
		force         bool
		outputFile    string
		expectedError string
	}{
		{
			name:          "without force",
			args:          []string{},
			force:         false,
			outputFile:    config.DEFAULT_CONFIG_FILE_NAME,
			expectedError: "Error: config file already exists at " + config.DEFAULT_CONFIG_FILE_NAME + " (use --force to overwrite)\n",
		},
		{
			name:          "with force",
			args:          []string{},
			force:         true,
			outputFile:    config.DEFAULT_CONFIG_FILE_NAME,
			expectedError: "",
		},
		{
			name:          "custom dir with force",
			force:         true,
			args:          []string{"docs"},
			outputFile:    filepath.Join("docs", config.DEFAULT_CONFIG_FILE_NAME),
			expectedError: "",
		},
	}

	for _, tt := range tests { //nolint:paralleltest // `t.Chdir` does not support parallel tests
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			t.Chdir(tmpDir)

			baseConfigPath := filepath.Join(tmpDir, "base.json")
			if err := config.GetDefaultConfig().Save(baseConfigPath, false); err != nil {
				t.Fatalf("failed to create base config: %v", err)
			}

			cfg := config.GetDefaultConfig()
			_ = cfg.Save(tt.outputFile, true)

			args := append([]string{"--config", baseConfigPath, "init"}, tt.args...)
			if tt.force {
				args = append(args, "--force")
			}
			_, cliBuf, err := utils.ExecuteTestCommandWithContext(t, cli.GetRootCmd(), args, false, false)

			if tt.force {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if _, err := os.Stat(tt.outputFile); os.IsNotExist(err) {
					t.Errorf("expected config file at %s", tt.outputFile)
				}
			} else if err == nil {
				t.Fatal("expected error, got none")
			}

			if !strings.Contains(cliBuf.String(), tt.expectedError) {
				t.Errorf("expected %q, got %q", tt.expectedError, cliBuf.String())
			}
		})
	}
}

func TestInitCommand_FlagPersistsToConfig(t *testing.T) { //nolint:paralleltest // `t.Chdir` does not support parallel tests
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	configDir := "old"
	author := "John"
	languages := "ruby=rb,python=py"
	logLevel := "debug"
	outputDir := "generated"
	templatesDir := "my-templates"

	args := []string{configDir, "--author", author, "--lang", languages, "--log-level", logLevel, "--output-dir", outputDir, "--templates-dir", templatesDir}
	configFilePath := filepath.Join(configDir, config.DEFAULT_CONFIG_FILE_NAME)

	baseConfigPath := filepath.Join(tmpDir, "base.json")
	if err := config.GetDefaultConfig().Save(baseConfigPath, false); err != nil {
		t.Fatalf("failed to create base config: %v", err)
	}

	execArgs := append([]string{"--config", baseConfigPath, "init"}, args...)

	_, _, err := utils.ExecuteTestCommandWithContext(t, cli.GetRootCmd(), execArgs, false, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	loadedCfg, err := config.LoadConfigFromFile(configFilePath)
	if err != nil {
		t.Fatalf("failed to load generated config: %v", err)
	}

	if loadedCfg.Author != author {
		t.Fatalf("expected author %q, got %q", author, loadedCfg.Author)
	}

	if loadedCfg.Languages["ruby"] != "rb" {
		t.Fatalf("expected ruby language extension 'rb', got %q", loadedCfg.Languages["ruby"])
	}

	if loadedCfg.Languages["python"] != "py" {
		t.Fatalf("expected python language extension 'py', got %q", loadedCfg.Languages["python"])
	}

	if loadedCfg.LogLevel != logLevel {
		t.Fatalf("expected log level %q, got %q", logLevel, loadedCfg.LogLevel)
	}

	if loadedCfg.OutputDir != outputDir {
		t.Fatalf("expected output dir %q, got %q", outputDir, loadedCfg.OutputDir)
	}

	if loadedCfg.TemplatesDir != templatesDir {
		t.Fatalf("expected templates dir %q, got %q", templatesDir, loadedCfg.TemplatesDir)
	}
}
