package commands

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/25prabhu10/scaffy/internal/config"
	"github.com/25prabhu10/scaffy/internal/program"
	"github.com/25prabhu10/scaffy/internal/ui"
)

// init command errors.
var (
	ErrConfigPathIsDir = errors.New("config file path is a directory")
)

var force bool

// NewInitCmd creates the init subcommand for initializing scaffy configuration.
func NewInitCmd() *cobra.Command {
	defaultCfg := config.GetDefaultConfig()

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a configuration file",
		Args:  cobra.MaximumNArgs(1),
		Long: fmt.Sprintf(`
Initialize creates a configuration file for scaffy.

It creates a local config file "%s" in the specified directory (default: current directory).`, config.DEFAULT_CONFIG_FILE_NAME),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigInitializerE(cmd, args, defaultCfg, force)
		},
	}

	initCmd.Flags().BoolVar(&force, "force", false, "allow write operations that overwrite files")

	// config overrides
	initCmd.Flags().StringVarP(&defaultCfg.Author, "author", "a", defaultCfg.Author, "author name")
	initCmd.Flags().StringToStringVar(&defaultCfg.Languages, "lang", defaultCfg.Languages, "programming languages to generate (e.g. --lang ruby=rb,python=py)")
	initCmd.Flags().StringVar(&defaultCfg.LogLevel, "log-level", defaultCfg.LogLevel, "slog level (e.g., debug, info, warn, error)")
	initCmd.Flags().StringVar(&defaultCfg.OutputDir, "output-dir", defaultCfg.OutputDir, "directory to output generated code")
	initCmd.Flags().StringVar(&defaultCfg.TemplatesDir, "templates-dir", defaultCfg.TemplatesDir, "directory containing template files")

	return initCmd
}

// runConfigInitializerE executes the logic for the init command, creating a config file with default settings.
func runConfigInitializerE(cmd *cobra.Command, args []string, cfg *config.Config, force bool) error {
	ctx := cmd.Context()

	userUI := ui.FromContext(ctx)
	app := program.FromContext(ctx)

	outputDir := ""
	if len(args) > 0 {
		outputDir = args[0]
	}

	if force {
		userUI.Warningf("existing config will be overwritten (--force)")
	}

	configFilePath, err := app.InitializeConfig(cfg, outputDir, force)
	if err != nil {
		if errors.Is(err, program.ErrConfigExists) {
			return fmt.Errorf("%w at %s (use --force to overwrite)", err, configFilePath)
		}

		return fmt.Errorf("failed to initialize config: %w", err)
	}

	userUI.Infof("initialized scaffy with default config at %s", configFilePath)

	return nil
}
