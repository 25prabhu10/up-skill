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

// GetLlmCommand returns a Cobra command that generates documentation for the CLI in a format suitable for language models. It supports multiple output formats and includes validation for the output path and format flags.
func GetLlmCommand() *cobra.Command {
	var (
		docOutputPath string
		docFormat     string
		frontMatter   bool
	)

	defaultDocOutputPath := filepath.Join(".agents", "skills", "cli", "reference")

	llmCmd := &cobra.Command{
		Use:   "llm",
		Short: "Generate llm ready documentation of the CLI",
		Long: `
This command generates documentation for the CLI in a format that can be easily
consumed by language models. It supports multiple output formats, including
Markdown, Man, and ReStructuredText. The generated documentation includes
details about each command, its flags, and usage examples.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if utils.IsStringEmpty(docOutputPath) {
				return fmt.Errorf("%w: %q", ErrInvalidOutputPath, docOutputPath)
			}

			return validateDocFormat(docFormat)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLlmCommandE(cmd, docOutputPath, docFormat, frontMatter)
		},
	}

	// register flags
	llmCmd.Flags().StringVarP(&docOutputPath, "output", "o", defaultDocOutputPath, "output path for the generated docs")
	llmCmd.Flags().StringVarP(&docFormat, "format", "f", program.Markdown, fmt.Sprintf("format to generate docs (%s)", program.AllowedDocFormatsStr))
	llmCmd.Flags().BoolVar(&frontMatter, "front-matter", false, "add frontmatter to markdown files")

	return llmCmd
}

// runLlmCommandE executes the logic for the llm command, generating documentation for the CLI in the specified format and output path.
func runLlmCommandE(cmd *cobra.Command, docOutputPath, docFormat string, frontMatter bool) error {
	ctx := cmd.Context()

	userUI := ui.FromContext(ctx)
	app := program.FromContext(ctx)

	if err := app.GenerateLLMDocs(cmd.Root(), docOutputPath, docFormat, frontMatter); err != nil {
		return err
	}

	userUI.Infof("generated documentation for CLI at %s in %s format", docOutputPath, docFormat)

	return nil
}

func validateDocFormat(docFormat string) error {
	if !slices.Contains(program.AllowedDocFormats, docFormat) {
		return fmt.Errorf("%w: %q (valid: %s)", program.ErrInvalidDocFormat, docFormat, program.AllowedDocFormatsStr)
	}

	return nil
}
