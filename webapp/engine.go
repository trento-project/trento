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

	renderer := LayoutRenderer()
	renderer.InitLayout()
	renderer.AddLayoutData("title", "SUSE Console for SAP Applications")
	renderer.AddLayoutData("footer", "© 2019-2020 SUSE, all rights reserved.")
	renderer.AddTemplateFromFS(
		"home", "templates/home.html.tmpl")
	engine.HTMLRender = renderer

	engine.StaticFS("/static", http.FS(assetsFS))
	engine.GET("/", homeHandler)

	return engine
}
