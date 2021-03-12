package webapp

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

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
	r.Get("/", renderTemplate)

	r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi home"))
	})
	listenAddress := fmt.Sprintf("%s:%d", host, port)
	err := http.ListenAndServe(listenAddress, r)
	log.Println("Error listening and serving:", err)
}

// Index data is used for the home template
type Index struct {
	Title string
}

func renderTemplate(w http.ResponseWriter, r *http.Request) {
	data := Index{
		Title: "Sapconsole",
	}
	parsedTemplate, err := template.ParseFiles("webapp/templates/home.html.tmpl")
	if err != nil {
		log.Println("Error parsing template :", err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	err = parsedTemplate.Execute(w, data)
	if err != nil {
		log.Println("Error executing template :", err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
}
