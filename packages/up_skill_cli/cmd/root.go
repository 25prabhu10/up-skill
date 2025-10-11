package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var flagInteractive bool

var rootCmd = &cobra.Command{
	Use:     "up-skill",
	Short:   "A CLI program to scaffold files for different languages.",
	Long:    `This CLI tool helps you scaffold files for different programming languages and frameworks. `,
	Version: "0.0.1",
	Example: "up-skill \"Inverse Binary Tree\"",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&flagInteractive, "interactive", "i", false, "launch interactive mode")
}
