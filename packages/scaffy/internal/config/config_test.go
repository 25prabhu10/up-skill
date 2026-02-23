package config_test

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/25prabhu10/scaffy/internal/config"
	"github.com/25prabhu10/scaffy/internal/constants"
	"github.com/25prabhu10/scaffy/internal/utils"
	"github.com/25prabhu10/scaffy/internal/utils/test_utils"
	"github.com/25prabhu10/scaffy/pkg/build_info"
	"github.com/spf13/viper"
)

const (
	testAppName   = "scaffy"
	logLevelDebug = "debug"
)

func init() { //nolint:gochecknoinits // needed to set test constants before tests run
	build_info.APP_NAME = testAppName
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

func TestConstantsAndDefaults(t *testing.T) {
	t.Parallel()

	t.Run("GetDefaultLanguages", func(t *testing.T) {
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
	})

	t.Run("GetDefaultLogLevel", func(t *testing.T) {
		t.Parallel()

		expectedLogLevel := "error"

		if got := config.GetDefaultLogLevel(); got != expectedLogLevel {
			t.Errorf("expected %q, got %q", expectedLogLevel, got)
		}
	})

	t.Run("AllLogLevelsStr", func(t *testing.T) {
		t.Parallel()

		expected := "debug|info|warn|error"
		if got := config.AllLogLevelsStr(); got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})
}

func TestGetDefaultConfigDir(t *testing.T) {
	t.Parallel()

	home, _ := os.UserHomeDir()
	expectedValue := filepath.Join(home, ".config")

	tests := []struct {
		name           string
		os             string
		userConfigDir  string
		expectedSubstr string
	}{
		{
			name:           "windows with user config dir",
			os:             test_utils.GoosWindows,
			userConfigDir:  "C:\\Users\\TestUser\\AppData\\Roaming",
			expectedSubstr: "C:\\Users\\TestUser\\AppData\\Roaming",
		},
		{
			name:           "darwin with user config dir",
			os:             test_utils.GoosDarwin,
			userConfigDir:  "/Users/testuser/Library/Application Support",
			expectedSubstr: "/Users/testuser/Library/Application Support",
		},
		{
			name:           "linux with user config dir",
			os:             test_utils.GoosLinux,
			userConfigDir:  "/home/testuser/.config",
			expectedSubstr: "/home/testuser/.config",
		},
		{
			name:           "windows without user config dir",
			os:             test_utils.GoosWindows,
			userConfigDir:  "",
			expectedSubstr: expectedValue,
		},
		{
			name:           "darwin without user config dir",
			os:             test_utils.GoosDarwin,
			userConfigDir:  "",
			expectedSubstr: expectedValue,
		},
		{
			name:           "linux without user config dir",
			os:             test_utils.GoosLinux,
			userConfigDir:  "",
			expectedSubstr: expectedValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dir := config.GetDefaultConfigDir(&test_utils.MockOSInfo{MockGOOS: tt.os, MockUserConfigDir: tt.userConfigDir})
			if utils.IsStringEmpty(dir) {
				t.Error("expected non-empty config directory")
			}

			if !strings.Contains(dir, tt.expectedSubstr) {
				t.Errorf("expected config directory to contain %q, got %q", tt.expectedSubstr, dir)
			}
		})
	}
}

func TestGetDefaultConfigPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		os             string
		expectedSubstr string
	}{
		{
			name:           "windows",
			os:             test_utils.GoosWindows,
			expectedSubstr: "%USERPROFILE%",
		},
		{
			name:           "darwin",
			os:             test_utils.GoosDarwin,
			expectedSubstr: "$HOME",
		},
		{
			name:           "linux",
			os:             test_utils.GoosLinux,
			expectedSubstr: "$XDG_CONFIG_HOME",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := config.GetDefaultConfigPath(&test_utils.MockOSInfo{MockGOOS: tt.os})
			if utils.IsStringEmpty(path) {
				t.Error("expected non-empty config path")
			}

			if !strings.Contains(path, tt.expectedSubstr) {
				t.Errorf("expected path to contain %q, got %q", tt.expectedSubstr, path)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	t.Parallel()

	longString := strings.Repeat("a", constants.MAX_NAME_LENGTH+1)

	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr error
	}{
		{
			name:    "nil config",
			cfg:     nil,
			wantErr: config.ErrNilConfig,
		},
		{
			name:    "valid default config",
			cfg:     config.GetDefaultConfig(),
			wantErr: nil,
		},
		{
			name: "empty author",
			cfg: &config.Config{
				Author:    "   ",
				Languages: map[string]string{"go": "go"},
				LogLevel:  "info",
				OutputDir: ".",
			},
			wantErr: config.ErrAuthorEmpty,
		},
		{
			name: "author too long",
			cfg: &config.Config{
				Author:    longString,
				Languages: map[string]string{"go": "go"},
				LogLevel:  "info",
				OutputDir: ".",
			},
			wantErr: config.ErrAuthorInvalidLength,
		},
		{
			name: "invalid log level",
			cfg: &config.Config{
				Languages: map[string]string{"go": "go"},
				LogLevel:  "invalid",
				OutputDir: ".",
			},
			wantErr: config.ErrInvalidLogLevel,
		},
		{
			name: "empty languages",
			cfg: &config.Config{
				Languages: map[string]string{},
				LogLevel:  "info",
				OutputDir: ".",
			},
			wantErr: config.ErrLanguagesEmpty,
		},
		{
			name: "empty language key",
			cfg: &config.Config{
				Languages: map[string]string{"   ": "go"},
				LogLevel:  "info",
				OutputDir: ".",
			},
			wantErr: config.ErrLanguageEmpty,
		},
		{
			name: "language key too long",
			cfg: &config.Config{
				Languages: map[string]string{longString: "go"},
				LogLevel:  "info",
				OutputDir: ".",
			},
			wantErr: config.ErrLanguageInvalidLength,
		},
		{
			name: "empty file extension",
			cfg: &config.Config{
				Languages: map[string]string{"go": "   "},
				LogLevel:  "info",
				OutputDir: ".",
			},
			wantErr: config.ErrFileExtensionEmpty,
		},
		{
			name: "file extension too long",
			cfg: &config.Config{
				Languages: map[string]string{"go": longString},
				LogLevel:  "info",
				OutputDir: ".",
			},
			wantErr: config.ErrFileExtensionInvalidLength,
		},
		{
			name: "empty output dir",
			cfg: &config.Config{
				Languages: map[string]string{"go": "go"},
				LogLevel:  "info",
				OutputDir: "   ",
			},
			wantErr: config.ErrOutputDirEmpty,
		},
		{
			name: "non-existent templates dir",
			cfg: &config.Config{
				Languages:    map[string]string{"go": "go"},
				LogLevel:     "info",
				OutputDir:    ".",
				TemplatesDir: "/path/that/does/not/exist/12345",
			},
			wantErr: os.ErrNotExist,
		},
		{
			name: "valid config with normalization",
			cfg: &config.Config{
				Languages: map[string]string{"  GO  ": "  GO  "},
				LogLevel:  "  INFO  ",
				OutputDir: ".",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.cfg.Validate()
			if tt.wantErr != nil { //nolint:nestif // test has multiple error cases
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if tt.cfg != nil && tt.name == "valid config with normalization" {
					if tt.cfg.LogLevel != "info" {
						t.Errorf("expected log level to be normalized to 'info', got %q", tt.cfg.LogLevel)
					}

					if tt.cfg.Languages["go"] != "go" {
						t.Errorf("expected language to be normalized to 'go', got %q", tt.cfg.Languages["go"])
					}
				}
			}
		})
	}
}

func TestConfig_Save(t *testing.T) { //nolint:gocognit // test has multiple subtests and error cases
	t.Parallel()

	t.Run("save and load valid config", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "test-config.json")

		cfg := config.GetDefaultConfig()
		cfg.Author = "test-author"
		cfg.OutputDir = "/custom/output"
		cfg.TemplatesDir = tmpDir
		cfg.LogLevel = logLevelDebug

		if err := cfg.Save(path, false, utils.NewFileSystem()); err != nil {
			t.Fatalf("failed to save config: %v", err)
		}

		loaded, err := config.LoadConfigFromFile(viper.New(), path)
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}

		if loaded.Author != "test-author" {
			t.Errorf("expected Author 'test-author', got %q", loaded.Author)
		}

		if loaded.OutputDir != "/custom/output" {
			t.Errorf("expected OutputDir '/custom/output', got %q", loaded.OutputDir)
		}

		if loaded.TemplatesDir != tmpDir {
			t.Errorf("expected TemplatesDir '%s', got %q", tmpDir, loaded.TemplatesDir)
		}

		if loaded.LogLevel != "debug" {
			t.Errorf("expected LogLevel 'debug', got %q", loaded.LogLevel)
		}

		if loaded.Languages["go"] != "go" {
			t.Errorf("expected Languages['go'] = 'go', got %q", loaded.Languages["go"])
		}
	})

	t.Run("overwrite existing config", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "overwrite-test.json")

		if err := config.GetDefaultConfig().Save(path, false, utils.NewFileSystem()); err != nil {
			t.Fatalf("failed to save initial config: %v", err)
		}

		err := config.GetDefaultConfig().Save(path, false, utils.NewFileSystem())
		if err == nil {
			t.Error("expected error when saving without overwrite")
		}

		cfg := config.GetDefaultConfig()
		cfg.Author = "updated-author"

		if err := cfg.Save(path, true, utils.NewFileSystem()); err != nil {
			t.Errorf("expected no error with overwrite=true, got %v", err)
		}

		loaded, err := config.LoadConfigFromFile(viper.New(), path)
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}

		if loaded.Author != "updated-author" {
			t.Errorf("expected updated author, got %q", loaded.Author)
		}
	})

	t.Run("creates parent directory", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "nested", "deep", "dir", "config.json")

		cfg := config.GetDefaultConfig()
		if err := cfg.Save(path, false, utils.NewFileSystem()); err != nil {
			t.Fatalf("failed to save config with nested dirs: %v", err)
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Error("expected config file to be created")
		}
	})

	t.Run("empty path", func(t *testing.T) {
		t.Parallel()

		cfg := config.GetDefaultConfig()

		err := cfg.Save("", false, utils.NewFileSystem())
		if err == nil {
			t.Error("expected error for empty path")
		}
	})

	t.Run("invalid config", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "invalid-config.json")

		cfg := config.GetDefaultConfig()
		cfg.LogLevel = "invalid"

		err := cfg.Save(path, false, utils.NewFileSystem())
		if !errors.Is(err, config.ErrInvalidConfig) {
			t.Errorf("expected ErrInvalidConfig, got %v", err)
		}
	})
}

