package program

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/25prabhu10/scaffy/internal/config"
	"github.com/25prabhu10/scaffy/internal/logger"
)

// contextKey is a private type for context keys to avoid collisions.
type contextKey struct{}

// programKey is the context key for storing the Program instance.
var programKey = contextKey{}

// Program encapsulates the core state and configuration of the application. It holds the loaded configuration and flags for verbose and quiet modes.
type Program struct {
	config  *config.Config
	verbose bool
	quiet   bool
	logger  *slog.Logger
}

// New creates a new Program instance by loading the configuration from the specified file or default locations. It also sets the verbose and quiet flags based on the provided parameters.
func New(configFile string, verbose bool, quiet bool) (*Program, error) {
	cfg, err := loadConfiguration(configFile, verbose, quiet)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return &Program{
		config:  cfg,
		verbose: verbose,
		quiet:   quiet,
		logger:  logger.New(cfg.LogLevel, verbose, quiet),
	}, nil
}

func (p *Program) CreateNewConfig(cfg *config.Config, configFile string, force bool) error {
	p.logger.Debug("initializing scaffy config", "path", configFile, "force", force)

	if err := cfg.Save(configFile, force); err != nil {
		return err
	}

	p.logger.Debug("config initialization completed", "path", configFile, "force", force)

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
		if errors.Is(err, config.ErrConfigNotFound) {
			var createErr error

			cfg, createErr = config.EnsureDefaultConfig()
			if createErr != nil {
				return nil, fmt.Errorf("failed to load or create config: %w", createErr)
			}

			// User-facing message about config creation
			localLogger := logger.New(slog.LevelDebug.String(), verbose, quiet)
			localLogger.Debug("created default configuration", "path",
				config.GetDefaultConfigDir()+"/"+config.DEFAULT_CONFIG_FILE_NAME)
		} else {
			return nil, fmt.Errorf("%w", err)
		}
	}

	return cfg, nil
}

// loadConfig loads configuration from file or default locations.
func loadConfig(configFile string) (*config.Config, error) {
	if configFile != "" {
		return config.LoadConfigFromFile(configFile)
	}

	return config.LoadConfigFromDefaultFile()
}
