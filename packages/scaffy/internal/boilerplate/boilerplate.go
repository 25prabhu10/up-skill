package boilerplate

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"

	"github.com/25prabhu10/scaffy/internal/templates"
	"github.com/25prabhu10/scaffy/internal/utils"
)

var (
	ErrDirectoryAlreadyExists = errors.New("target directory already exists")
)

type Boilerplate interface {
	Scaffold(ctx context.Context, opts Options) error
}

type Options struct {
	Name      string
	OutputDir string
	Languages map[string]string
	Force     bool
}

type boilerplate struct {
	templateManager templates.TemplateManager
	fs              utils.FileSystem
	logger          *slog.Logger
}

func New(templateManager templates.TemplateManager, fs utils.FileSystem, logger *slog.Logger) Boilerplate {
	return &boilerplate{
		templateManager: templateManager,
		fs:              fs,
		logger:          logger,
	}
}

func (b *boilerplate) Scaffold(ctx context.Context, opts Options) error {
	targetDir := filepath.Join(opts.OutputDir, opts.Name)

	createdRootDir, err := b.prepareDirectory(targetDir, opts.Force)
	if err != nil {
		return err
	}

	createdFiles, errs := b.generateFiles(ctx, targetDir, opts)

	if len(errs) > 0 {
		return b.handleScaffoldErrors(targetDir, createdFiles, createdRootDir, errs)
	}

	return nil
}

func (b *boilerplate) prepareDirectory(targetDir string, force bool) (bool, error) {
	exists, err := utils.IsDirectoryExists(targetDir, b.fs)
	if err != nil {
		return false, err
	} else if exists {
		if !force {
			return false, fmt.Errorf("%w: %s", ErrDirectoryAlreadyExists, targetDir)
		}
	}

	if err = utils.CreateDirectoryIfNotExists(targetDir, b.fs); err != nil {
		return false, err
	}

	return !exists, nil
}

func (b *boilerplate) generateFiles(ctx context.Context, targetDir string, opts Options) ([]string, []error) {
	var createdFiles []string

	var mu sync.Mutex

	var wg sync.WaitGroup

	errCh := make(chan error, len(opts.Languages))

	for lang, fileExt := range opts.Languages {
		wg.Add(1)

		go func(lang, fileExt string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
			}

			templateContent, err := b.templateManager.RenderTemplate(lang)
			if err != nil {
				errCh <- err
				return
			}

			filePath := filepath.Join(targetDir, fmt.Sprintf("%s.%s", opts.Name, fileExt))

			if err = b.fs.WriteFile(filePath, templateContent.Bytes(), 0600); err != nil {
				errCh <- err
				return
			}

			mu.Lock()
			defer mu.Unlock()

			createdFiles = append(createdFiles, filePath)
		}(lang, fileExt)
	}

	wg.Wait()
	close(errCh)

	var errs []error

	for err := range errCh {
		if err != nil {
			errs = append(errs, err)
		}
	}

	return createdFiles, errs
}

func (b *boilerplate) handleScaffoldErrors(targetDir string, createdFiles []string, createdRootDir bool, errs []error) error {
	rollbackFailures := b.rollback(targetDir, createdFiles, createdRootDir)

	if len(rollbackFailures) > 0 {
		b.logger.Error("scaffolding failed with errors and rollback had failures", "errors", len(errs), "rollback_failures", len(rollbackFailures))

		return fmt.Errorf("scaffolding failed with %d error(s) and rollback had %d failure(s): %w; rollback failures: %s",
			len(errs), len(rollbackFailures), errors.Join(errs...), strings.Join(rollbackFailures, "; "))
	}

	b.logger.Error("failed to scaffold boilerplate with errors", "errors", len(errs))

	return fmt.Errorf("failed to scaffold boilerplate with errors: %w, all created files have been cleaned up", errors.Join(errs...))
}

func (b *boilerplate) rollback(targetDir string, createdFiles []string, createdRootDir bool) []string {
	var rollbackFailures []string

	for _, file := range createdFiles {
		if err := b.fs.Remove(file); err != nil {
			rollbackFailures = append(rollbackFailures, fmt.Sprintf("%s: %s", file, err.Error()))
		}
	}

	if createdRootDir {
		if err := b.fs.RemoveAll(targetDir); err != nil {
			rollbackFailures = append(rollbackFailures, fmt.Sprintf("%s: %s", targetDir, err.Error()))
		}

		return rollbackFailures
	}

	// Try to remove problem directory if empty
	if entries, err := b.fs.ReadDir(targetDir); err == nil && len(entries) == 0 {
		if err := b.fs.Remove(targetDir); err != nil {
			rollbackFailures = append(rollbackFailures, fmt.Sprintf("%s: %s", targetDir, err.Error()))
		}
	}

	return rollbackFailures
}
