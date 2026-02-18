// Package commands provides the implementation of the 'llm' command, which generates documentation for the CLI in various formats suitable for language models. It includes validation for output path and format flags, and supports
package commands

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/25prabhu10/up-skill/internal/ui"
	"github.com/25prabhu10/up-skill/internal/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// Documentation formats supported by the llm command.
const (
	Markdown = "markdown"
	Man      = "man"
	Rest     = "rest"
)

// Flags for the llm command, including output path, format, and front matter option.
var (
	flagDocOutputPath string
	flagDocFormat     string
	flagFrontMatter   bool
)

// Default output path for generated documentation and allowed formats for validation.
var (
	defaultDocOutputPath = filepath.Join(".agents", "skills", "cli", "reference")
	allowedDocFormats    = []string{Markdown, Man, Rest}
	allowedDocFormatsStr = strings.Join(allowedDocFormats, "|")
)

// GetLlmCommand returns a Cobra command that generates documentation for the CLI in a format suitable for language models. It supports multiple output formats and includes validation for the output path and format flags.
func GetLlmCommand() *cobra.Command {
	llmCmd := &cobra.Command{
		Use:   "llm",
		Short: "Generate llm ready documentation of the CLI",
		Long:  `This command generates documentation for the CLI in a format that can be easily consumed by language models. It supports multiple output formats, including Markdown, Man, and ReStructuredText. The generated documentation includes details about each command, its flags, and usage examples.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if utils.IsStringEmpty(flagDocOutputPath) {
				return fmt.Errorf("'%s' is not a valid output path. Please choose a different output path", flagDocOutputPath)
			}

			if !slices.Contains(allowedDocFormats, flagDocFormat) {
				return fmt.Errorf("'%s' is not a valid format. Please choose a different format (%s)", flagDocFormat, allowedDocFormatsStr)
			}

			return nil
		},
		RunE: runLlmCommandE,
	}

	// register flags
	llmCmd.Flags().StringVarP(&flagDocOutputPath, "output", "o", defaultDocOutputPath, "output path for the generated docs")
	llmCmd.Flags().StringVarP(&flagDocFormat, "format", "f", Markdown, fmt.Sprintf("format to generate docs (%s)", allowedDocFormatsStr))
	llmCmd.Flags().BoolVar(&flagFrontMatter, "front-matter", false, "add frontmatter to markdown files")

	return llmCmd
}

// runLlmCommandE executes the logic for the llm command, generating documentation for the CLI in the specified format and output path.
func runLlmCommandE(cmd *cobra.Command, args []string) error {
	if err := utils.CreateDirectoryIfNotExists(flagDocOutputPath); err != nil {
		return err
	}

	ctx := cmd.Context()

	userUI := ui.FromContext(ctx)

	root := cmd.Root()
	root.DisableAutoGenTag = true

	switch flagDocFormat {
	case "markdown":
		if flagFrontMatter {
			prep := func(filename string) string {
				base := filepath.Base(filename)
				name := strings.TrimSuffix(base, filepath.Ext(base))
				title := strings.ReplaceAll(name, "_", " ")
				return fmt.Sprintf("---\ntitle: %q\nslug: %q\ndescription: \"CLI reference for %s\"\n---\n\n", title, name, title)
			}
			link := func(name string) string { return strings.ToLower(name) }
			if err := doc.GenMarkdownTreeCustom(root, flagDocOutputPath, prep, link); err != nil {
				return fmt.Errorf("error in generating markdown docs with frontmatter: %w", err)
			}
		} else {
			if err := doc.GenMarkdownTree(root, flagDocOutputPath); err != nil {
				return fmt.Errorf("error in generating markdown docs: %w", err)
			}
		}
	case "man":
		hdr := &doc.GenManHeader{Title: strings.ToUpper(root.Name()), Section: "1"}
		if err := doc.GenManTree(root, hdr, flagDocOutputPath); err != nil {
			return fmt.Errorf("error in generating man pages: %w", err)
		}
	case "rest":
		if err := doc.GenReSTTree(root, flagDocOutputPath); err != nil {
			return fmt.Errorf("error in generating reStructuredText docs: %w", err)

		}
	default:
		return fmt.Errorf("unsupported doc format: %s", flagDocFormat)
	}

	userUI.Infof("generated documentation for CLI at %s in %s format", flagDocOutputPath, flagDocFormat)

	return nil
}

// func generateSKILLSReferenceDoc(cmd *cobra.Command) string {
// 	return ""
// }