func TestLoadConfigFromFile(t *testing.T) {
	t.Parallel()

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		_, err := config.LoadConfigFromFile(viper.New(), "/nonexistent/path/config.json")
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected os.ErrNotExist, got %v", err)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "invalid.json")

		invalidJSON := `{"author": "test", "languages": invalid}`
		if err := os.WriteFile(path, []byte(invalidJSON), 0600); err != nil {
			t.Fatalf("failed to write invalid JSON: %v", err)
		}

		_, err := config.LoadConfigFromFile(viper.New(), path)
		if !errors.Is(err, config.ErrReadConfig) {
			t.Errorf("expected ErrReadConfig, got %v", err)
		}
	})

	t.Run("directory path", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()

		_, err := config.LoadConfigFromFile(viper.New(), tmpDir)
		if !errors.Is(err, config.ErrReadConfig) {
			t.Errorf("expected ErrReadConfig, got %v", err)
		}
	})

	t.Run("invalid log level in file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "invalid-log-level.json")

		invalidConfig := `{"log-level": "trace", "languages": {"go": "go"}, "output-dir": "."}`
		if err := os.WriteFile(path, []byte(invalidConfig), 0600); err != nil {
			t.Fatalf("failed to write invalid config: %v", err)
		}

		_, err := config.LoadConfigFromFile(viper.New(), path)
		if !errors.Is(err, config.ErrInvalidConfig) {
			t.Errorf("expected ErrInvalidConfig, got %v", err)
		}
	})
}

