package web

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SUSE/console-for-sap-applications/web/api"
	"github.com/SUSE/console-for-sap-applications/web/envronments"
)

//go:embed frontend/assets
var assetsFS embed.FS

//go:embed templates
var templatesFS embed.FS

var layoutData = gin.H{
	"title":     "SUSE Console for SAP Applications",
	"copyright": "© 2019-2020 SUSE, all rights reserved.",
}

func NewEngine() *gin.Engine {

	engine := gin.Default()
	engine.HTMLRender = NewLayoutRender(templatesFS, layoutData, "templates/*.tmpl")

	engine.StaticFS("/static", http.FS(assetsFS))
	engine.GET("/", homeHandler)
	engine.GET("/environments", envronments.ListHandler)

	apiGroup := engine.Group("/api")
	{
		apiGroup.GET("/ping", api.PingHandler)
	}

	return engine
}
