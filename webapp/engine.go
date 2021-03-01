package webapp

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed templates
var templatesFS embed.FS

//go:embed frontend/assets
var assetsFS embed.FS

func NewEngine() *gin.Engine {

	templates := template.Must(template.New("").ParseFS(templatesFS, "templates/*.tmpl"))

	engine := gin.Default()
	engine.SetHTMLTemplate(templates)

	engine.StaticFS("/static", http.FS(assetsFS))
	engine.GET("/", homeHandler)

	return engine
}
