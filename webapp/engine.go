package webapp

import (
	"embed"
	"html/template"

	"github.com/gin-gonic/gin"
)

//go:embed templates
var FS embed.FS

func Engine() *gin.Engine {

	templates := template.Must(template.New("").ParseFS(FS, "templates/*.tmpl"))

	engine := gin.Default()
	engine.SetHTMLTemplate(templates)
	engine.GET("/", homeHandler)

	return engine
}
