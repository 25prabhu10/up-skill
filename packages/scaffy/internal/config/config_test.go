package config_test

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/25prabhu10/scaffy/internal/config"
	"github.com/25prabhu10/scaffy/internal/utils"
	"github.com/25prabhu10/scaffy/pkg/build_info"
	"github.com/spf13/viper"
)

const (
	testAppName   = "scaffy"
	logLevelDebug = "debug"
	goosWindows   = "windows"
	goosDarwin    = "darwin"
	goosIOS       = "ios"
)

func init() { //nolint:gochecknoinits // needed to set test constants before tests run
	build_info.AppName = testAppName
	config.DEFAULT_CONFIG_FILE_NAME = testAppName + ".json"
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := config.GetDefaultConfig()

	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	if cfg.Author != "" {
		t.Errorf("expected Author to be empty, got %s", cfg.Author)
	}

	if cfg.OutputDir != "." {
		t.Errorf("expected OutputDir '.', got %q", cfg.OutputDir)
	}

	if cfg.TemplatesDir != "" {
		t.Errorf("expected TemplatesDir to be empty, got %s", cfg.TemplatesDir)
	}

	if len(cfg.Languages) == 0 {
		t.Fatal("expected non-empty Languages list")
	}

	expectedLangs := config.GetDefaultLanguages()
	if len(cfg.Languages) != len(expectedLangs) {
		t.Errorf("expected %d languages, got %d", len(expectedLangs), len(cfg.Languages))
	}

	for lang, fileExt := range expectedLangs {
		if cfg.Languages[lang] != fileExt {
			t.Errorf("expected language %q with %q file extension, got %q", lang, fileExt, cfg.Languages[lang])
		}
	}

	if cfg.LogLevel != strings.ToLower(slog.LevelError.String()) {
		t.Errorf("expected LogLevel '%s', got '%q'", strings.ToLower(slog.LevelError.String()), cfg.LogLevel)
	}
}

func TestGetDefaultLanguages(t *testing.T) {
	t.Parallel()

	langs := config.GetDefaultLanguages()

	if len(langs) == 0 {
		t.Fatal("expected non-empty languages map")
	}

	expectedLangs := []string{"go", "c"}
	for _, lang := range expectedLangs {
		if _, ok := langs[lang]; !ok {
			t.Errorf("expected language %q in default languages", lang)
		}
	}

	if langs["go"] != "go" {
		t.Errorf("expected go extension 'go', got %q", langs["go"])
	}

	if langs["c"] != "c" {
		t.Errorf("expected c extension 'c', got %q", langs["c"])
	}
}

