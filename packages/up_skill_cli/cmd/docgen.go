package cmd

import (
	"25prabhu10/up-skill/cli/internal/options"
	"25prabhu10/up-skill/cli/internal/ui/single_select"
	textinput "25prabhu10/up-skill/cli/internal/ui/text_input"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	flagDocOutputPath string
	flagDocFormat     string
	flagFrontMatter   bool
)

type DocGenOptions struct {
	DocOutputPath *textinput.Output
	DocFormat     *single_select.Output
	FrontMatter   *single_select.Output
}

var docgenCmd = &cobra.Command{
	Use:   "docgen",
	Short: "A brief description of your command",
	Long:  ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		flagDocOutputPath = strings.TrimSpace(flagDocOutputPath)

		if !flagInteractive && flagDocOutputPath == "" {
			err = fmt.Errorf("'%s' is not a valid output path. Please choose a different output path", flagDocOutputPath)
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}

		if !slices.Contains(options.AllowedDocFormats, flagDocFormat) {
			err = fmt.Errorf("'%s' is not a valid format. Please choose a different format (%s)", flagDocFormat, options.AllowedDocFormatsStr())
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		docGenOptions := DocGenOptions{
			DocOutputPath: &textinput.Output{},
			DocFormat:     &single_select.Output{},
			FrontMatter:   &single_select.Output{},
		}

		if flagInteractive {
			tea_program := tea.NewProgram(textinput.InitialTextInputModelWithPlaceholder(docGenOptions.DocOutputPath, "Where to put the docs?", flagDocOutputPath, flagDocOutputPath))

			if _, err := tea_program.Run(); err != nil {
				log.Printf("Path contains an error: %v", err)
				cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
			}

			if docGenOptions.DocOutputPath.Output != "" {
				flagDocOutputPath = docGenOptions.DocOutputPath.Output
			}

			var formatOptions []list.Item = []list.Item{
				options.NewItem(options.Markdown, "Generate docs in Markdown format"),
				options.NewItem(options.Man, "Generate docs in Man format"),
				options.NewItem(options.Rest, "Generate docs in ReStructuredText format"),
			}
			var formatTitle = "Select the format to generate docs"

			tea_program = tea.NewProgram(single_select.InitializeModel(&formatOptions, &formatTitle, docGenOptions.DocFormat))
			if _, err := tea_program.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
			}
			if docGenOptions.DocFormat.Output != "" {
				flagDocFormat = docGenOptions.DocFormat.Output
			}

			if flagDocFormat == string(options.Markdown) {
				var frontMatterOptions []list.Item = []list.Item{
					options.NewItem(options.Yes, "Add frontmatter to markdown files"),
					options.NewItem(options.No, "Do not add frontmatter to markdown files"),
				}
				var frontMatterTitle = "Do you want to add frontmatter to markdown files?"

				tea_program = tea.NewProgram(single_select.InitializeModel(&frontMatterOptions, &frontMatterTitle, docGenOptions.FrontMatter))
				if _, err := tea_program.Run(); err != nil {
					fmt.Printf("Alas, there's been an error: %v", err)
					cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
				}
				if docGenOptions.FrontMatter.Output != "" {
					flagFrontMatter = docGenOptions.FrontMatter.Output == string(options.Yes)
				}
			}
		}

		if err := os.MkdirAll(flagDocOutputPath, 0o755); err != nil {
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}

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
					cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
				}
			} else {
				if err := doc.GenMarkdownTree(root, flagDocOutputPath); err != nil {
					cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
				}
			}
		case "man":
			hdr := &doc.GenManHeader{Title: strings.ToUpper(root.Name()), Section: "1"}
			if err := doc.GenManTree(root, hdr, flagDocOutputPath); err != nil {
				cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
			}
		case "rest":
			if err := doc.GenReSTTree(root, flagDocOutputPath); err != nil {
				cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
			}
		default:
			cobra.CheckErr(textinput.CreateErrorInputModel(fmt.Errorf("unknown format: %s", flagDocFormat)).Err())
		}

		fmt.Printf("Docs Generated at o:%s --> f:%s --> M:%t - Global: %t", flagDocOutputPath, flagDocFormat, flagFrontMatter, flagInteractive)
	},
}

func init() {
	rootCmd.AddCommand(docgenCmd)

	docgenCmd.Flags().StringVarP(&flagDocOutputPath, "out", "o", "./docs/cli", "output path for the generated docs")
	docgenCmd.Flags().StringVarP(&flagDocFormat, "format", "f", "markdown", fmt.Sprintf("format to generate docs (%s)", options.AllowedDocFormatsStr()))
	docgenCmd.Flags().BoolVar(&flagFrontMatter, "frontmatter", false, "add frontmatter to markdown files")
}
