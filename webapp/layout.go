package webapp

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

//go:embed templates
var templatesFS embed.FS

type Render struct {
	Data      map[string]interface{}
	root	  string
	patterns  []string
	templates map[string]*template.Template
}

func NewLayout() *Render {
	r := &Render{
		Data: map[string]interface{}{},
		root: "templates/layout/layout.html.tmpl",
		patterns: []string{
			"templates/layout/*",
		},
		//patterns: []string{"templates/layout/*.tmpl"},
		templates: map[string]*template.Template{},
	}

	r.Data["title"] = "SUSE Console for SAP Applications"
	r.Data["copyright"] = "© 2019-2020 SUSE, all rights reserved."

	return r
}

// Instance returns a render.HTML instance with the associated named Template
func (r Render) Instance(name string, data interface{}) render.Render {
	r.addLayoutData(data)
	return render.HTML{
		Template: r.templates[name],
		Data:     data,
	}
}

// Add new template
func (r Render) Add(name string, tmpl *template.Template) {
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

func (r Render) AddFromEmbeddedFS(patterns ...string) {
	for _, pattern := range patterns {
		var tmpl *template.Template

		name := filepath.Base(pattern)
		tmpl = template.New(filepath.Base(r.root))
		tmpl = tmpl.Funcs(template.FuncMap{
			"escapedTemplate": func(name string, data interface{}) string {
				var out bytes.Buffer
				_ = tmpl.ExecuteTemplate(&out, name, data)
				return out.String()
			},
		})
		tmpl = template.Must(tmpl.ParseFS(templatesFS, append([]string{r.root, pattern}, r.patterns...)...))

		r.Add(name, tmpl)
	}
}

// Add layout data to the data interface
func (r Render) addLayoutData(data interface{}) {
	for key, value := range r.Data {
		data.(gin.H)[key] = value
	}
}