func TestVerboseLogLevel(t *testing.T) {
	t.Parallel()

	expected := logLevelDebug

	got := config.VerboseLogLevel()
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestQuietLogLevel(t *testing.T) {
	t.Parallel()

	expected := "error"

	got := config.QuietLogLevel()
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestAllLogLevelsStr(t *testing.T) {
	t.Parallel()

	expected := "debug,info,warn,error"

	got := config.AllLogLevelsStr()
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestGetDefaultConfigDir(t *testing.T) {
	t.Parallel()

	dir := config.GetDefaultConfigDir()
	if utils.IsStringEmpty(dir) {
		t.Error("expected non-empty config directory")
	}
}

func TestGetDefaultConfigPath(t *testing.T) {
	t.Parallel()

	path := config.GetDefaultConfigPath()
	if utils.IsStringEmpty(path) {
		t.Error("expected non-empty config path")
	}

	switch runtime.GOOS {
	case goosWindows:
		if !strings.Contains(path, "%USERPROFILE%") {
			t.Errorf("expected path to contain %%USERPROFILE%%, got %q", path)
		}
	case goosDarwin, goosIOS:
		if !strings.Contains(path, "$HOME") {
			t.Errorf("expected path to contain $HOME, got %q", path)
		}
	default:
		if !strings.Contains(path, "$XDG_CONFIG_HOME") {
			t.Errorf("expected path to contain $XDG_CONFIG_HOME, got %q", path)
		}
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	t.Parallel()

	viper.Reset()

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test-config.json")

	cfg := config.GetDefaultConfig()
	cfg.Author = "test-author"
	cfg.OutputDir = "/custom/output"
	cfg.TemplatesDir = "/custom/templates"
	cfg.LogLevel = logLevelDebug

	if err := cfg.Save(path, false); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	loaded, err := config.LoadConfigFromFile(path)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loaded.Author != "test-author" {
		t.Errorf("expected Author 'test-author', got %q", loaded.Author)
	}

	if loaded.OutputDir != "/custom/output" {
		t.Errorf("expected OutputDir '/custom/output', got %q", loaded.OutputDir)
	}

	if loaded.TemplatesDir != "/custom/templates" {
		t.Errorf("expected TemplatesDir '/custom/templates', got %q", loaded.TemplatesDir)
	}

	if loaded.LogLevel != "debug" {
		t.Errorf("expected LogLevel 'debug', got %q", loaded.LogLevel)
	}

	if loaded.Languages["go"] != "go" {
		t.Errorf("expected Languages['go'] = 'go', got %q", loaded.Languages["go"])
	}
}

func TestSaveAndLoadConfig_AllFields(t *testing.T) { //nolint:paralleltest // uses global viper state (viper.Reset)
	tests := []struct {
		name string
		cfg  *config.Config
	}{
		{
			name: "empty_fields",
			cfg: &config.Config{
				Author:       "",
				OutputDir:    ".",
				TemplatesDir: "",
				LogLevel:     "error",
				Languages:    map[string]string{"go": "go"},
			},
		},
		{
			name: "all_fields_set",
			cfg: &config.Config{
				Author:       "test-user",
				OutputDir:    "/tmp/output",
				TemplatesDir: "/tmp/templates",
				LogLevel:     "debug",
				Languages:    map[string]string{"go": "go", "python": "py", "rust": "rs"},
			},
		},
		{
			name: "special_characters",
			cfg: &config.Config{
				Author:       "user@example.com",
				OutputDir:    "/path/with spaces/dir",
				TemplatesDir: "/path/with-underscore",
				LogLevel:     "info",
				Languages:    map[string]string{},
			},
		},
	}

	//nolint:paralleltest // uses global viper state (viper.Reset)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()

			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, "config-"+tt.name+".json")

			if err := tt.cfg.Save(path, false); err != nil {
				t.Fatalf("failed to save config: %v", err)
			}

			loaded, err := config.LoadConfigFromFile(path)
			if err != nil {
				t.Fatalf("failed to load config: %v", err)
			}

			if loaded.Author != tt.cfg.Author {
				t.Errorf("Author: expected %q, got %q", tt.cfg.Author, loaded.Author)
			}

			if loaded.OutputDir != tt.cfg.OutputDir {
				t.Errorf("OutputDir: expected %q, got %q", tt.cfg.OutputDir, loaded.OutputDir)
			}

			if loaded.TemplatesDir != tt.cfg.TemplatesDir {
				t.Errorf("TemplatesDir: expected %q, got %q", tt.cfg.TemplatesDir, loaded.TemplatesDir)
			}

			if loaded.LogLevel != tt.cfg.LogLevel {
				t.Errorf("LogLevel: expected %q, got %q", tt.cfg.LogLevel, loaded.LogLevel)
			}

			for lang, ext := range tt.cfg.Languages {
				if loaded.Languages[lang] != ext {
					t.Errorf("Languages[%q]: expected %q, got %q", lang, ext, loaded.Languages[lang])
				}
			}
		})
	}
}

func TestSaveConfig_Overwrite(t *testing.T) { //nolint:paralleltest // uses global viper state (viper.Reset)
	viper.Reset()

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "overwrite-test.json")

	if err := config.GetDefaultConfig().Save(path, false); err != nil {
		t.Fatalf("failed to save initial config: %v", err)
	}

	err := config.GetDefaultConfig().Save(path, false)
	if err == nil {
		t.Error("expected error when saving without overwrite")
	}

	cfg := config.GetDefaultConfig()

	cfg.Author = "updated-author"
	if err := cfg.Save(path, true); err != nil {
		t.Errorf("expected no error with overwrite=true, got %v", err)
	}

	viper.Reset()

	loaded, err := config.LoadConfigFromFile(path)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loaded.Author != "updated-author" {
		t.Errorf("expected updated author, got %q", loaded.Author)
	}
}

func TestSaveConfig_CreatesParentDirectory(t *testing.T) { //nolint:paralleltest // uses global viper state (viper.Reset)
	viper.Reset()

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "nested", "deep", "dir", "config.json")

	cfg := config.GetDefaultConfig()
	if err := cfg.Save(path, false); err != nil {
		t.Fatalf("failed to save config with nested dirs: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected config file to be created")
	}
}

