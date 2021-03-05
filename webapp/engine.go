package webapp

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/assets
var assetsFS embed.FS

//go:embed templates
var templatesFS embed.FS

var layoutData = gin.H{
	"title": "SUSE Console for SAP Applications",
	"copyright": "© 2019-2020 SUSE, all rights reserved.",
}

func NewEngine() *gin.Engine {

	engine := gin.Default()
	engine.HTMLRender = NewLayoutRender(templatesFS, layoutData, "templates/*.tmpl")

	engine.StaticFS("/static", http.FS(assetsFS))
	engine.GET("/", homeHandler)

	return engine
}
