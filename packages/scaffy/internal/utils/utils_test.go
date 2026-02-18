package utils_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/25prabhu10/scaffy/internal/utils"
)

func TestIsStringEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		inputStr string
		expected bool
	}{
		{
			name:     "string is not empty",
			inputStr: "hello world",
			expected: false,
		},
		{
			name:     "string empty",
			inputStr: "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := utils.IsStringEmpty(tt.inputStr)

			if tt.expected != result {
				t.Errorf("expected result %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestIsStringOverMaxLength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		inputStr string
		expected bool
	}{
		{
			name:     "string is under max length",
			inputStr: "hello",
			expected: false,
		},
		{
			name:     "string is exactly max length",
			inputStr: utils.GetLongStringChars(),
			expected: false,
		},
		{
			name:     "string is over max length",
			inputStr: utils.GetLongString256Chars(),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := utils.IsStringOverMaxLength(tt.inputStr)

			if tt.expected != result {
				t.Errorf("expected result %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestNormalizeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		inputStr string
		expected string
	}{
		{
			name:     "string with spaces and special characters",
			inputStr: "  Hello, World!  ",
			expected: "hello_world",
		},
		{
			name:     "string with multiple special characters",
			inputStr: "Go@Lang#2024-",
			expected: "go_lang_2024",
		},
		{
			name:     "string with leading/trailing underscores",
			inputStr: "__Hello__World__",
			expected: "hello_world",
		},
		{
			name:     "string with only special characters",
			inputStr: "@#$%^&*()",
			expected: "",
		},
		{
			name:     "string with no special characters",
			inputStr: "HelloWorld",
			expected: "helloworld",
		},
		{
			name:     "string with spaces only",
			inputStr: "     ",
			expected: "",
		},
		{
			name:     "empty string",
			inputStr: "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := utils.NormalizeString(tt.inputStr)

			if tt.expected != result {
				t.Errorf("expected result '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestCreateDirectoryIfNotExists(t *testing.T) { //nolint:tparallel,paralleltest // `t.Chdir` does not support
	tempDir := t.TempDir()
	t.Chdir(tempDir)

	// create a test file
	testFilePath := "testfile.txt"
	if err := os.WriteFile(testFilePath, []byte("test content"), 0600); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	tests := []struct {
		name      string
		dirPath   string
		expectErr bool
	}{
		{
			name:      "directory already exists",
			dirPath:   tempDir,
			expectErr: false,
		},
		{
			name:      "directory does not exist",
			dirPath:   filepath.Join(tempDir, "testdata/newdir"),
			expectErr: false,
		},
		{
			name:      "path exists but is a file",
			dirPath:   testFilePath,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := utils.CreateDirectoryIfNotExists(tt.dirPath)

			if tt.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			} else if !tt.expectErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
