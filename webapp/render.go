package webapp

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
)

const (
	root = "templates/base.html.tmpl"
)

// TemplateRender wraps user templates into a root one which has it's own data and a bunch of inner blocks
type TemplateRender struct {
	templates map[string]*template.Template
}

// The default constructor expects an FS, some data, and user templates;
// user templates are the ones that can be referenced by the Gin context.
func NewTemplateRender(templatesFS fs.FS, templates ...string) *TemplateRender {
	r := &TemplateRender{
		templates: map[string]*template.Template{},
	}
	for _, pattern := range templates {
		files, err := fs.Glob(templatesFS, pattern)
		if err != nil {
			// we exceptionally hard panic in case of glob errors, these should never happen.
			panic(err)
		}
		for _, file := range files {
			if file == root {
				continue
			}
			r.addFileFromFS(templatesFS, file)
		}
	}
	return r
}

// addFileFromFS parses the root template with the user
func (r *TemplateRender) addFileFromFS(templatesFS fs.FS, file string) {
	var tmpl *template.Template
	// use the base template first
	name := filepath.Base(file)
	tmpl = template.New(filepath.Base(root))

	// we "extend" the templates by adding custom functions
	tmpl = tmpl.Funcs(template.FuncMap{
		"escapedTemplate": func(name string, data interface{}) string {
			var out bytes.Buffer
			_ = tmpl.ExecuteTemplate(&out, name, data)
			return out.String()
		},
	})
	// parse all templates
	patterns := append([]string{root, file}, []string{"templates/blocks/*.html.tmpl"}...)
	tmpl = template.Must(tmpl.ParseFS(templatesFS, patterns...))

	// add template to template map, consumed by handlers
	if tmpl == nil {
		panic("template can not be nil")
	}
	if len(name) == 0 {
		panic("template name cannot be empty")
	}
	if _, ok := r.templates[name]; ok {
		panic(fmt.Sprintf("template %s already exists", name))
	}
	r.templates[name] = tmpl
}
