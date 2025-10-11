package cmd

import (
	"25prabhu10/up-skill/cli/internal/program"
	"25prabhu10/up-skill/cli/internal/ui/text_input"
	"25prabhu10/up-skill/cli/internal/utils"
	"fmt"
	"log"
	"strings"

	// "github.com/25prabhu10/up-skill/pkg/options"
	// multiSelect "github.com/25prabhu10/up-skill/cli/internal/ui/multi_select"
	// tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// var (
// 	logoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6"))
// 	tipMsgStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("190")).Italic(true)
// 	endingMsgStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170")).Bold(true)
// )

var (
	name          string
	flagOutput    string
	flagLanguage  string
	flagLanguages []string
)

type UserOptions struct {
	TopicName  *textinput.Output
	OutputPath *textinput.Output
	// Languages  *multiSelect.Selection
	// Language   *multiSelect.Selection
}

var createCmd = &cobra.Command{
	Use:   "create <topic-name>",
	Short: "Create files for the specified options.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		name = strings.TrimSpace(args[0])
		flagOutput = strings.TrimSpace(flagOutput)

		if name == "" || (name != "" && !utils.ValidateName(name)) {
			err = fmt.Errorf("'%s' is not a valid topic name. Please choose a different name (string length 1 to 255 chars)", name)
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}

		if !flagInteractive && flagOutput == "" {
			err = fmt.Errorf("'%s' is not a valid output path. Please choose a different output path", flagOutput)
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}

		fmt.Printf("Tags: %s --> %v\n", flagLanguage, flagLanguages)

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("standard logger")
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("with file/line")
		// var err error

		fmt.Printf("create called with args: %v\n", args)

		// userOptions := UserOptions{
		// 	TopicName:  &textinput.Output{},
		// 	OutputPath: &textinput.Output{},
		// 	// Languages:  &multiSelect.Selection{},
		// 	// Language:   &multiSelect.Selection{},
		// }

		topic := &program.Topic{
			Name:      name,
			Output:    flagOutput,
			Languages: flagLanguages,
		}

		// if flagInteractive {
		// 	// options := options.InitOptions(options.Language(flagLanguage))
		// 	tea_program := tea.NewProgram(textinput.InitialTextInputTopicModel(userOptions.TopicName, "What is the name of the topic?", topic, topic.Name))

		// 	if _, err := tea_program.Run(); err != nil {
		// 		log.Printf("Name of topic contains an error: %v", err)
		// 		cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		// 	}

		// 	if userOptions.TopicName.Output != "" && !utils.ValidateNameWithTrim(userOptions.TopicName.Output) {
		// 		err = fmt.Errorf("'%s' is not a valid topic name. Please choose a different name (string length 1 to 255 chars)", userOptions.TopicName.Output)
		// 		cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		// 	}

		// 	topic.Name = utils.NormalizeName(userOptions.TopicName.Output)

		// 	tea_program = tea.NewProgram(textinput.InitialTextInputTopicModel(userOptions.OutputPath, "What is the output directory?", topic, topic.Output))

		// 	if _, err := tea_program.Run(); err != nil {
		// 		log.Printf("Output path contains an error: %v", err)
		// 		cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		// 	}

		// 	// if userOptions.OutputPath.Output != "" {
		// 	// 	outputExists, err := utils.DoesDirectoryExistAndIsNotEmpty(userOptions.OutputPath.Output)

		// 	// 	if err != nil || !outputExists {
		// 	// 		err = fmt.Errorf("directory '%s' already exists and is not empty. Please choose a different output path", userOptions.OutputPath.Output)
		// 	// 		cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		// 	// 	}
		// 	// }

		// 	topic.Output = userOptions.OutputPath.Output

		// 	// option := options.Options["language"]
		// 	// tea_program = tea.NewProgram((multiSelect.InitialModelMultiSelect(option.Items, userOptions.Language, option.Headers, topic)))

		// 	// if _, err := tea_program.Run(); err != nil {
		// 	// 	log.Printf("Output path contains an error: %v", err)
		// 	// 	cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		// 	// }

		// 	// for key, opt := range userOptions.Language.Choices {
		// 	// 	// topic.Languages[strings.ToLower(key)] = opt
		// 	// 	// err := cmd.Flag("feature").Value.Set(strings.ToLower(key))
		// 	// 	fmt.Printf("Languagasasfas: %s <--> %t", key, opt)
		// 	// 	if err != nil {
		// 	// 		log.Fatal("failed to set the feature flag value", err)
		// 	// 	}
		// 	// }

		// } else {
		topic.Name = utils.NormalizeName(topic.Name)
		// }

		fmt.Printf("Topic: %+v", topic)

		fmt.Printf("Completed!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&flagOutput, "output", "o", "./src", "path of the output directory")
	createCmd.Flags().StringVarP(&flagLanguage, "language", "l", "", "") // fmt.Sprintf("languages (repeatable) [allowed: %s]", strings.Join(options.AllowedLanguageTypes, ", ")))
	// createCmd.Flags().StringSliceVarP(&flagLanguages, "languages", "ls", nil, fmt.Sprintf("languages (repeatable) [allowed: %s]", strings.Join(options.AllowedLanguageTypes, ", ")))
}
