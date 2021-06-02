package web

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/gin-gonic/gin/render"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// LayoutRender wraps user templates into a root one which has it's own data and a bunch of inner blocks
type LayoutRender struct {
	data      LayoutData
	root      string   // the root template is separate because it has to be parsed first
	blocks    []string // blocks are used by the root template and can be redefined in user templates
	templates map[string]*template.Template
}

type LayoutData struct {
	Title     string
	Copyright string
	Submenu   Submenu
	Content   interface{}
}

type Submenu []SubmenuItem

type SubmenuItem struct {
	Label string
	URL   string
}

var defaultLayoutData = LayoutData{
	Title:     "Trento Console for SAP Applications",
	Copyright: "© 2020-2021 SUSE LLC",
}

// The default constructor expects an FS, some data, and user templates;
// user templates are the ones that can be referenced by the Gin context.
func NewLayoutRender(templatesFS fs.FS, templates ...string) *LayoutRender {
	r := &LayoutRender{
		data:      defaultLayoutData,
		root:      "templates/layout.html.tmpl",
		blocks:    []string{"templates/blocks/*.html.tmpl"},
		templates: map[string]*template.Template{},
	}

	r.addGlobFromFS(templatesFS, templates...)

	return r
}

// Instance returns a render.HTML instance with the associated named Template
func (r *LayoutRender) Instance(name string, data interface{}) render.Render {
	r.data.Content = data
	tmpl, ok := r.templates[name]
	if !ok {
		panic("template not found")
	}
	return render.HTML{
		Template: tmpl,
		Data:     r.data,
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
		"sum": func(a int, b int) int {
			return a + b
		},
		"markdown": markdownToHTML,
	})
	patterns := append([]string{r.root, file}, r.blocks...)
	tmpl = template.Must(tmpl.ParseFS(templatesFS, patterns...))

	r.addTemplate(name, tmpl)
}

func markdownToHTML(md string) template.HTML {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	markdownParser := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	htmlOptions := html.RendererOptions{Flags: htmlFlags}
	markdownRenderer := html.NewRenderer(htmlOptions)
	h := markdown.ToHTML([]byte(md), markdownParser, markdownRenderer)
	return template.HTML(h)
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
