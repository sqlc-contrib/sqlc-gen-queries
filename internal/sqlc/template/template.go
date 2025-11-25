// Package template provides embedded SQL templates for code generation.
package template

import (
	"embed"
	"strings"
	"text/template"

	"github.com/go-openapi/inflect"
)

//go:embed *.tmpl
var fs embed.FS

// Open loads and parses a template file from the embedded filesystem.
func Open(name string, opts ...map[string]any) (*template.Template, error) {
	file := template.New(name)
	file.Funcs(template.FuncMap{
		"singular": inflect.Singularize,
		"plural":   inflect.Pluralize,
		"camelize": inflect.Camelize,
		"join":     strings.Join,
	})
	for _, item := range opts {
		file.Funcs(item)
	}

	return file.ParseFS(fs, name)
}
