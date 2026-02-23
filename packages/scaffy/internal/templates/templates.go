package templates

import (
	"bytes"
	"embed"
	"html/template"
	"path/filepath"

	"github.com/25prabhu10/scaffy/internal/utils"
)

//go:embed *.tmpl
var embeddedTemplates embed.FS

// Data holds the dynamic data to be injected into the templates.
type Data struct {
	Date   string
	Author string
	URL    string
}

type TemplateManager interface {
	RenderTemplate(lang string) (*bytes.Buffer, error)
}

// TemplateManager manages the loading and rendering of templates.
type templateManager struct {
	templatesDir string
	templateData Data
}

// New creates a new instance of TemplateManager with the specified templates directory and template data.
func New(templatesDir string, templateData Data) TemplateManager {
	return &templateManager{
		templatesDir: templatesDir,
		templateData: templateData,
	}
}

// RenderTemplate renders the template for the specified language and returns the resulting string.
func (tm *templateManager) RenderTemplate(lang string) (*bytes.Buffer, error) {
	templateName := lang + ".tmpl"

	tmpl := tm.loadTemplate(templateName)

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, templateName, tm.templateData); err != nil {
		return nil, err
	}

	return &buf, nil
}

func (tm *templateManager) loadTemplate(templateName string) *template.Template {
	if !utils.IsStringEmpty(tm.templatesDir) {
		return template.Must(template.New(templateName).ParseFiles(filepath.Join(tm.templatesDir, templateName)))
	}

	return template.Must(template.New(templateName).ParseFS(embeddedTemplates, templateName))
}