func TestSaveConfig_EmptyPath(t *testing.T) { //nolint:paralleltest // uses global viper state (viper.Reset)
	viper.Reset()

	cfg := config.GetDefaultConfig()

	err := cfg.Save("", false)
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestLoadConfigFromFile_NotFound(t *testing.T) { //nolint:paralleltest // uses global viper state (viper.Reset)
	viper.Reset()

	_, err := config.LoadConfigFromFile("/nonexistent/path/config.json")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestLoadConfigFromFile_InvalidJSON(t *testing.T) { //nolint:paralleltest // uses global viper state (viper.Reset)
	viper.Reset()

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "invalid.json")

	invalidJSON := `{"author": "test", "languages": invalid}`
	if err := os.WriteFile(path, []byte(invalidJSON), 0600); err != nil {
		t.Fatalf("failed to write invalid JSON: %v", err)
	}

	_, err := config.LoadConfigFromFile(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoadConfigFromFile_DirectoryPath(t *testing.T) { //nolint:paralleltest // uses global viper state (viper.Reset)
	viper.Reset()

	tmpDir := t.TempDir()

	_, err := config.LoadConfigFromFile(tmpDir)
	if err == nil {
		t.Error("expected error when path is a directory")
	}
}

func TestLoadConfigFromDefaultFile_NotFound(t *testing.T) { //nolint:paralleltest // t.Setenv used
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

	_, err := config.LoadConfigFromDefaultFile()
	if err == nil {
		t.Error("expected error when no config file found")
	}
}

func TestLoadConfigFromDefaultFile_FoundInCurrentDir(t *testing.T) { //nolint:paralleltest // changes working directory
	viper.Reset()

	tmpDir := t.TempDir()

	cfg := config.GetDefaultConfig()
	cfg.Author = "current-dir-author"

	configPath := filepath.Join(tmpDir, "scaffy.json")
	if err := cfg.Save(configPath, false); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	t.Chdir(tmpDir)

	viper.Reset()

	loaded, err := config.LoadConfigFromDefaultFile()
	if err != nil {
		t.Fatalf("failed to load config from default: %v", err)
	}

	if loaded.Author != "current-dir-author" {
		t.Errorf("expected Author 'current-dir-author', got %q", loaded.Author)
	}
}

func TestEnsureDefaultConfig(t *testing.T) { //nolint:paralleltest // t.Setenv used
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

	cfg, err := config.EnsureDefaultConfig()
	if err != nil {
		t.Fatalf("failed to ensure default config: %v", err)
	}

	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	configPath := filepath.Join(config.GetDefaultConfigDir(), "scaffy.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("expected config file at %s", configPath)
	}
}

func TestEnsureDefaultConfig_AlreadyExists(t *testing.T) { //nolint:paralleltest // t.Setenv used, uses global viper state (viper.Reset)
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

	cfg, err := config.EnsureDefaultConfig()
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}

	_, err = config.EnsureDefaultConfig()
	if err == nil {
		t.Error("expected error when config already exists")
	}

	if cfg == nil {
		t.Error("expected config from first call")
	}
}

func TestEnvVariableOverride(t *testing.T) {
	viper.Reset()

	tmpDir := t.TempDir()

	t.Setenv("SCAFFY_AUTHOR", "env-author")
	t.Setenv("SCAFFY_LOG_LEVEL", "warn")

	cfg := config.GetDefaultConfig()
	cfg.Author = "file-author"
	cfg.LogLevel = logLevelDebug

	path := filepath.Join(tmpDir, "config.json")
	if err := cfg.Save(path, false); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	loaded, err := config.LoadConfigFromFile(path)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loaded.Author != "env-author" {
		t.Errorf("expected env override Author 'env-author', got %q", loaded.Author)
	}

	if loaded.LogLevel != "warn" {
		t.Errorf("expected env override LogLevel 'warn', got %q", loaded.LogLevel)
	}
}

func TestEnvVariableOverride_OutputDir(t *testing.T) {
	viper.Reset()

	tmpDir := t.TempDir()

	t.Setenv("SCAFFY_OUTPUT_DIR", "/env/output")
	t.Setenv("SCAFFY_TEMPLATES_DIR", "/env/templates")

	cfg := config.GetDefaultConfig()
	cfg.OutputDir = "/file/output"
	cfg.TemplatesDir = "/file/templates"

	path := filepath.Join(tmpDir, "config.json")
	if err := cfg.Save(path, false); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	loaded, err := config.LoadConfigFromFile(path)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loaded.OutputDir != "/env/output" {
		t.Errorf("expected env override OutputDir '/env/output', got %q", loaded.OutputDir)
	}

	if loaded.TemplatesDir != "/env/templates" {
		t.Errorf("expected env override TemplatesDir '/env/templates', got %q", loaded.TemplatesDir)
	}
}

func TestLoadConfigFromFile_InvalidLogLevel(t *testing.T) { //nolint:paralleltest // uses global viper state (viper.Reset)
	viper.Reset()

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "invalid-log-level.json")

	cfg := config.GetDefaultConfig()
	cfg.LogLevel = "trace"

	err := cfg.Save(path, false)
	if err == nil {
		t.Fatal("expected save to fail for invalid log level")
	}
}

func TestConfigJSONFormat(t *testing.T) { //nolint:paralleltest // uses global viper state (viper.Reset)
	viper.Reset()

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "format-test.json")

	cfg := config.GetDefaultConfig()
	cfg.Author = "json-test"
	cfg.Languages = map[string]string{"go": "go", "rust": "rs"}

	if err := cfg.Save(path, false); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	data, err := os.ReadFile(path) //nolint:gosec // test needs to read the file it just wrote
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	if !json.Valid(data) {
		t.Error("config file is not valid JSON")
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if raw["author"] != "json-test" {
		t.Errorf("expected author 'json-test' in JSON, got %v", raw["author"])
	}
}
