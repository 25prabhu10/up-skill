package cli_test

import (
	"strings"
	"testing"

	"github.com/25prabhu10/scaffy/cmd/cli"
	"github.com/25prabhu10/scaffy/internal/utils"
	"github.com/25prabhu10/scaffy/pkg/build_info"
)

func TestRootCmd_HelpText(t *testing.T) {
	t.Parallel()

	buf, err := utils.ExecuteTestCommand(t, cli.GetRootCmd(), []string{"--help"})
	if err != nil {
		t.Fatalf("failed to execute root command: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "This CLI tool helps you scaffold files for different programming languages and") {
		t.Errorf("unexpected help output: %s", output)
	}
}

func TestGetRootCmd_BasicConfiguration(t *testing.T) {
	t.Parallel()

	cmd := cli.GetRootCmd()
	if cmd == nil {
		t.Fatal("expected root command to be initialized, got nil")
	}

	if cmd.Use != build_info.APP_NAME {
		t.Errorf("unexpected command use: got %s, want %s", cmd.Use, build_info.APP_NAME)
	}

	if cmd.Short != "A CLI program to scaffold files for different languages." {
		t.Errorf("unexpected command short description: got %s", cmd.Short)
	}

	if !strings.Contains(cmd.Long, "This CLI tool helps you scaffold files for different programming languages") {
		t.Errorf("unexpected command long description: got %s", cmd.Long)
	}
}

func TestRootCmd_PersistentFlagsRegistered(t *testing.T) {
	t.Parallel()

	cmd := cli.GetRootCmd()
	pFlags := cmd.PersistentFlags()

	if pFlags.Lookup("verbose") == nil {
		t.Fatal("expected --verbose persistent flag to be registered")
	}

	if pFlags.Lookup("quiet") == nil {
		t.Fatal("expected --quiet persistent flag to be registered")
	}

	if pFlags.Lookup("config") == nil {
		t.Fatal("expected --config persistent flag to be registered")
	}
}

func TestRootCmd_SubcommandsRegistered(t *testing.T) { //nolint:paralleltest // subcommand registration is not parallelizable
	cmd := cli.GetRootCmd()

	subcommands := []string{"init"}
	registeredCommands := make(map[string]bool)

	for _, c := range cmd.Commands() {
		registeredCommands[c.Name()] = true
	}

	for _, required := range subcommands {
		if !registeredCommands[required] {
			t.Fatalf("expected %q subcommand to be registered", required)
		}
	}
}
