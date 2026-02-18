package cli_test

import (
	"bytes"
	"testing"

	"github.com/25prabhu10/up-skill/cmd/cli"
	"github.com/25prabhu10/up-skill/internal/utils"
	"github.com/25prabhu10/up-skill/pkg/build_info"
)

func TestRootCmd_HelpText(t *testing.T) {
	t.Parallel()

	buf, err := utils.ExecuteTestCommand(t, cli.GetRootCmd(), []string{"--help"})

	if err != nil {
		t.Fatalf("failed to execute root command: %v", err)
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("This CLI tool helps you scaffold files for different programming languages and frameworks.")) {
		t.Errorf("unexpected help output: %s", output)
	}
}

func TestGetRootCmd_BasicConfiguration(t *testing.T) {
	t.Parallel()

	cmd := cli.GetRootCmd()
	if cmd == nil {
		t.Fatal("expected root command to be initialized, got nil")
	}

	if cmd.Use != build_info.AppName {
		t.Errorf("unexpected command use: got %s, want %s", cmd.Use, build_info.AppName)
	}

	if cmd.Short != "A CLI program to scaffold files for different languages." {
		t.Errorf("unexpected command short description: got %s", cmd.Short)
	}

	if cmd.Long != "This CLI tool helps you scaffold files for different programming languages and frameworks. It supports multiple languages and provides a simple interface to generate boilerplate code for your projects." {
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

func TestRootCmd_SubcommandsRegistered(t *testing.T) {
	t.Parallel()

	cmd := cli.GetRootCmd()

	hasInit := false

	for _, c := range cmd.Commands() {
		switch c.Name() {
		case "init":
			hasInit = true
		}
	}

	if !hasInit {
		t.Fatal("expected init subcommand to be registered")
	}
}
