package web

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudquery/sqlite"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"

	"github.com/trento-project/trento/web/projectors"

	"github.com/gin-gonic/gin"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/web/service"
)

//go:embed frontend/assets
var assetsFS embed.FS

//go:embed templates
var templatesFS embed.FS

type App struct {
	host string
	port int
	Dependencies
}

type Dependencies struct {
	consul       consul.Client
	engine       *gin.Engine
	db           *gorm.DB
	hostsService service.IHostsService
}

func DefaultDependencies() Dependencies {
	consulClient, _ := consul.DefaultClient()
	engine := gin.Default()

	db, err := gorm.Open(sqlite.Open("trento.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		panic("failed to connect database")
	}

	// NOTE: sqlite drive does not support multiple threads
	dbConfig, _ := db.DB()
	dbConfig.SetMaxOpenConns(1)

	err = db.AutoMigrate(models.Host{}, projectors.Subscription{})
	if err != nil {
		panic("failed to migrate the database")
	}

	hostsService := service.NewHostsService(db)

	return Dependencies{consul: consulClient, engine: engine, db: db, hostsService: hostsService}
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

	engine := deps.engine
	engine.HTMLRender = NewLayoutRender(templatesFS, "templates/*.tmpl")
	engine.Use(ErrorHandler)
	engine.StaticFS("/static", http.FS(assetsFS))
	engine.GET("/", HomeHandler)
	engine.GET("/hosts", NewHostListHandler(deps.hostsService))
	engine.GET("/hosts/:name", NewHostHandler(deps.consul))
	engine.GET("/hosts/:name/ha-checks", NewHAChecksHandler(deps.consul))
	engine.GET("/clusters", NewClusterListHandler(deps.consul))
	engine.GET("/clusters/:id", NewClusterHandler(deps.consul))
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
	hostsProjectorHandler := projectors.NewHostsHandler("hosts", 5*time.Second, a.Dependencies.consul)
	hostsProjector := projectors.NewProjector(hostsProjectorHandler, a.db)
	hostsProjector.Run()

	hostsHealthProjectorHandler := projectors.NewHostsHealthHandler("hosts_health", 5*time.Second, a.Dependencies.consul)
	hostsHealthProjector := projectors.NewProjector(hostsHealthProjectorHandler, a.db)
	hostsHealthProjector.Run()

	return s.ListenAndServe()
}

func (a *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.engine.ServeHTTP(w, req)
}
