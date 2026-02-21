// Package commands provides the implementation of the 'llm' command, which generates documentation for the CLI in various formats suitable for language models. It includes validation for output path and format flags, and supports
package commands

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/25prabhu10/scaffy/internal/program"
	"github.com/25prabhu10/scaffy/internal/ui"
	"github.com/25prabhu10/scaffy/internal/utils"

	"github.com/spf13/cobra"
)

// llm command errors.
var (
	ErrInvalidOutputPath = errors.New("error: invalid output path. Please choose a different output path")
)

// Flags for the llm command, including output path, format, and front matter option.
var (
	flagDocOutputPath string
	flagDocFormat     string
	flagFrontMatter   bool
)

// GetLlmCommand returns a Cobra command that generates documentation for the CLI in a format suitable for language models. It supports multiple output formats and includes validation for the output path and format flags.
func GetLlmCommand() *cobra.Command {
	defaultDocOutputPath := filepath.Join(".agents", "skills", "cli", "reference")

	llmCmd := &cobra.Command{
		Use:   "llm",
		Short: "Generate llm ready documentation of the CLI",
		Long:  `This command generates documentation for the CLI in a format that can be easily consumed by language models. It supports multiple output formats, including Markdown, Man, and ReStructuredText. The generated documentation includes details about each command, its flags, and usage examples.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if utils.IsStringEmpty(flagDocOutputPath) {
				return fmt.Errorf("%w: %q", ErrInvalidOutputPath, flagDocOutputPath)
			}

			return validateDocFormat()
		},
		RunE: runLlmCommandE,
	}

	// register flags
	llmCmd.Flags().StringVarP(&flagDocOutputPath, "output", "o", defaultDocOutputPath, "output path for the generated docs")
	llmCmd.Flags().StringVarP(&flagDocFormat, "format", "f", program.Markdown, fmt.Sprintf("format to generate docs (%s)", program.AllowedDocFormatsStr))
	llmCmd.Flags().BoolVar(&flagFrontMatter, "front-matter", false, "add frontmatter to markdown files")

	return llmCmd
}

// runLlmCommandE executes the logic for the llm command, generating documentation for the CLI in the specified format and output path.
func runLlmCommandE(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	userUI := ui.FromContext(ctx)
	app := program.FromContext(ctx)

	if err := app.GenerateLLMDocs(cmd.Root(), flagDocOutputPath, flagDocFormat, flagFrontMatter); err != nil {
		return err
	}

	userUI.Infof("generated documentation for CLI at %s in %s format", flagDocOutputPath, flagDocFormat)

	return nil
}

func validateDocFormat() error {
	if !slices.Contains(program.AllowedDocFormats, flagDocFormat) {
		return fmt.Errorf("%w: %q (valid: %s)", program.ErrInvalidDocFormat, flagDocFormat, program.AllowedDocFormatsStr)
	}

	return nil
}
