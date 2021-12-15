package web

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/version"

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
	Version   string
	Flavor    string
	Submenu   Submenu
	Content   interface{}
}

type Submenu []SubmenuItem

type SubmenuItem struct {
	Label string
	URL   string
}

var defaultLayoutData = LayoutData{
	Title:     "Trento Console",
	Copyright: "© 2020-2021 SUSE LLC",
	Version:   version.Version,
	Flavor:    version.Flavor,
}

type LayoutHTML struct {
	Templates    map[string]*template.Template
	TemplateName string
	Data         interface{}
}

func (r LayoutHTML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	tmpl, ok := r.Templates[r.TemplateName]
	if !ok {
		err := fmt.Errorf("template %s not found", r.TemplateName)
		r.RenderErrorPage(InternalServerError(err.Error()), w)
		return err
	}

	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, r.Data)
	if err != nil {
		r.RenderErrorPage(InternalServerError(err.Error()), w)
		return err
	}

	_, err = w.Write(buf.Bytes())
	return err
}

func (r LayoutHTML) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{"text/html; charset=utf-8"}
	}
}

func (r LayoutHTML) RenderErrorPage(e *HttpError, w http.ResponseWriter) {
	tmpl, ok := r.Templates[e.template]
	if !ok {
		panic("error page template not found")
	}
	w.WriteHeader(e.code)
	err := tmpl.Execute(w, r.Data)

	if err != nil {
		log.Fatal("Error while rendering error page template", err)
	}
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

	return LayoutHTML{
		Templates:    r.templates,
		TemplateName: name,
		Data:         r.data,
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
		"split":    strings.Split,
		"script":   script,
	})
	patterns := append([]string{r.root, file}, r.blocks...)
	tmpl = template.Must(tmpl.ParseFS(templatesFS, patterns...))

	r.addTemplate(name, tmpl)
}

func script(filename string) template.HTML {
	scriptTag := fmt.Sprintf("<script src=\"/static/frontend/assets/js/%s\"></script>", filename)
	return template.HTML(scriptTag)
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
