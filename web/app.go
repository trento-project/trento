package web

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/web/services"
	"github.com/trento-project/trento/web/services/ara"
)

//go:embed frontend/assets
var assetsFS embed.FS

//go:embed templates
var templatesFS embed.FS

const araAddrDefault = "127.0.0.1:8000"

type App struct {
	host string
	port int
	Dependencies
}

type Dependencies struct {
	consul        consul.Client
	engine        *gin.Engine
	store         cookie.Store
	checksService services.ChecksService
}

func DefaultDependencies() Dependencies {
	consulClient, _ := consul.DefaultClient()
	engine := gin.Default()
	store := cookie.NewStore([]byte("secret"))

	araService := ara.NewAraService(araAddrDefault)
	checksService := services.NewChecksService(araService)

	return Dependencies{consulClient, engine, store, checksService}
}

func (d *Dependencies) SetAraAddr(araAddr string) {
	araService := ara.NewAraService(araAddr)
	d.checksService = services.NewChecksService(araService)
}

// shortcut to use default dependencies
func NewApp(host string, port int) (*App, error) {
	return NewAppWithDeps(host, port, DefaultDependencies())
}

func NewAppWithDeps(host string, port int, deps Dependencies) (*App, error) {
	app := &App{
		Dependencies: deps,
		host:         host,
		port:         port,
	}

	InitAlerts()
	engine := deps.engine
	engine.HTMLRender = NewLayoutRender(templatesFS, "templates/*.tmpl")
	engine.Use(ErrorHandler)
	engine.Use(sessions.Sessions("session", deps.store))
	engine.StaticFS("/static", http.FS(assetsFS))
	engine.GET("/", HomeHandler)
	engine.GET("/hosts", NewHostListHandler(deps.consul))
	engine.GET("/hosts/:name", NewHostHandler(deps.consul))
	engine.GET("/catalog", NewChecksCatalogHandler(deps.checksService))
	engine.GET("/clusters", NewClusterListHandler(deps.consul, deps.checksService))
	engine.GET("/clusters/:id", NewClusterHandler(deps.consul, deps.checksService))
	engine.POST("/clusters/:id/checks", NewSaveChecksHandler(deps.consul))
	engine.GET("/environments", NewEnvironmentListHandler(deps.consul))
	engine.GET("/environments/:env", NewEnvironmentHandler(deps.consul))
	engine.GET("/landscapes", NewLandscapeListHandler(deps.consul))
	engine.GET("/landscapes/:land", NewLandscapeHandler(deps.consul))
	engine.GET("/sapsystems", NewSAPSystemListHandler(deps.consul))
	engine.GET("/sapsystems/:sys", NewSAPSystemHandler(deps.consul))

	apiGroup := engine.Group("/api")
	{
		apiGroup.GET("/ping", ApiPingHandler)
	}

	return app, nil
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
