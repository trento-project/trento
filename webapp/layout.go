package webapp

import (
	"embed"
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin/render"
)

//go:embed templates
var templatesFS embed.FS

type Render struct {
	Templates map[string]*template.Template
}

func LayoutRenderer() Render {
	r := Render{
		Templates: map[string]*template.Template{},
	}
	return r
}

// Add new template
func (r Render) Add(name string, tmpl *template.Template) {
	if tmpl == nil {
		panic("template can not be nil")
	}
	if len(name) == 0 {
		panic("template name cannot be empty")
	}
	if _, ok := r.Templates[name]; ok {
		panic(fmt.Sprintf("template %s already exists", name))
	}
	r.Templates[name] = tmpl
}

// AddFromFiles supply add template from files
func (r Render) AddFromFS(name string, files ...string) *template.Template {
	tmpl := template.Must(template.ParseFS(templatesFS, files...))
	r.Add(name, tmpl)
	return tmpl
}

// Instance supply render string
func (r Render) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: r.Templates[name],
		Data:     data,
	}
}
