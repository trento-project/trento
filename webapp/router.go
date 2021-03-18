package webapp

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

//go:embed templates
var templateFS embed.FS

//go:embed frontend/assets
var assetsFS embed.FS

// InitRouter initialize the http router
func InitRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", IndexHandler(templateFS))
	filesDir, err := fs.Sub(assetsFS, "frontend/assets")
	if err != nil {
		panic(err)
	}
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
