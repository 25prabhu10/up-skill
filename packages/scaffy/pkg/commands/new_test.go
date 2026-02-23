package commands_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/25prabhu10/scaffy/internal/utils/test_utils"
	"github.com/25prabhu10/scaffy/pkg/commands"
)

func TestNewCmd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		args        []string
		setupFunc   func(t *testing.T, tmpDir string)
		wantErr     bool
		errContains string
		verifyFunc  func(t *testing.T, tmpDir string)
	}{
		{
			name:    "creates project successfully",
			args:    []string{"my-project"},
			wantErr: false,
			verifyFunc: func(t *testing.T, tmpDir string) {
				t.Helper()

				expectedDir := filepath.Join(tmpDir, "my_project")
				if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
					t.Errorf("expected project directory not created")
				}
			},
		},
		{
			name:        "fails with empty project name",
			args:        []string{""},
			wantErr:     true,
			errContains: "cannot be empty",
		},
		{
			name:        "fails with project name too long",
			args:        []string{test_utils.GetLongString256Chars()},
			wantErr:     true,
			errContains: "exceeds maximum length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := test_utils.SetupTestEnv(t)

			if tt.setupFunc != nil {
				tt.setupFunc(t, tmpDir)
			}

			cmd := commands.GetNewCmd()
			_, _, err := test_utils.ExecuteTestCommandWithContext(t, cmd, tt.args, false, false)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.verifyFunc != nil {
				tt.verifyFunc(t, tmpDir)
			}
		})
	}
}
