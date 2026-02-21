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
	"slices"
	"strings"

	"github.com/spf13/viper"

	"github.com/25prabhu10/scaffy/internal/constants"
	"github.com/25prabhu10/scaffy/internal/utils"
	"github.com/25prabhu10/scaffy/pkg/build_info"
)

// Configuration-related errors.
var (
	ErrReadConfig                 = errors.New("failed to read config")
	ErrInvalidConfig              = errors.New("invalid config")
	ErrNilConfig                  = errors.New("config is nil")
	ErrAuthorEmpty                = errors.New("author cannot be empty")
	ErrAuthorInvalidLength        = fmt.Errorf("author name cannot be longer than %d characters", constants.MAX_NAME_LENGTH)
	ErrInvalidLogLevel            = errors.New("invalid log-level")
	ErrLanguagesEmpty             = errors.New("languages cannot be empty")
	ErrLanguageEmpty              = errors.New("language cannot be empty")
	ErrLanguageInvalidLength      = fmt.Errorf("language name cannot be longer than %d characters", constants.MAX_NAME_LENGTH)
	ErrFileExtensionEmpty         = errors.New("language file extension cannot be empty")
	ErrFileExtensionInvalidLength = fmt.Errorf("file extension cannot be longer than %d characters", constants.MAX_NAME_LENGTH)
	ErrOutputDirEmpty             = errors.New("output-dir cannot be empty")
)

// Log Levels.
var (
	logLevelDebug  = strings.ToLower(slog.LevelDebug.String())
	logLevelInfo   = strings.ToLower(slog.LevelInfo.String())
	logLevelWarn   = strings.ToLower(slog.LevelWarn.String())
	logLevelError  = strings.ToLower(slog.LevelError.String())
	validLogLevels = []string{logLevelDebug, logLevelInfo, logLevelWarn, logLevelError}
)

// CURRENT_DIR is the current working directory, used as the default output directory in the config.
const CURRENT_DIR = "."

// CONFIG_FORMAT is the configuration file format.
const CONFIG_FORMAT = "json"

// DEFAULT_CONFIG_FILE_NAME is the default configuration file name.
var DEFAULT_CONFIG_FILE_NAME = build_info.APP_NAME + "." + CONFIG_FORMAT

// configuration keys.
const (
	KEY_AUTHOR        = "author"
	KEY_LANGUAGES     = "languages"
	KEY_LOG_LEVEL     = "log-level"
	KEY_OUTPUT_DIR    = "output-dir"
	KEY_TEMPLATES_DIR = "templates-dir"
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
		OutputDir:    CURRENT_DIR,
		TemplatesDir: "",
	}
}

// GetDefaultConfigPath returns the default configuration file path as a
// string. This is used for display purposes in help messages.
func GetDefaultConfigPath(info utils.OSInfo) string {
	var homePath string

	osname := info.GetOS()

	switch osname {
	case "windows":
		homePath = "%USERPROFILE%"
	case "darwin", "ios":
		homePath = "$HOME"
	default:
		homePath = "$XDG_CONFIG_HOME"
	}

	return filepath.Join(homePath, build_info.APP_NAME, DEFAULT_CONFIG_FILE_NAME)
}

// GetDefaultConfigDir returns the actual default configuration directory path.
//
// On Unix systems, it returns $XDG_CONFIG_HOME as specified by
//
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if
// non-empty, else $HOME/.config.
// On Darwin, it returns $HOME/.config/Application Support.
// On Windows, it returns %AppData%.
func GetDefaultConfigDir(osInfo utils.OSInfo) string {
	userConfigDir, err := osInfo.GetUserConfigDir()
	if err != nil {
		// Fallback to a reasonable default if we can't get user config
		// dir
		home, _ := os.UserHomeDir()
		userConfigDir = filepath.Join(home, ".config")
	}

	return filepath.Join(userConfigDir, build_info.APP_NAME)
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
	v.AddConfigPath(CURRENT_DIR)
	v.AddConfigPath(GetDefaultConfigDir(utils.NewOSInfo()))
	v.SetConfigType(CONFIG_FORMAT)
	v.SetConfigName(build_info.APP_NAME)

	return readConfig(v)
}

