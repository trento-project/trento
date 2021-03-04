package webapp

import (
	"bytes"
	"embed"
	"fmt"
	"html"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

//go:embed templates
var templatesFS embed.FS

// Get layout templates contants
func getLayoutTemplates() []string {
	// I cannot make the templates/layout/*.tmpl work for some reason
	return []string{
		"templates/layout/layout.html.tmpl",
		"templates/layout/footer.html.tmpl",
		"templates/layout/source_footer.html.tmpl",
		"templates/layout/sidebar.html.tmpl",
		"templates/layout/header.html.tmpl",
		"templates/layout/submenu.html.tmpl",
	}
}

type Render struct {
	CommonData      map[string]string
	LayoutTemplates []string
	Templates       map[string]*template.Template
}

func LayoutRenderer() Render {
	r := Render{
		CommonData:      map[string]string{},
		LayoutTemplates: []string{},
		Templates:       map[string]*template.Template{},
	}
	return r
}

// Create the source footer html content
// This method is not correclty working as the html is totally unescaped
// The tooltip should have the url references properly
func sourceFooter() string {
	var result bytes.Buffer
	tmpl := template.Must(
		template.New("source_footer").ParseFS(
			templatesFS, "templates/layout/source_footer.html.tmpl"))
	tmpl.Execute(&result, "")
	return html.UnescapeString(result.String())
}

// Initialize the layout
// 1. Set the layout templates
// 2. Set the source footer data
func (r *Render) InitLayout() {
	r.AddLayoutTemplates(getLayoutTemplates()...)
	r.CommonData["source_footer"] = sourceFooter()
}

// Add layout templates
func (r *Render) AddLayoutTemplates(files ...string) {
	r.LayoutTemplates = append(r.LayoutTemplates, files...)
}

// Add new commond data that is used to render the layout template
func (r *Render) AddLayoutData(key string, value string) {
	r.CommonData[key] = value
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
func (r Render) AddTemplateFromFS(name string, files ...string) *template.Template {
	combined_files := append(r.LayoutTemplates, files...)
	tmpl := template.Must(template.ParseFS(templatesFS, combined_files...))
	r.Add(name, tmpl)
	return tmpl
}

// Set the added common data to the data interface
func (r Render) setCommonData(data interface{}) {
	for key, value := range r.CommonData {
		data.(gin.H)[key] = value
	}
}

// Instance supply render string
func (r Render) Instance(name string, data interface{}) render.Render {
	r.setCommonData(data)
	return render.HTML{
		Template: r.Templates[name],
		Data:     data,
	}
}
