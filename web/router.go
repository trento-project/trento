package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

//go:embed templates
var templatesFS embed.FS

//go:embed frontend/assets
var assetsFS embed.FS

func InitRouter() chi.Router {
	r := chi.NewRouter()
	// parse all templates and return a map, which is consumed by each handler
	templs := NewTemplateRender(templatesFS, "templates/*.tmpl")
	// filesystem for static file
	filesDir, err := fs.Sub(assetsFS, "frontend/assets")
	if err != nil {
		panic(err)
	}

	r.Get("/", IndexHandler(templs.templates))

	FileServer(r, "/static", http.FS(filesDir))
	return r
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
