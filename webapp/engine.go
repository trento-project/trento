package webapp

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/assets
var assetsFS embed.FS

func NewEngine() *gin.Engine {

	engine := gin.Default()

	renderer := NewLayout()
	renderer.AddFromEmbeddedFS("templates/home.html.tmpl")

	engine.HTMLRender = renderer

	engine.StaticFS("/static", http.FS(assetsFS))
	engine.GET("/", homeHandler)

	return engine
}
