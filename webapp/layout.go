package webapp

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

// LayoutRender wraps user templates into a root one which has it's own data and a bunch of inner blocks
type LayoutRender struct {
	Data      gin.H
	root      string   // the root template is separate because it has to be parsed first
	blocks    []string // blocks are used by the root template and can be redefined in user templates
	templates map[string]*template.Template
}

// The default constructor expects an FS, some data, and user templates;
// user templates are the ones that can be referenced by the Gin context.
func NewLayoutRender(templatesFS fs.FS, data gin.H, templates ...string) *LayoutRender {
	r := &LayoutRender{
		Data:      data,
		root:      "templates/layout.html.tmpl",
		blocks:    []string{"templates/includes/*.html.tmpl"},
		templates: map[string]*template.Template{},
	}

	r.addGlobFromFS(templatesFS, templates...)

	return r
}

// Instance returns a render.HTML instance with the associated named Template
func (r *LayoutRender) Instance(name string, data interface{}) render.Render {
	r.addLayoutData(data)
	return render.HTML{
		Template: r.templates[name],
		Data:     data,
	}
}

// addGlobFromFS expands globs so that each user template is added under a name
func (r *LayoutRender) addGlobFromFS(templatesFS fs.FS, patterns ...string) {
	for _, pattern := range patterns {
		files, err := fs.Glob(templatesFS, pattern)
		if err != nil {
			// we exceptionally hard panic in case of glob errors, these should never happen.
			panic(err)
		}
		for _, file := range files {
			if file == r.root {
				continue
			}
			r.addFileFromFS(templatesFS, file)
		}
	}
}

// addFileFromFS parses the root template with the user
func (r *LayoutRender) addFileFromFS(templatesFS fs.FS, file string) {
	var tmpl *template.Template

	name := filepath.Base(file)
	tmpl = template.New(filepath.Base(r.root))
	tmpl = tmpl.Funcs(template.FuncMap{
		"escapedTemplate": func(name string, data interface{}) string {
			var out bytes.Buffer
			_ = tmpl.ExecuteTemplate(&out, name, data)
			return out.String()
		},
	})
	patterns := append([]string{r.root, file}, r.blocks...)
	tmpl = template.Must(tmpl.ParseFS(templatesFS, patterns...))

	r.addTemplate(name, tmpl)
}

// addTemplate adds a new user template to the render
func (r *LayoutRender) addTemplate(name string, tmpl *template.Template) {
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

// addTemplate adds the root template data to the data passed to the user template
func (r *LayoutRender) addLayoutData(data interface{}) {
	for key, value := range r.Data {
		data.(gin.H)[key] = value
	}
}
