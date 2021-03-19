package web

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
)

const (
	base = "templates/base.html.tmpl"
)

// the main motivation of this code is to parse all templates contained in templates directory,
// respecting the "base" template which is the basis, where other custom templates will inherit and override content
// we also extend the template and build the base template with blocks templated (contained in /templates/blocks directory)

type TemplateRender struct {
	templates map[string]*template.Template
}

// TemplateRender is a map of parsed html templates which is consumed by http handler, accessed via key
// from the embeded filesystems,
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
			if file == base {
				continue
			}
			r.addFileFromFS(templatesFS, file)
		}
	}
	return r
}

// addFileFromFS parses the base template with the user
func (r *TemplateRender) addFileFromFS(templatesFS fs.FS, file string) {
	var tmpl *template.Template
	// use the base template first
	name := filepath.Base(file)
	tmpl = template.New(filepath.Base(base))

	// we "extend" the templates by adding custom functions
	tmpl = tmpl.Funcs(template.FuncMap{
		"escapedTemplate": func(name string, data interface{}) string {
			var out bytes.Buffer
			_ = tmpl.ExecuteTemplate(&out, name, data)
			return out.String()
		},
	})
	// parse all templates
	patterns := append([]string{base, file}, []string{"templates/blocks/*.html.tmpl"}...)
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
