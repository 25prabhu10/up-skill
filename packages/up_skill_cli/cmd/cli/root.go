package cli

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/25prabhu10/up-skill/internal/config"
	"github.com/25prabhu10/up-skill/pkg/buildinfo"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	quiet   bool
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:     buildinfo.AppName,
	Short:   "A CLI program to scaffold files for different languages.",
	Long:    `This CLI tool helps you scaffold files for different programming languages and frameworks. `,
	Version: buildinfo.Version,
	Example: `  ` + buildinfo.AppName + ` "Inverse Binary Tree"`,
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	// Initialize configuration after cobra is initialized
	cobra.OnInitialize(initConfig)

	// Persistent flags
	configMsg := fmt.Sprintf("config file (default %s,%s)", config.DEFAULT_CONFIG_FILE_NAME, config.GetDefaultConfigPath())
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", configMsg)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "suppress non-error output")
}

func initConfig() {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Set log level based on config and flags
	if quiet {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})))
	} else if verbose {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	} else {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: cfg.LogLevel})))
	}

	slog.Info("Configuration loaded", "file", cfgFile)
}