// LoadConfigFromFile loads configuration from the specified file path.
// Environment variables with SCAFFY_ prefix can override config values.
func LoadConfigFromFile(path string) (*Config, error) {
	v := newReadViper()
	v.SetConfigFile(path)

	return readConfig(v)
}

// EnsureDefaultConfig creates a default configuration file if none exists.
// It creates the file at the default config directory and logs the action.
// Returns the created config and any error encountered.
func EnsureDefaultConfig() (*Config, error) {
	cfg := GetDefaultConfig()
	configPath := filepath.Join(GetDefaultConfigDir(utils.NewOSInfo()), DEFAULT_CONFIG_FILE_NAME)

	if err := cfg.Save(configPath, false); err != nil {
		return nil, fmt.Errorf("failed to create default config: %w", err)
	}

	return cfg, nil
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

// Validate the configuration values.
func (c *Config) Validate() error {
	if c == nil {
		return ErrNilConfig
	}

	if c.Author != "" {
		if utils.IsStringEmpty(c.Author) {
			return ErrAuthorEmpty
		} else if utils.IsStringOverMaxLength(c.Author) {
			return fmt.Errorf("%w, must be at most %d characters", ErrAuthorInvalidLength, constants.MAX_NAME_LENGTH)
		}
	}

	c.LogLevel = strings.ToLower(strings.TrimSpace(c.LogLevel))

	if !isValidLogLevel(c.LogLevel) {
		return fmt.Errorf("%w: '%s', must be one of: %s", ErrInvalidLogLevel, c.LogLevel, AllLogLevelsStr())
	}

	if len(c.Languages) == 0 {
		return ErrLanguagesEmpty
	}

	var normalizedLanguages = make(map[string]string)

	for language, fileExtension := range c.Languages {
		if utils.IsStringEmpty(language) {
			return ErrLanguageEmpty
		} else if utils.IsStringOverMaxLength(language) {
			return fmt.Errorf("%w, must be at most %d characters", ErrLanguageInvalidLength, constants.MAX_NAME_LENGTH)
		}

		if utils.IsStringEmpty(fileExtension) {
			return ErrFileExtensionEmpty
		} else if utils.IsStringOverMaxLength(fileExtension) {
			return fmt.Errorf("%w, must be at most %d characters", ErrFileExtensionInvalidLength, constants.MAX_NAME_LENGTH)
		}

		normalizedLanguages[utils.NormalizeString(language)] = utils.NormalizeString(fileExtension)
	}

	c.Languages = normalizedLanguages

	if utils.IsStringEmpty(c.OutputDir) {
		return ErrOutputDirEmpty
	}

	if !utils.IsStringEmpty(c.TemplatesDir) {
		if _, err := os.Stat(c.TemplatesDir); os.IsNotExist(err) {
			return fmt.Errorf("templates-dir '%s': %w", c.TemplatesDir, err)
		}
	}

	return nil
}

func setDefaultConfigValues(c *Config) *viper.Viper {
	v := viper.New()
	v.Set(KEY_AUTHOR, c.Author)
	v.Set(KEY_LANGUAGES, c.Languages)
	v.Set(KEY_LOG_LEVEL, c.LogLevel)
	v.Set(KEY_OUTPUT_DIR, c.OutputDir)
	v.Set(KEY_TEMPLATES_DIR, c.TemplatesDir)

	return v
}

func newReadViper() *viper.Viper {
	v := viper.New()
	setViperDefaults(v)

	v.SetEnvPrefix(strings.ToUpper(utils.NormalizeString(build_info.APP_NAME)))
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	return v
}

func readConfig(v *viper.Viper) (*Config, error) {
	if err := v.ReadInConfig(); err != nil {
		if isConfigNotFound(err) {
			return nil, err
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

func setViperDefaults(v *viper.Viper) {
	defaults := GetDefaultConfig()
	v.SetDefault(KEY_LANGUAGES, defaults.Languages)
	v.SetDefault(KEY_OUTPUT_DIR, defaults.OutputDir)
	v.SetDefault(KEY_TEMPLATES_DIR, defaults.TemplatesDir)
	v.SetDefault(KEY_AUTHOR, defaults.Author)
	v.SetDefault(KEY_LOG_LEVEL, defaults.LogLevel)
}
