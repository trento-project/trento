package web

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/trento-project/trento/web/api"
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

type App struct {
	engine *gin.Engine
	host   string
	port   int
}

func NewApp(host string, port int) *App {
	engine := gin.Default()
	engine.HTMLRender = NewLayoutRender(templatesFS, layoutData, "templates/*.tmpl")

	engine.StaticFS("/static", http.FS(assetsFS))
	engine.GET("/", homeHandler)
	engine.GET("/environments", envronments.ListHandler)

	apiGroup := engine.Group("/api")
	{
		apiGroup.GET("/ping", api.PingHandler)
	}

	return &App{engine, host, port}
}

func (a *App) Start() error {
	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", a.host, a.port),
		Handler:        a,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s.ListenAndServe()
}

func (a *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.engine.ServeHTTP(w, req)
}
