package commands

import (
	"errors"
	"fmt"

	"github.com/25prabhu10/scaffy/internal/boilerplate"
	"github.com/25prabhu10/scaffy/internal/config"
	"github.com/25prabhu10/scaffy/internal/constants"
	"github.com/25prabhu10/scaffy/internal/program"
	"github.com/25prabhu10/scaffy/internal/ui"
	"github.com/25prabhu10/scaffy/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrInvalidProjectName = errors.New("project name cannot be empty")
	ErrProjectNameTooLong = errors.New("project name exceeds maximum length")
)

func GetNewCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "new [project-name]",
		Short: "Create a new project from a template",
		Long: `
The new command creates a new project by generating code from templates based on
the provided configuration.`,
		Example: `  scaffy new "Inverse Binary Search Tree"
  scaffy new "Linear Search" -o "my-search-algorithms"
  scaffy new "Graph Algorithms" --lang ruby=rb,python=py,go=go`,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			}

			return validateProjectName(args[0])
		},
		RunE: runNewCommandE,
	}

	newCmd.Flags().Bool("force", false, "allow write operations that overwrite files")

	newCmd.Flags().StringToStringP("lang", "l", config.GetDefaultLanguages(),
		"programming languages to generate")
	cobra.CheckErr(viper.BindPFlag("languages", newCmd.Flags().Lookup("lang")))

	newCmd.Flags().StringP("author", "a", "", "author name")
	cobra.CheckErr(viper.BindPFlag("author", newCmd.Flags().Lookup("author")))

	newCmd.Flags().StringP("output-dir", "o", "", "output directory")
	cobra.CheckErr(viper.BindPFlag("output-dir", newCmd.Flags().Lookup("output-dir")))

	newCmd.Flags().String("templates-dir", "", "user defined templates directory")
	cobra.CheckErr(viper.BindPFlag("templates-dir", newCmd.Flags().Lookup("templates-dir")))

	return newCmd
}

func runNewCommandE(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	userUI := ui.FromContext(ctx)
	app := program.FromContext(ctx)

	name := utils.NormalizeString(args[0])

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return err
	}

	if force {
		userUI.Warnf("Force flag is enabled. Existing files may be overwritten.")
	}

	if err := app.GenerateFilesFromTemplates(ctx, name, force); err != nil {
		if errors.Is(err, boilerplate.ErrDirectoryAlreadyExists) {
			userUI.Errorf("Directory for project %q already exists. Use --force to overwrite.", name)
			return nil
		}

		return err
	}

	userUI.Infof("Project %q created successfully!", name)

	return nil
}

func validateProjectName(name string) error {
	if utils.IsStringEmpty(name) {
		return ErrInvalidProjectName
	} else if utils.IsStringOverMaxLength(name) {
		return fmt.Errorf("%w: max length is %d characters", ErrProjectNameTooLong, constants.MAX_NAME_LENGTH)
	}

	return nil
}