func TestLoadConfigFromDefaultFile(t *testing.T) { //nolint:paralleltest // test modifies environment variables and working directory
	t.Run("not found", func(t *testing.T) { //nolint:paralleltest // modifies environment variables
		tmpHome := t.TempDir()
		test_utils.SetHomeEnv(t, tmpHome)

		_, err := config.LoadConfigFromDefaultFile(viper.New())
		if !errors.Is(err, os.ErrNotExist) && !strings.Contains(err.Error(), "Not Found") {
			t.Errorf("expected os.ErrNotExist or 'Not Found' error, got %v", err)
		}
	})

	t.Run("found in current dir", func(t *testing.T) { //nolint:paralleltest // modifies working directory
		tmpDir := t.TempDir()
		t.Chdir(tmpDir)

		cfg := config.GetDefaultConfig()
		cfg.Author = "current-dir-author"

		configPath := filepath.Join(tmpDir, "scaffy.json")
		if err := cfg.Save(configPath, false, utils.NewFileSystem()); err != nil {
			t.Fatalf("failed to save config: %v", err)
		}

		loaded, err := config.LoadConfigFromDefaultFile(viper.New())
		if err != nil {
			t.Fatalf("failed to load config from default: %v", err)
		}

		if loaded.Author != "current-dir-author" {
			t.Errorf("expected Author 'current-dir-author', got %q", loaded.Author)
		}
	})
}

func TestEnsureDefaultConfig(t *testing.T) { //nolint:paralleltest // test modifies environment variables and working directory
	t.Run("creates new config", func(t *testing.T) { //nolint:paralleltest // modifies environment variables
		tmpHome := t.TempDir()
		test_utils.SetHomeEnv(t, tmpHome)

		configMgr := config.NewConfigManager(utils.NewOSInfo(), utils.NewFileSystem())

		cfg, err := configMgr.EnsureDefaultConfig()
		if err != nil {
			t.Fatalf("failed to ensure default config: %v", err)
		}

		if cfg == nil {
			t.Fatal("expected non-nil config")
		}

		configPath := filepath.Join(config.GetDefaultConfigDir(utils.NewOSInfo()), "scaffy.json")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Errorf("expected config file at %s", configPath)
		}
	})

	t.Run("already exists", func(t *testing.T) { //nolint:paralleltest // modifies environment variables
		tmpHome := t.TempDir()
		test_utils.SetHomeEnv(t, tmpHome)

		configMgr := config.NewConfigManager(utils.NewOSInfo(), utils.NewFileSystem())

		cfg, err := configMgr.EnsureDefaultConfig()
		if err != nil {
			t.Fatalf("first call failed: %v", err)
		}

		_, err = configMgr.EnsureDefaultConfig()
		if err == nil {
			t.Error("expected error when config already exists")
		}

		if cfg == nil {
			t.Error("expected config from first call")
		}
	})
}

func TestEnvVariableOverride(t *testing.T) {
	t.Run("override author and log level", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SCAFFY_AUTHOR", "env-author")
		t.Setenv("SCAFFY_LOG_LEVEL", "warn")

		cfg := config.GetDefaultConfig()
		cfg.Author = "file-author"
		cfg.LogLevel = logLevelDebug

		path := filepath.Join(tmpDir, "config.json")
		if err := cfg.Save(path, false, utils.NewFileSystem()); err != nil {
			t.Fatalf("failed to save config: %v", err)
		}

		loaded, err := config.LoadConfigFromFile(viper.New(), path)
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}

		if loaded.Author != "env-author" {
			t.Errorf("expected env override Author 'env-author', got %q", loaded.Author)
		}

		if loaded.LogLevel != "warn" {
			t.Errorf("expected env override LogLevel 'warn', got %q", loaded.LogLevel)
		}
	})

	t.Run("override output dir and templates dir", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("SCAFFY_OUTPUT_DIR", tmpDir)
		t.Setenv("SCAFFY_TEMPLATES_DIR", tmpDir)

		cfg := config.GetDefaultConfig()
		cfg.OutputDir = "/file/output"
		cfg.TemplatesDir = ""

		path := filepath.Join(tmpDir, "config.json")
		if err := cfg.Save(path, false, utils.NewFileSystem()); err != nil {
			t.Fatalf("failed to save config: %v", err)
		}

		loaded, err := config.LoadConfigFromFile(viper.New(), path)
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}

		if loaded.OutputDir != tmpDir {
			t.Errorf("expected env override OutputDir '%s', got %q", tmpDir, loaded.OutputDir)
		}

		if loaded.TemplatesDir != tmpDir {
			t.Errorf("expected env override TemplatesDir '%s', got %q", tmpDir, loaded.TemplatesDir)
		}
	})
}

func TestConfigJSONFormat(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "format-test.json")

	cfg := config.GetDefaultConfig()
	cfg.Author = "json-test"
	cfg.Languages = map[string]string{"go": "go", "rust": "rs"}

	if err := cfg.Save(path, false, utils.NewFileSystem()); err != nil {
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
