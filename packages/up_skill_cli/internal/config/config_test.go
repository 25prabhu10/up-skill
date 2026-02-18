package config_test

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/25prabhu10/up-skill/internal/config"
	"github.com/25prabhu10/up-skill/internal/utils"
)

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
		t.Error("expected non-empty Languages list")
	}

	expectedLangs := config.GetDefaultLanguages()
	if len(cfg.Languages) != len(*expectedLangs) {
		t.Errorf("expected %d languages, got %d", len(*expectedLangs), len(cfg.Languages))
	} else {
		for lang, fileExt := range *expectedLangs {
			if cfg.Languages[lang] != fileExt {
				t.Errorf("expected language %q with %q file extension, got %q", lang, fileExt, cfg.Languages[lang])
			}
		}
	}

	if cfg.LogLevel != strings.ToLower(slog.LevelError.String()) {
		t.Errorf("expected LogLevel '%s', got '%q'", strings.ToLower(slog.LevelError.String()), cfg.LogLevel)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, config.DEFAULT_CONFIG_FILE_NAME)

	config.DEFAULT_CONFIG_FILE_NAME = "test-app-1.json"

	// Save config with all valid fields
	cfg := config.GetDefaultConfig()

	cfg.Author = "test-author"

	if err := cfg.Save(path, false); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Load config
	loaded, err := config.LoadConfigFromFile(path)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loaded.Author != "test-author" {
		t.Errorf("expected Author 'test-author', got \n\n%v, \n\n%q \n\n%v", loaded, path, cfg)
	}

	// Verify all required fields were loaded correctly
	if utils.IsStringEmpty(loaded.OutputDir) {
		t.Error("OutputDir should not be empty after load")
	}
}

func TestSaveConfig_Overwrite(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, config.DEFAULT_CONFIG_FILE_NAME)

	config.DEFAULT_CONFIG_FILE_NAME = "test-app-2.json"

	// Create first config
	if err := config.GetDefaultConfig().Save(path, false); err != nil {
		t.Fatalf("failed to save initial config: %v", err)
	}

	// Try to save without overwrite - should fail
	if err := config.GetDefaultConfig().Save(path, false); err == nil {
		t.Error("expected error when saving without overwrite")
	}

	// Save with overwrite - should succeed
	if err := config.GetDefaultConfig().Save(path, true); err != nil {
		t.Errorf("expected no error with overwrite=true, got %v", err)
	}
}

func TestLoadConfigFromFile_NotFound(t *testing.T) {
	t.Parallel()

	_, err := config.LoadConfigFromFile("/nonexistent/path.json")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestEnsureDefaultConfig(t *testing.T) { //nolint:paralleltest // `t.Setenv` does not support
	// Save original HOME and create isolated test environment
	tmpHome := t.TempDir()

	// Override
	switch runtime.GOOS {
	case "windows":
		t.Setenv("AppData", tmpHome)
		t.Setenv("USERPROFILE", tmpHome)
	case "darwin", "ios":
		t.Setenv("HOME", tmpHome)
	default:
		t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmpHome, ".config"))
	}

	config.DEFAULT_CONFIG_FILE_NAME = "test-app-3.json"

	cfg, err := config.EnsureDefaultConfig()
	if err != nil {
		t.Fatalf("failed to ensure default config: %v", err)
	}

	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	// Verify file was created
	configPath := filepath.Join(config.GetDefaultConfigDir(), config.DEFAULT_CONFIG_FILE_NAME)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("expected config file at %s", configPath)
	}
}

func TestGetDefaultConfigDir(t *testing.T) {
	t.Parallel()

	dir := config.GetDefaultConfigDir()
	if utils.IsStringEmpty(dir) {
		t.Error("expected non-empty config directory")
	}
}
