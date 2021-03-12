package webapp

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/SUSE/console-for-sap-applications/webapp"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
)

var host string
var port int

func NewWebappCmd() *cobra.Command {
	webappCmd := &cobra.Command{
		Use:   "webapp",
		Short: "Command tree related to the web application component",
	}

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Starts the web application",
		Run:   serve,
	}

	serveCmd.Flags().StringVar(&host, "host", "0.0.0.0", "The host to bind the HTTP service to")
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "The port for the HTTP service to listen at")

	webappCmd.AddCommand(serveCmd)

	return webappCmd
}

func serve(cmd *cobra.Command, args []string) {
	r := chi.NewRouter()

	r.Get("/", webapp.IndexHandler)
	r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi home"))
	})

	// Create a route along /files that will serve contents from
	// the ./data/ folder.
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "webapp/frontend/assets/"))
	FileServer(r, "/static", filesDir)

	listenAddress := fmt.Sprintf("%s:%d", host, port)
	err := http.ListenAndServe(listenAddress, r)
	if err != nil {
		log.Println("Error while serving HTTP:", err)
	}
	log.Printf("serving on port %s", listenAddress)

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
