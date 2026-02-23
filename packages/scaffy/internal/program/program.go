package program

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/25prabhu10/scaffy/internal/boilerplate"
	"github.com/25prabhu10/scaffy/internal/config"
	"github.com/25prabhu10/scaffy/internal/logger"
	"github.com/25prabhu10/scaffy/internal/templates"
	"github.com/25prabhu10/scaffy/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

// contextKey is a private type for context keys to avoid collisions.
type contextKey struct{}

// programKey is the context key for storing the Program instance.
var programKey = contextKey{}

// Error definitions for the Program.
var (
	ErrConfigExists     = errors.New("config file already exists")
	ErrInvalidDocFormat = errors.New("invalid documentation format")
)

// Documentation formats supported by the llm command.
const (
	Markdown = "markdown"
	Man      = "man"
	Rest     = "rest"
)

var (
	AllowedDocFormats    = []string{Markdown, Man, Rest}
	AllowedDocFormatsStr = strings.Join(AllowedDocFormats, "|")
)

// Program encapsulates the core state and configuration of the application. It holds the loaded configuration and flags for verbose and quiet modes.
type Program struct {
	config *config.Config
	logger *slog.Logger
	fs     utils.FileSystem
	osInfo utils.OSInfo
}

type ProgramOption func(*Program)

func WithFileSystem(fs utils.FileSystem) ProgramOption {
	return func(p *Program) { p.fs = fs }
}

func WithOSInfo(osInfo utils.OSInfo) ProgramOption {
	return func(p *Program) { p.osInfo = osInfo }
}

// New creates a new Program instance by loading the configuration from the specified file or default locations. It also sets the verbose and quiet flags based on the provided parameters.
func New(configFile string, verbose bool, quiet bool, options ...ProgramOption) (*Program, error) {
	cfg, err := loadConfiguration(configFile, verbose, quiet)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	p := &Program{
		config: cfg,
		logger: logger.New(cfg.LogLevel, verbose, quiet),
		fs:     utils.NewFileSystem(),
		osInfo: utils.NewOSInfo(),
	}

	for _, option := range options {
		option(p)
	}

	return p, nil
}

// InitializeConfig handles directory creation, path resolution, and config saving.
func (p *Program) InitializeConfig(cfg *config.Config, outputDir string, force bool) (string, error) {
	p.logger.Debug("initializing new config", "force", force)

	configFilePath := config.DEFAULT_CONFIG_FILE_NAME

	if outputDir != "" {
		if err := utils.CreateDirectoryIfNotExists(outputDir, p.fs); err != nil {
			return "", fmt.Errorf("failed to create output directory: %w", err)
		}

		configFilePath = filepath.Join(outputDir, configFilePath)
	}

	if !force {
		if _, err := os.Stat(configFilePath); err == nil {
			return configFilePath, ErrConfigExists
		}
	}

	if err := cfg.Save(configFilePath, force, p.fs); err != nil {
		return configFilePath, fmt.Errorf("failed to create config: %w", err)
	}

	p.logger.Debug("new config initialized successfully", "path", configFilePath)

	return configFilePath, nil
}

func (p *Program) GenerateFilesFromTemplates(ctx context.Context, name string, force bool) error {
	p.logger.Debug("Generating files from templates...")
	p.logger.Debug("Values", "config", p.config, "name", name)

	tm := templates.New(p.config.TemplatesDir, templates.Data{
		Date:   utils.GetCurrentDate(),
		Author: p.config.Author,
		URL:    "",
	})

	bp := boilerplate.New(tm, p.fs, p.logger)

	options := boilerplate.Options{
		Name:      name,
		OutputDir: p.config.OutputDir,
		Languages: p.config.Languages,
		Force:     force,
	}

	return bp.Scaffold(ctx, options)
}

// GenerateLLMDocs handles directory creation and documentation generation for various formats.
func (p *Program) GenerateLLMDocs(root *cobra.Command, outputDir string, format string, frontMatter bool) error {
	p.logger.Debug("generating documentation", "path", outputDir, "format", format, "frontMatter", frontMatter)

	if err := utils.CreateDirectoryIfNotExists(outputDir, p.fs); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	root.DisableAutoGenTag = true

	switch format {
	case Markdown:
		if frontMatter {
			prep := func(filename string) string {
				base := filepath.Base(filename)
				name := strings.TrimSuffix(base, filepath.Ext(base))
				title := strings.ReplaceAll(name, "_", " ")

				return fmt.Sprintf("---\ntitle: %q\nslug: %q\ndescription: \"CLI reference for %s\"\n---\n\n", title, name, title)
			}
			if err := doc.GenMarkdownTreeCustom(root, outputDir, prep, strings.ToLower); err != nil {
				return fmt.Errorf("error in generating markdown docs with frontmatter: %w", err)
			}
		} else {
			if err := doc.GenMarkdownTree(root, outputDir); err != nil {
				return fmt.Errorf("error in generating markdown docs: %w", err)
			}
		}
	case Man:
		hdr := &doc.GenManHeader{Title: strings.ToUpper(root.Name()), Section: "1"}
		if err := doc.GenManTree(root, hdr, outputDir); err != nil {
			return fmt.Errorf("error in generating man pages: %w", err)
		}
	case Rest:
		if err := doc.GenReSTTree(root, outputDir); err != nil {
			return fmt.Errorf("error in generating reStructuredText docs: %w", err)
		}
	default:
		return fmt.Errorf("%w: %q (valid: %s)", ErrInvalidDocFormat, format, AllowedDocFormatsStr)
	}

	p.logger.Debug("documentation generated successfully", "path", outputDir, "format", format)

	return nil
}

// WithProgram returns a new context with the Program instance attached.
func WithProgram(ctx context.Context, p *Program) context.Context {
	return context.WithValue(ctx, programKey, p)
}

// FromContext retrieves the Program instance from the context. If not found, it returns a default Program instance.
func FromContext(ctx context.Context) *Program {
	if u, ok := ctx.Value(programKey).(*Program); ok {
		return u
	}

	// Return default Program if not found in context
	p, err := New("", false, false)
	if err != nil {
		panic(fmt.Errorf("failed to create default program: %w", err))
	}

	return p
}

// loadConfiguration loads the configuration from the specified file or default locations. If the config file is not found, it attempts to create a default config. It returns the loaded configuration or an error if loading fails.
func loadConfiguration(configFile string, verbose bool, quiet bool) (*config.Config, error) {
	cfg, err := loadConfig(configFile)
	if err != nil {
		// if config not found, try to create default config
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			var createErr error

			configMgr := config.NewConfigManager(utils.NewOSInfo(), utils.NewFileSystem())

			cfg, createErr = configMgr.EnsureDefaultConfig()
			if createErr != nil {
				return nil, fmt.Errorf("failed to load or create config: %w", createErr)
			}

			// User-facing message about config creation
			localLogger := logger.New(slog.LevelDebug.String(), verbose, quiet)
			localLogger.Debug("created default configuration", "path",
				config.GetDefaultConfigDir(utils.NewOSInfo())+"/"+config.DEFAULT_CONFIG_FILE_NAME)
		} else {
			return nil, fmt.Errorf("%w", err)
		}
	}

	return cfg, nil
}

// loadConfig loads configuration from file or default locations.
func loadConfig(configFile string) (*config.Config, error) {
	if configFile != "" {
		return config.LoadConfigFromFile(viper.GetViper(), configFile)
	}

	return config.LoadConfigFromDefaultFile(viper.GetViper())
}
