// Package config provides configuration management for the scaffy CLI tool.
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
	"slices"
	"strings"

	"github.com/spf13/viper"

	"github.com/25prabhu10/scaffy/internal/utils"
	"github.com/25prabhu10/scaffy/pkg/build_info"
)

// Error variables for configuration-related errors.
var (
	// ErrConfigNotFound indicates no configuration file was found and no
	// environment-based configuration was provided.
	ErrConfigNotFound = errors.New("config not found")
	// ErrReadConfig indicates configuration could not be read or unmarshaled.
	ErrReadConfig = errors.New("failed to read config")
	// ErrInvalidConfig indicates configuration values are invalid.
	ErrInvalidConfig = errors.New("invalid config")
	// ErrNilConfig indicates a nil Config pointer was provided where a non-nil value is required.
	ErrNilConfig = errors.New("config is nil")
	// ErrInvalidLogLevel indicates the log level specified in the configuration is not valid.
	ErrInvalidLogLevel = errors.New("invalid log-level")
	// ErrLanguageEmpty indicates the language is empty in the configuration.
	ErrLanguageEmpty = errors.New("language cannot be empty")
	// ErrFileExtensionEmpty indicates a language file extension in the configuration is empty.
	ErrFileExtensionEmpty = errors.New("language file extension cannot be empty")
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

// DEFAULT_CONFIG_FILE_NAME is the default configuration file name.
var DEFAULT_CONFIG_FILE_NAME = build_info.AppName + "." + CONFIG_FORMAT

// configuration keys.
const (
	keyAuthor       = "author"
	keyLanguages    = "languages"
	keyLogLevel     = "log-level"
	keyOutputDir    = "output-dir"
	keyTemplatesDir = "templates-dir"
)

// Config represents the scaffy configuration settings. It can be loaded from a
// JSON file or environment variables prefixed with SCAFFY_.
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
		Languages:    GetDefaultLanguages(),
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

// GetDefaultLanguages returns the default list of supported programming languages.
func GetDefaultLanguages() map[string]string {
	return map[string]string{
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
//  1. Current directory (./scaffy.json)
//  2. XDG config directory (~/.config/scaffy/scaffy.json)
//
// If no config file is found, it returns an error.
// Environment variables with SCAFFY_ prefix can override config values.
func LoadConfigFromDefaultFile() (*Config, error) {
	v := newReadViper()
	v.AddConfigPath(currentDir)
	v.AddConfigPath(GetDefaultConfigDir())
	v.SetConfigType(CONFIG_FORMAT)
	v.SetConfigName(build_info.AppName)

	return readConfig(v)
}

// LoadConfigFromFile loads configuration from the specified file path.
// Environment variables with SCAFFY_ prefix can override config values.
func LoadConfigFromFile(path string) (*Config, error) {
	v := newReadViper()
	v.SetConfigFile(path)

	return readConfig(v)
}

// Save saves the configuration to the specified file path.
// If overwrite is true, it will overwrite an existing file; otherwise, it
// returns an error.
func (c *Config) Save(path string, overwrite bool) error {
	if err := c.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidConfig, err)
	}

	v := setDefaultConfigValues(c)

	configDir := filepath.Dir(path)

	if err := utils.CreateDirectoryIfNotExists(configDir); err != nil {
		return err
	}

	if overwrite {
		return v.WriteConfigAs(path)
	}

	return v.SafeWriteConfigAs(path)
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

// Validate validates the configuration values.
func (c *Config) Validate() error {
	if c == nil {
		return ErrNilConfig
	}

	c.LogLevel = strings.ToLower(strings.TrimSpace(c.LogLevel))

	if !isValidLogLevel(c.LogLevel) {
		return fmt.Errorf("%w: '%s', must be one of: %s", ErrInvalidLogLevel, c.LogLevel, AllLogLevelsStr())
	}

	for language, extension := range c.Languages {
		if utils.IsStringEmpty(strings.TrimSpace(language)) {
			return ErrLanguageEmpty
		}

		if utils.IsStringEmpty(strings.TrimSpace(extension)) {
			return ErrFileExtensionEmpty
		}
	}

	return nil
}

// setDefaultConfigValues sets default values in a viper instance based on the provided Config struct. This is used to ensure that all config fields have
// default values when saving or validating the config.
func setDefaultConfigValues(c *Config) *viper.Viper {
	v := viper.New()
	v.Set(keyAuthor, c.Author)
	v.Set(keyLanguages, c.Languages)
	v.Set(keyLogLevel, c.LogLevel)
	v.Set(keyOutputDir, c.OutputDir)
	v.Set(keyTemplatesDir, c.TemplatesDir)

	return v
}

func newReadViper() *viper.Viper {
	v := viper.New()
	setViperDefaults(v)

	v.SetEnvPrefix(strings.ToUpper(utils.NormalizeString(build_info.AppName)))
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	return v
}

// readConfig reads and unmarshals the configuration from a viper
// instance. It returns an error if the config file is not found or cannot be
// read.
func readConfig(v *viper.Viper) (*Config, error) {
	if err := v.ReadInConfig(); err != nil {
		if isConfigNotFound(err) {
			return nil, fmt.Errorf("config file not found: %w", err)
		}

		return nil, fmt.Errorf("%w: %w", ErrReadConfig, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadConfig, err)
	}

	normalizeConfig(&cfg)

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfig, err)
	}

	return &cfg, nil
}

func isConfigNotFound(err error) bool {
	var configNotFound viper.ConfigFileNotFoundError

	return os.IsNotExist(err) || errors.As(err, &configNotFound)
}

func normalizeConfig(cfg *Config) {
	cfg.LogLevel = strings.ToLower(strings.TrimSpace(cfg.LogLevel))
}

func isValidLogLevel(level string) bool {
	return slices.Contains(validLogLevels, level)
}

// setViperDefaults sets default values in a viper instance for all config
// fields.
func setViperDefaults(v *viper.Viper) {
	defaults := GetDefaultConfig()
	v.SetDefault(keyLanguages, defaults.Languages)
	v.SetDefault(keyOutputDir, defaults.OutputDir)
	v.SetDefault(keyTemplatesDir, defaults.TemplatesDir)
	v.SetDefault(keyAuthor, defaults.Author)
	v.SetDefault(keyLogLevel, defaults.LogLevel)
}
