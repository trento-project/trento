package cmd

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
)

// webappCmd represents the webapp command
var webappCmd = &cobra.Command{
	Use:   "webapp",
	Short: "Command tree related to the web application component",
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the web application",
	Run:   serve,
}

var host string
var port int

func init() {
	rootCmd.AddCommand(webappCmd)
	webappCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVar(&host, "host", "0.0.0.0", "The host to bind the HTTP service to")
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "The port for the HTTP service to listen at")
}

func serve(cmd *cobra.Command, args []string) {
	r := chi.NewRouter()

	r.Get("/", renderTemplate)

	r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi home"))
	})

	listenAddress := fmt.Sprintf("%s:%d", host, port)

	http.ListenAndServe(listenAddress, r)

}

// Index data is used for the home template
type Index struct {
	Title string
}

func renderTemplate(w http.ResponseWriter, r *http.Request) {
	data := Index{
		Title: "Sapconsole",
	}
	parsedTemplate, err := template.ParseFiles("webapp/templates/home.html")
	if err != nil {
		log.Println("Error parsing template :", err)
		return
	}
	err = parsedTemplate.Execute(w, data)
	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}
