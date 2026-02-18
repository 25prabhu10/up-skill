package utils

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/25prabhu10/up-skill/internal/constants"
	"github.com/25prabhu10/up-skill/internal/ui"
	"github.com/spf13/cobra"
)

func ExecuteTestCommand(t *testing.T, cmd *cobra.Command, args []string) (*bytes.Buffer, error) {
	t.Helper()

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	cmd.SilenceUsage = true

	// clean up after test to reset command state.
	t.Cleanup(func() {
		cmd.SetArgs(nil)
		cmd.SetOut(nil)
		cmd.SetErr(nil)
	})

	return buf, cmd.Execute()
}

func ExecuteTestCommandWithContext(t *testing.T,
	cmd *cobra.Command,
	args []string,
	verbose bool,
	quiet bool) (*bytes.Buffer, *bytes.Buffer, error) {
	t.Helper()

	// create command with context.
	ctx := context.Background()

	// All UI output goes to stderr (POSIX compliance).
	var uiStderr bytes.Buffer

	var userUI *ui.UI

	if quiet {
		userUI = ui.New(ui.WithOutput(&uiStderr), ui.WithQuiet(quiet))
	} else {
		userUI = ui.New(ui.WithOutput(&uiStderr))
	}

	ctx = ui.WithUI(ctx, userUI)
	cmd.SetContext(ctx)

	// capture command output to a buffer instead of stdout/stderr.
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	cmd.SilenceUsage = true

	// clean up after test to reset command state.
	t.Cleanup(func() {
		cmd.SetArgs(nil)
		cmd.SetOut(nil)
		cmd.SetErr(nil)
	})

	return &uiStderr, buf, cmd.Execute()
}

func AssertPanics(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic, got none")
		}
	}()

	fn()
}

func GetLongStringChars() string {
	return strings.Repeat("A", constants.MAX_NAME_LENGTH)
}

func GetLongString256Chars() string {
	return strings.Repeat("a", constants.MAX_NAME_LENGTH+1)
}
