// Package config provides configuration management for the up-skill CLI tool.
//
// It handles loading, saving, and validating configuration from files and
// environment variables using Viper for configuration management with support
// for JSON format and XDG base directories.
package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/25prabhu10/up-skill/internal/utils"
	"github.com/25prabhu10/up-skill/pkg/build_info"
	"github.com/spf13/viper"
)

// Log Levels.
var (
	logLevelDebug  = strings.ToLower(slog.LevelDebug.String())
	logLevelInfo   = strings.ToLower(slog.LevelInfo.String())
	logLevelWarn   = strings.ToLower(slog.LevelWarn.String())
	logLevelError  = strings.ToLower(slog.LevelError.String())
	validLogLevels = []string{logLevelDebug, logLevelInfo, logLevelWarn, logLevelError}
	currentDir     = "."
)

// CONFIG_FORMAT is the configuration file format.
const CONFIG_FORMAT = "json"

// DEFAULT_CONFIG_FILE is the default configuration file name.
var DEFAULT_CONFIG_FILE_NAME = build_info.AppName + "." + CONFIG_FORMAT

// Config represents the up-skill configuration settings. It can be loaded from a
// JSON file or environment variables prefixed with UP_SKILL_.
type Config struct {
	Author       string            `mapstructure:"author"`
	Languages    map[string]string `mapstructure:"languages"`
	LogLevel     string            `mapstructure:"log-level"`
	OutputDir    string            `mapstructure:"output-dir"`
	TemplatesDir string            `mapstructure:"templates-dir"`
}

// GetDefaultConfig returns a new Config with sensible default values.
func GetDefaultConfig() *Config {
	return &Config{
		Author:       "",
		Languages:    *GetDefaultLanguages(),
		LogLevel:     logLevelError,
		OutputDir:    currentDir,
		TemplatesDir: "",
	}
}

// GetDefaultConfigPath returns the default configuration file path as a
// string. This is used for display purposes in help messages.
func GetDefaultConfigPath() string {
	var homePath string

	switch runtime.GOOS {
	case "windows":
		homePath = "%USERPROFILE%"
	case "darwin", "ios":
		homePath = "$HOME"
	default:
		homePath = "$XDG_CONFIG_HOME"
	}

	return fmt.Sprintf("%s/%s/%s", homePath, build_info.AppName, DEFAULT_CONFIG_FILE_NAME)
}

// GetDefaultConfigDir returns the actual default configuration directory path.
//
// On Unix systems, it returns $XDG_CONFIG_HOME as specified by
//
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if
// non-empty, else $HOME/.config.
// On Darwin, it returns $HOME/.config/Application Support.
// On Windows, it returns %AppData%.
func GetDefaultConfigDir() string {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to a reasonable default if we can't get user config
		// dir
		home, _ := os.UserHomeDir()
		userConfigDir = filepath.Join(home, ".config")
	}

	return filepath.Join(userConfigDir, build_info.AppName)
}

// DefaultLanguages returns the default list of supported programming languages.
func GetDefaultLanguages() *map[string]string {
	return &map[string]string{
		"go":         "go",
		"c":          "c",
		"python":     "py",
		"javascript": "js",
	}
}

// VerboseLogLevel returns the log level string for verbose mode (debug).
func VerboseLogLevel() string {
	return logLevelDebug
}

// QuietLogLevel returns the log level string for quiet mode (error).
func QuietLogLevel() string {
	return logLevelError
}

// AllLogLevelsStr returns a comma-separated string of all valid log levels.
func AllLogLevelsStr() string {
	return strings.Join(validLogLevels, ",")
}

// LoadConfigFromDefaultFile loads configuration from the default search paths.
// It searches for config files in the following order:
//  1. Current directory (./up-skill.json)
//  2. XDG config directory (~/.config/up-skill/up-skill.json)
//
// If no config file is found, it returns an error.
// Environment variables with UP_SKILL_ prefix can override config values.
func LoadConfigFromDefaultFile() (*Config, error) {
	viper.AddConfigPath(currentDir)
	viper.AddConfigPath(GetDefaultConfigDir())
	viper.SetConfigType(CONFIG_FORMAT)
	viper.SetConfigName(build_info.AppName)

	return readConfig()
}

// LoadConfigFromFile loads configuration from the specified file path.
// Environment variables with UP_SKILL_ prefix can override config values.
func LoadConfigFromFile(path string) (*Config, error) {
	viper.SetConfigFile(path)

	return readConfig()
}

// Save saves the configuration to the specified file path.
// If overwrite is true, it will overwrite an existing file; otherwise, it
// returns an error.
func (c *Config) Save(path string, overwrite bool) error {
	v := setDefaultConfigValues(c)

	configDir := filepath.Dir(path)

	if err := utils.CreateDirectoryIfNotExists(configDir); err != nil {
		return err
	}

	if overwrite {
		return v.WriteConfigAs(path)
	}

	return viper.SafeWriteConfigAs(path)
}

// EnsureDefaultConfig creates a default configuration file if none exists.
// It creates the file at the default config directory and logs the action.
// Returns the created config and any error encountered.
func EnsureDefaultConfig() (*Config, error) {
	cfg := GetDefaultConfig()
	configPath := filepath.Join(GetDefaultConfigDir(), DEFAULT_CONFIG_FILE_NAME)

	if err := cfg.Save(configPath, false); err != nil {
		return nil, fmt.Errorf("failed to create default config: %w", err)
	}

	return cfg, nil
}

// setDefaultConfigValues sets default values in a viper instance based on the provided Config struct. This is used to ensure that all config fields have
// default values when saving or validating the config.
func setDefaultConfigValues(c *Config) *viper.Viper {
	v := viper.New()
	v.Set("author", c.Author)
	v.Set("languages", c.Languages)
	v.Set("log-level", c.LogLevel)
	v.Set("output-dir", c.OutputDir)
	v.Set("templates-dir", c.TemplatesDir)

	return v
}

// readConfig reads and unmarshals the configuration from a viper
// instance. It returns an error if the config file is not found or cannot be
// read.
func readConfig() (*Config, error) {
	setViperDefaults()

	viper.SetEnvPrefix(strings.ToUpper(build_info.AppName))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var configNotFound viper.ConfigFileNotFoundError
		if ok := os.IsNotExist(err); ok || errors.As(err, &configNotFound) {
			return nil, fmt.Errorf("config file not found: %w", err)
		}

		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// if err := cfg.Validate(); err != nil {
	// 	return nil, fmt.Errorf("invalid options: %w", err)
	// }

	return &cfg, nil
}

// setViperDefaults sets default values in a viper instance for all config
// fields.
func setViperDefaults() {
	defaults := GetDefaultConfig()
	viper.SetDefault("languages", defaults.Languages)
	viper.SetDefault("output-dir", defaults.OutputDir)
	viper.SetDefault("templates-dir", defaults.TemplatesDir)
	viper.SetDefault("author", defaults.Author)
	viper.SetDefault("log-level", defaults.LogLevel)
}
