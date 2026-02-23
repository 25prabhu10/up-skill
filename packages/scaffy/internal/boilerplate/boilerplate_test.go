package boilerplate_test

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/25prabhu10/scaffy/internal/boilerplate"
	"github.com/25prabhu10/scaffy/internal/logger"
	"github.com/25prabhu10/scaffy/internal/utils"
	"github.com/25prabhu10/scaffy/internal/utils/test_utils"
)

var errMock = errors.New("mock error")

func TestScaffold(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		opts        boilerplate.Options
		setupFS     func(t *testing.T, tmpDir string) *test_utils.MockFileSystem
		setupTM     func() *test_utils.MockTemplateManager
		setupCtx    func() (context.Context, context.CancelFunc)
		wantErr     bool
		errContains string
	}{
		{
			name: "successful scaffold",
			opts: boilerplate.Options{
				Name: "myapp",
				Languages: map[string]string{
					"go": "go",
				},
			},
			setupFS: func(t *testing.T, tmpDir string) *test_utils.MockFileSystem {
				t.Helper()

				return &test_utils.MockFileSystem{FileSystem: utils.NewFileSystem()}
			},
			setupTM: func() *test_utils.MockTemplateManager {
				return &test_utils.MockTemplateManager{}
			},
			setupCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			wantErr: false,
		},
		{
			name: "directory already exists without force",
			opts: boilerplate.Options{
				Name: "myapp",
				Languages: map[string]string{
					"go": "go",
				},
				Force: false,
			},
			setupFS: func(t *testing.T, tmpDir string) *test_utils.MockFileSystem {
				t.Helper()

				fs := utils.NewFileSystem()

				err := fs.MkdirAll(filepath.Join(tmpDir, "myapp"), 0750)
				if err != nil {
					t.Fatalf("failed to create dir: %v", err)
				}

				return &test_utils.MockFileSystem{FileSystem: fs}
			},
			setupTM: func() *test_utils.MockTemplateManager {
				return &test_utils.MockTemplateManager{}
			},
			setupCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			wantErr:     true,
			errContains: boilerplate.ErrDirectoryAlreadyExists.Error(),
		},
		{
			name: "directory already exists with force",
			opts: boilerplate.Options{
				Name: "myapp",
				Languages: map[string]string{
					"go": "go",
				},
				Force: true,
			},
			setupFS: func(t *testing.T, tmpDir string) *test_utils.MockFileSystem {
				t.Helper()

				fs := utils.NewFileSystem()

				err := fs.MkdirAll(filepath.Join(tmpDir, "myapp"), 0750)
				if err != nil {
					t.Fatalf("failed to create dir: %v", err)
				}

				return &test_utils.MockFileSystem{FileSystem: fs}
			},
			setupTM: func() *test_utils.MockTemplateManager {
				return &test_utils.MockTemplateManager{}
			},
			setupCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			wantErr: false,
		},
		{
			name: "stat error",
			opts: boilerplate.Options{
				Name: "myapp",
			},
			setupFS: func(t *testing.T, tmpDir string) *test_utils.MockFileSystem {
				t.Helper()

				return &test_utils.MockFileSystem{
					FileSystem: utils.NewFileSystem(),
					StatFunc: func(name string) (fs.FileInfo, error) {
						return nil, errMock
					},
				}
			},
			setupTM: func() *test_utils.MockTemplateManager {
				return &test_utils.MockTemplateManager{}
			},
			setupCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			wantErr:     true,
			errContains: "failed to check directory",
		},
		{
			name: "mkdir error",
			opts: boilerplate.Options{
				Name: "myapp",
			},
			setupFS: func(t *testing.T, tmpDir string) *test_utils.MockFileSystem {
				t.Helper()

				return &test_utils.MockFileSystem{
					FileSystem: utils.NewFileSystem(),
					MkdirAllFunc: func(path string, perm fs.FileMode) error {
						return errMock
					},
				}
			},
			setupTM: func() *test_utils.MockTemplateManager {
				return &test_utils.MockTemplateManager{}
			},
			setupCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			wantErr:     true,
			errContains: "failed to create directory",
		},
		{
			name: "template render error",
			opts: boilerplate.Options{
				Name: "myapp",
				Languages: map[string]string{
					"go": "go",
				},
			},
			setupFS: func(t *testing.T, tmpDir string) *test_utils.MockFileSystem {
				t.Helper()

				return &test_utils.MockFileSystem{FileSystem: utils.NewFileSystem()}
			},
			setupTM: func() *test_utils.MockTemplateManager {
				return &test_utils.MockTemplateManager{
					RenderTemplateFunc: func(lang string) (*bytes.Buffer, error) {
						return nil, errMock
					},
				}
			},
			setupCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			wantErr:     true,
			errContains: "failed to scaffold boilerplate with errors",
		},
		{
			name: "write file error",
			opts: boilerplate.Options{
				Name: "myapp",
				Languages: map[string]string{
					"go": "go",
				},
			},
			setupFS: func(t *testing.T, tmpDir string) *test_utils.MockFileSystem {
				t.Helper()

				return &test_utils.MockFileSystem{
					FileSystem: utils.NewFileSystem(),
					WriteFileFunc: func(name string, data []byte, perm fs.FileMode) error {
						return errMock
					},
				}
			},
			setupTM: func() *test_utils.MockTemplateManager {
				return &test_utils.MockTemplateManager{}
			},
			setupCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			wantErr:     true,
			errContains: "failed to scaffold boilerplate with errors",
		},
		{
			name: "context cancelled",
			opts: boilerplate.Options{
				Name: "myapp",
				Languages: map[string]string{
					"go": "go",
				},
			},
			setupFS: func(t *testing.T, tmpDir string) *test_utils.MockFileSystem {
				t.Helper()

				return &test_utils.MockFileSystem{FileSystem: utils.NewFileSystem()}
			},
			setupTM: func() *test_utils.MockTemplateManager {
				return &test_utils.MockTemplateManager{}
			},
			setupCtx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // Cancel immediately

				return ctx, cancel
			},
			wantErr:     true,
			errContains: context.Canceled.Error(),
		},
		{
			name: "rollback failure on file remove",
			opts: boilerplate.Options{
				Name: "myapp",
				Languages: map[string]string{
					"go": "go",
					"py": "py",
				},
			},
			setupFS: func(t *testing.T, tmpDir string) *test_utils.MockFileSystem {
				t.Helper()

				return &test_utils.MockFileSystem{
					FileSystem: utils.NewFileSystem(),
					WriteFileFunc: func(name string, data []byte, perm fs.FileMode) error {
						if strings.HasSuffix(name, ".py") {
							return errMock // Trigger error on second file
						}

						return utils.NewFileSystem().WriteFile(name, data, perm)
					},
					RemoveFunc: func(name string) error {
						return errMock // Fail to remove the first file
					},
				}
			},
			setupTM: func() *test_utils.MockTemplateManager {
				return &test_utils.MockTemplateManager{}
			},
			setupCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			wantErr:     true,
			errContains: "rollback had 1 failure(s)",
		},
		{
			name: "rollback failure on dir removeAll",
			opts: boilerplate.Options{
				Name: "myapp",
				Languages: map[string]string{
					"go": "go",
				},
			},
			setupFS: func(t *testing.T, tmpDir string) *test_utils.MockFileSystem {
				t.Helper()

				return &test_utils.MockFileSystem{
					FileSystem: utils.NewFileSystem(),
					WriteFileFunc: func(name string, data []byte, perm fs.FileMode) error {
						return errMock // Trigger error
					},
					RemoveAllFunc: func(path string) error {
						return errMock // Fail to remove dir
					},
				}
			},
			setupTM: func() *test_utils.MockTemplateManager {
				return &test_utils.MockTemplateManager{}
			},
			setupCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			wantErr:     true,
			errContains: "rollback had 1 failure(s)",
		},
		{
			name: "rollback failure on empty dir remove",
			opts: boilerplate.Options{
				Name: "myapp",
				Languages: map[string]string{
					"go": "go",
				},
				Force: true, // Directory already exists, so createdRootDir is false
			},
			setupFS: func(t *testing.T, tmpDir string) *test_utils.MockFileSystem {
				t.Helper()

				fsys := utils.NewFileSystem()

				err := fsys.MkdirAll(filepath.Join(tmpDir, "myapp"), 0750)
				if err != nil {
					t.Fatalf("failed to create dir: %v", err)
				}

				return &test_utils.MockFileSystem{
					FileSystem: fsys,
					WriteFileFunc: func(name string, data []byte, perm fs.FileMode) error {
						return errMock // Trigger error
					},
					RemoveFunc: func(name string) error {
						return errMock // Fail to remove empty dir
					},
				}
			},
			setupTM: func() *test_utils.MockTemplateManager {
				return &test_utils.MockTemplateManager{}
			},
			setupCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			wantErr:     true,
			errContains: "rollback had 1 failure(s)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()
			tt.opts.OutputDir = tmpDir

			fs := tt.setupFS(t, tmpDir)
			tm := tt.setupTM()

			ctx, cancel := tt.setupCtx()
			defer cancel()

			log := logger.New("error", false, true)
			b := boilerplate.New(tm, fs, log)

			err := b.Scaffold(ctx, tt.opts)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
