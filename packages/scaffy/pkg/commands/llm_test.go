package commands_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/25prabhu10/scaffy/internal/program"
	"github.com/25prabhu10/scaffy/internal/utils/test_utils"
	"github.com/25prabhu10/scaffy/pkg/commands"
)

func TestLlmCmd_Success(t *testing.T) { //nolint:paralleltest // t.Setenv used and global flags
	tests := []struct {
		name              string
		args              []string
		expectedPathParts []string
		expectedLog       string
	}{
		{
			name:              "default args",
			args:              []string{},
			expectedPathParts: []string{".agents", "skills", "cli", "reference", "llm.md"},
			expectedLog:       "generated documentation for CLI at",
		},
		{
			name:              "custom output dir",
			args:              []string{"--output", "custom_docs"},
			expectedPathParts: []string{"custom_docs", "llm.md"},
			expectedLog:       "generated documentation for CLI at custom_docs in markdown format",
		},
		{
			name:              "man format",
			args:              []string{"--output", "man_docs", "-f", program.Man},
			expectedPathParts: []string{"man_docs", "llm.1"},
			expectedLog:       "generated documentation for CLI at man_docs in man format",
		},
		{
			name:              "rest format",
			args:              []string{"-o", "rest_docs", "--format", program.Rest},
			expectedPathParts: []string{"rest_docs", "llm.rst"},
			expectedLog:       "generated documentation for CLI at rest_docs in rest format",
		},
	}

	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := test_utils.SetupTestEnv(t)
			cmd := commands.GetLlmCommand()

			uiStderr, _, err := test_utils.ExecuteTestCommandWithContext(t, cmd, tt.args, false, false)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			errOut := uiStderr.String()
			if !strings.Contains(errOut, tt.expectedLog) {
				t.Errorf("expected stderr to contain %q, got %q", tt.expectedLog, errOut)
			}

			expectedPath := filepath.Join(append([]string{tmpDir}, tt.expectedPathParts...)...)
			if _, err := os.Stat(expectedPath); err != nil {
				t.Errorf("expected file at %s, got error: %v", expectedPath, err)
			}
		})
	}
}

func TestLlmCmd_FrontMatter(t *testing.T) { //nolint:paralleltest // t.Setenv used and global flags
	tmpDir := test_utils.SetupTestEnv(t)
	cmd := commands.GetLlmCommand()

	args := []string{"--output", "fm_docs", "--front-matter"}

	_, _, err := test_utils.ExecuteTestCommandWithContext(t, cmd, args, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "fm_docs", "llm.md")

	content, err := os.ReadFile(expectedPath) //nolint:gosec // test file
	if err != nil {
		t.Fatalf("expected markdown file at %s, got error: %v", expectedPath, err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "---") || !strings.Contains(contentStr, "title:") {
		t.Errorf("expected front-matter in markdown file, got: %s", contentStr)
	}
}

func TestLlmCmd_Errors(t *testing.T) { //nolint:paralleltest // t.Setenv used and global flags
	tests := []struct {
		name        string
		args        []string
		setupFile   string
		expectedErr string
	}{
		{
			name:        "empty output path",
			args:        []string{"--output", ""},
			expectedErr: "invalid output path",
		},
		{
			name:        "invalid format",
			args:        []string{"--output", "docs", "--format", "invalid"},
			expectedErr: "invalid documentation format",
		},
		{
			name:        "output dir is a file",
			args:        []string{"--output", "file_as_dir"},
			setupFile:   "file_as_dir",
			expectedErr: "failed to create output directory",
		},
	}

	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := test_utils.SetupTestEnv(t)

			if tt.setupFile != "" {
				filePath := filepath.Join(tmpDir, tt.setupFile)
				if err := os.WriteFile(filePath, []byte("dummy"), 0600); err != nil {
					t.Fatalf("failed to create setup file: %v", err)
				}
			}

			cmd := commands.GetLlmCommand()

			_, _, err := test_utils.ExecuteTestCommandWithContext(t, cmd, tt.args, false, false)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("expected error to contain %q, got %q", tt.expectedErr, err.Error())
			}
		})
	}
}
