package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/25prabhu10/scaffy/internal/config"
	"github.com/25prabhu10/scaffy/internal/program"
	"github.com/25prabhu10/scaffy/internal/ui"
	"github.com/25prabhu10/scaffy/internal/utils"
)

var (
	ErrConfigExists    = errors.New("config file already exists")
	ErrConfigPathIsDir = errors.New("config file path is a directory")
)

// NewInitCmd creates the init subcommand for initializing scaffy configuration.
func NewInitCmd() *cobra.Command {
	defaultCfg := config.GetDefaultConfig()

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a configuration file",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
				return err
			}

			// create directory structure
			if len(args) == 1 {
				// check if the output path exists and is a directory
				return utils.CreateDirectoryIfNotExists(args[0])
			}

			return nil
		},
		Long: fmt.Sprintf(`Initialize creates a configuration file for scaffy.

It creates a local config file "%s" in the current directory.`, config.DEFAULT_CONFIG_FILE_NAME),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, force, err := buildConfigFromFlags(cmd)
			if err != nil {
				return err
			}

			return runConfigInitializerE(cmd, args, cfg, force)
		},
	}

	initCmd.Flags().Bool("force", false, "allow write operations that overwrite files")

	// config overrides
	initCmd.Flags().StringP("author", "a", defaultCfg.Author, "author name")
	initCmd.Flags().StringToString("lang", defaultCfg.Languages, "programming languages to generate (e.g. --lang ruby=rb,python=py)")
	initCmd.Flags().String("log-level", defaultCfg.LogLevel, "slog level (e.g., debug, info, warn, error)")
	initCmd.Flags().String("output-dir", defaultCfg.OutputDir, "directory to output generated code")
	initCmd.Flags().String("templates-dir", defaultCfg.TemplatesDir, "directory containing template files")

	return initCmd
}

func buildConfigFromFlags(cmd *cobra.Command) (*config.Config, bool, error) {
	author, err := cmd.Flags().GetString("author")
	if err != nil {
		return nil, false, err
	}

	languages, err := cmd.Flags().GetStringToString("lang")
	if err != nil {
		return nil, false, err
	}

	logLevel, err := cmd.Flags().GetString("log-level")
	if err != nil {
		return nil, false, err
	}

	outputDir, err := cmd.Flags().GetString("output-dir")
	if err != nil {
		return nil, false, err
	}

	templatesDir, err := cmd.Flags().GetString("templates-dir")
	if err != nil {
		return nil, false, err
	}

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return nil, false, err
	}

	return &config.Config{
		Author:       author,
		Languages:    languages,
		LogLevel:     logLevel,
		OutputDir:    outputDir,
		TemplatesDir: templatesDir,
	}, force, nil
}

// runConfigInitializerE executes the logic for the init command, creating a config file with default settings.
func runConfigInitializerE(cmd *cobra.Command, args []string, cfg *config.Config, force bool) error {
	ctx := cmd.Context()

	userUI := ui.FromContext(ctx)
	app := program.FromContext(ctx)

	configFilePath := config.DEFAULT_CONFIG_FILE_NAME

	if len(args) > 0 {
		configFilePath = filepath.Join(args[0], configFilePath)
	}

	overwrite := false

	if force {
		userUI.Warningf("existing config will be overwritten (--force)")

		overwrite = true
	} else {
		// check if config file already exists
		if _, err := os.Stat(configFilePath); err == nil {
			return fmt.Errorf("%w at %s (use --force to overwrite)", ErrConfigExists, configFilePath)
		}
	}

	if err := app.CreateNewConfig(cfg, configFilePath, overwrite); err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	userUI.Infof("initialized scaffy with default config at %s", configFilePath)

	return nil
}
