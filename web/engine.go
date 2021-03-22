package web

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/trento-project/trento/web/environments"
)

//go:embed frontend/assets
var assetsFS embed.FS

//go:embed templates
var templatesFS embed.FS

var layoutData = gin.H{
	"title":     "Trento Web Console for SAP Applications administrators",
	"copyright": "© 2019-2021 SUSE LLC",
}

func NewEngine() *gin.Engine {

	engine := gin.Default()
	engine.HTMLRender = NewLayoutRender(templatesFS, layoutData, "templates/*.tmpl")

	engine.StaticFS("/static", http.FS(assetsFS))
	engine.GET("/", homeHandler)
	engine.GET("/environments", envronments.ListHandler)

	return engine
}
