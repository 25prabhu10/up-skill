package test_utils

import (
	"bytes"
	"io/fs"
	"os"

	"github.com/25prabhu10/scaffy/internal/utils"
)

type MockOSInfo struct {
	MockGOOS          string
	MockUserConfigDir string
}

func (m *MockOSInfo) GetOS() string {
	return m.MockGOOS
}

func (m *MockOSInfo) GetUserConfigDir() (string, error) {
	if m.MockUserConfigDir != "" {
		return m.MockUserConfigDir, nil
	}

	return "", os.ErrNotExist
}

type MockFileSystem struct {
	utils.FileSystem

	StatFunc      func(name string) (fs.FileInfo, error)
	MkdirAllFunc  func(path string, perm fs.FileMode) error
	WriteFileFunc func(name string, data []byte, perm fs.FileMode) error
	RemoveFunc    func(name string) error
	RemoveAllFunc func(path string) error
	ReadDirFunc   func(name string) ([]fs.DirEntry, error)
}

func (m *MockFileSystem) Stat(name string) (fs.FileInfo, error) {
	if m.StatFunc != nil {
		return m.StatFunc(name)
	}

	return m.FileSystem.Stat(name)
}

func (m *MockFileSystem) MkdirAll(path string, perm fs.FileMode) error {
	if m.MkdirAllFunc != nil {
		return m.MkdirAllFunc(path, perm)
	}

	return m.FileSystem.MkdirAll(path, perm)
}

func (m *MockFileSystem) WriteFile(name string, data []byte, perm fs.FileMode) error {
	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(name, data, perm)
	}

	return m.FileSystem.WriteFile(name, data, perm)
}

func (m *MockFileSystem) Remove(name string) error {
	if m.RemoveFunc != nil {
		return m.RemoveFunc(name)
	}

	return m.FileSystem.Remove(name)
}

func (m *MockFileSystem) RemoveAll(path string) error {
	if m.RemoveAllFunc != nil {
		return m.RemoveAllFunc(path)
	}

	return m.FileSystem.RemoveAll(path)
}

func (m *MockFileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	if m.ReadDirFunc != nil {
		return m.ReadDirFunc(name)
	}

	return m.FileSystem.ReadDir(name)
}

type MockTemplateManager struct {
	RenderTemplateFunc func(lang string) (*bytes.Buffer, error)
}

func (m *MockTemplateManager) RenderTemplate(lang string) (*bytes.Buffer, error) {
	if m.RenderTemplateFunc != nil {
		return m.RenderTemplateFunc(lang)
	}

	return bytes.NewBufferString("mock content"), nil
}
