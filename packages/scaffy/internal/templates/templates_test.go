package templates_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/25prabhu10/scaffy/internal/templates"
)

func TestRenderTemplates(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	// Create a valid custom template
	validCustomTmpl := filepath.Join(tmpDir, "custom.tmpl")

	err := os.WriteFile(validCustomTmpl, []byte("Custom Template: {{.Author}}"), 0600)
	if err != nil {
		t.Fatalf("failed to create valid custom template: %v", err)
	}

	// Create an invalid custom template (fails at execution time)
	invalidCustomTmpl := filepath.Join(tmpDir, "invalid.tmpl")

	err = os.WriteFile(invalidCustomTmpl, []byte(`{{template "nonexistent"}}`), 0600)
	if err != nil {
		t.Fatalf("failed to create invalid custom template: %v", err)
	}

	// Create a template that fails at parse time (to test template.Must panic on parse)
	parseFailTmpl := filepath.Join(tmpDir, "parsefail.tmpl")

	err = os.WriteFile(parseFailTmpl, []byte(`{{if .Author}} missing end`), 0600)
	if err != nil {
		t.Fatalf("failed to create parse fail template: %v", err)
	}

	testData := templates.Data{
		Date:   "2026-02-22",
		Author: "Test Author",
		URL:    "https://example.com",
	}

	tests := []struct {
		name         string
		templatesDir string
		lang         string
		wantContains string
		wantErr      bool
		wantPanic    bool
	}{
		{
			name:         "valid embedded template",
			templatesDir: "",
			lang:         "go",
			wantContains: "Author: Test Author",
			wantErr:      false,
			wantPanic:    false,
		},
		{
			name:         "valid embedded template with whitespace dir",
			templatesDir: "   ",
			lang:         "go",
			wantContains: "Author: Test Author",
			wantErr:      false,
			wantPanic:    false,
		},
		{
			name:         "valid custom template",
			templatesDir: tmpDir,
			lang:         "custom",
			wantContains: "Custom Template: Test Author",
			wantErr:      false,
			wantPanic:    false,
		},
		{
			name:         "missing embedded template",
			templatesDir: "",
			lang:         "unknown",
			wantContains: "",
			wantErr:      false,
			wantPanic:    true,
		},
		{
			name:         "missing custom template",
			templatesDir: tmpDir,
			lang:         "unknown",
			wantContains: "",
			wantErr:      false,
			wantPanic:    true,
		},
		{
			name:         "execution error in custom template",
			templatesDir: tmpDir,
			lang:         "invalid",
			wantContains: "",
			wantErr:      true,
			wantPanic:    false,
		},
		{
			name:         "parse error in custom template",
			templatesDir: tmpDir,
			lang:         "parsefail",
			wantContains: "",
			wantErr:      false,
			wantPanic:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tm := templates.New(tt.templatesDir, testData)

			// Handle expected panics from template.Must
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("RenderTemplate() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()

			got, err := tm.RenderTemplate(tt.lang)

			if (err != nil) != tt.wantErr {
				t.Errorf("RenderTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !tt.wantPanic {
				if !strings.Contains(got.String(), tt.wantContains) {
					t.Errorf("RenderTemplate() got = %v, want to contain %v", got, tt.wantContains)
				}
			}
		})
	}
}
