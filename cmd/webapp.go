package cmd

import (
	"fmt"
	"html/template"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/SUSE/console-for-sap/webapp"
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

func makeEngine() *gin.Engine {
	templates := template.Must(template.New("").ParseFS(webapp.FS, "templates/*.tmpl"))

	engine := gin.Default()
	engine.SetHTMLTemplate(templates)
	engine.GET("/", webapp.Home)

	return engine
}

func serve(cmd *cobra.Command, args []string) {
	engine := makeEngine()

	listenAddress := fmt.Sprintf("%s:%d", host, port)
	log.Fatal(engine.Run(listenAddress))
}
