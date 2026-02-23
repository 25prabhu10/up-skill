package utils

import (
	"io/fs"
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

type FileSystem interface {
	WriteFile(name string, data []byte, perm fs.FileMode) error
	Remove(name string) error
	RemoveAll(path string) error
	ReadDir(name string) ([]fs.DirEntry, error)
	Stat(name string) (fs.FileInfo, error)
	MkdirAll(path string, perm fs.FileMode) error
}

type realFileSystem struct{}

// NewFileSystem creates an instance of the real file system provider.
func NewFileSystem() *realFileSystem {
	return &realFileSystem{}
}

// WriteFile writes data to a file with the specified name and permissions.
func (r *realFileSystem) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

// Remove deletes the named file or directory.
func (r *realFileSystem) Remove(name string) error {
	return os.Remove(name)
}

// RemoveAll removes path and any children it contains.
func (r *realFileSystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// ReadDir reads the directory named by name and returns a list of directory entries.
func (r *realFileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

// Stat returns a FileInfo describing the named file or directory.
func (r *realFileSystem) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

// MkdirAll creates a directory named path, along with any necessary parents, and returns nil, or else returns an error.
func (r *realFileSystem) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}
