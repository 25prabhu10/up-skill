package utils

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/25prabhu10/scaffy/internal/constants"
)

// utils errors.
var (
	ErrPathIsNotDir = errors.New("path is not a directory")
)

// IsStringEmpty checks if a string is empty or contains only whitespace.
func IsStringEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// IsStringOverMaxLength checks if a string exceeds the maximum allowed length.
func IsStringOverMaxLength(str string) bool {
	return utf8.RuneCountInString(str) > constants.MAX_NAME_LENGTH
}

// NormalizeString normalizes a string by trimming whitespace, replacing non-alphanumeric characters with underscores, and converting to lowercase.
func NormalizeString(str string) string {
	str = strings.TrimSpace(str)

	// Replace non-alphanumeric with underscores
	r := regexp.MustCompile(`\W`)
	str = r.ReplaceAllString(str, "_")

	// Remove leading/trailing underscores
	str = strings.Trim(str, "_")
	str = strings.ToLower(str)

	// Remove consecutive underscores
	for strings.Contains(str, "__") {
		str = strings.ReplaceAll(str, "__", "_")
	}

	return str
}

func GetCurrentDate() string {
	return time.Now().Format("02-01-2006")
}

func IsDirectoryExists(dirPath string, fs FileSystem) (bool, error) {
	info, err := fs.Stat(dirPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, fmt.Errorf("failed to check directory: %w", err)
	}

	return info.IsDir(), nil
}

// CreateDirectoryIfNotExists checks if a directory exists at the given path. If it does not exist, it creates the directory. If a file exists at the path, it returns an error.
func CreateDirectoryIfNotExists(dirPath string, fs FileSystem) error {
	// check if the output path exists and is a directory
	if info, err := fs.Stat(dirPath); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("%w: %s", ErrPathIsNotDir, dirPath)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		// create the directory if it does not exist
		return CreateDirectory(dirPath, fs)
	} else {
		return fmt.Errorf("failed to access directory: %w", err)
	}

	return nil
}

func CreateDirectory(dirPath string, fs FileSystem) error {
	if err := fs.MkdirAll(dirPath, 0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}
