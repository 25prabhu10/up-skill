package utils

import (
	"os"
	"runtime"
)

// OSInfo defines an interface for accessing OS information.
type OSInfo interface {
	GetOS() string
	GetUserConfigDir() (string, error)
}

type realOSInfo struct{}

// NewOSInfo creates an instance of the real OS info provider.
func NewOSInfo() *realOSInfo {
	return &realOSInfo{}
}

// GetOS returns the name of the operating system using runtime.GOOS.
func (r *realOSInfo) GetOS() string {
	return runtime.GOOS
}

// GetUserConfigDir returns the user configuration directory based on the operating system. It uses os.UserConfigDir() and falls back to a reasonable default if it encounters an error.
func (r *realOSInfo) GetUserConfigDir() (string, error) {
	return os.UserConfigDir()
}
