// Package cli implements the CLI commands for up-skill using Cobra.
//
// It provides the root command and all subcommands with proper flag handling,
// configuration management, and logging.
package cli

import (
	"errors"
	"fmt"

	"github.com/25prabhu10/up-skill/internal/config"
	"github.com/25prabhu10/up-skill/internal/program"
	"github.com/25prabhu10/up-skill/internal/ui"
	"github.com/25prabhu10/up-skill/pkg/build_info"
	"github.com/25prabhu10/up-skill/pkg/commands"

	"github.com/spf13/cobra"
)

// root command errors.
var (
	ErrVerboseQuietConflict = errors.New("verbose and quiet flags are mutually exclusive")
)

// variables to hold flag values.
var (
	flagCfgFile string
	flagVerbose bool
	flagQuiet   bool
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:               build_info.AppName,
	Short:             "A CLI program to scaffold files for different languages.",
	Long:              `This CLI tool helps you scaffold files for different programming languages and frameworks. It supports multiple languages and provides a simple interface to generate boilerplate code for your projects.`,
	Version:           fmt.Sprintf("%s (%s) built at %s", build_info.Version, build_info.GitCommit, build_info.BuildDate),
	SilenceUsage:      true,
	PersistentPreRunE: initApp,
	Example:           `  ` + build_info.AppName + ` "Inverse Binary Tree"`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

// GetRootCmd returns the root Cobra command.
func GetRootCmd() *cobra.Command {
	return rootCmd
}

// init function is called automatically when the package is imported. It sets up the Cobra command and flags.
func init() { //nolint:gochecknoinits // init is used to set up the root command
	// Persistent flags
	configMsg := fmt.Sprintf("config file (default %s or %s)", config.DEFAULT_CONFIG_FILE_NAME, config.GetDefaultConfigPath())
	rootCmd.PersistentFlags().StringVar(&flagCfgFile, "config", "", configMsg)

	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&flagQuiet, "quiet", "q", false, "suppress non-error output")

	// register subcommands
	rootCmd.AddCommand(commands.GetLlmCommand())
	rootCmd.AddCommand(commands.NewInitCmd())
}

// initApp initializes the application context and sets it in the root command. This allows all subcommands to access the program instance and configuration.
func initApp(cmd *cobra.Command, args []string) error {
	if flagVerbose && flagQuiet {
		return fmt.Errorf("%w", ErrVerboseQuietConflict)
	}

	app, err := program.New(flagCfgFile, flagVerbose, flagQuiet)
	cobra.CheckErr(err)

	// create user interface (writes to stdout for user-facing messages).
	userUI := ui.New(ui.WithQuiet(flagQuiet))

	// store data in context for access in subcommands
	ctx := cmd.Context()

	ctx = program.WithProgram(ctx, app)
	ctx = ui.WithUI(ctx, userUI)

	// set the context in the root command so that it can be accessed by subcommands
	cmd.SetContext(ctx)

	return nil
}
