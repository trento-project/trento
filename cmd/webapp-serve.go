package cmd

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/SUSE/console-for-sap/webapp"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the web application",
	Run:   serve,
}

var host string
var port int

func init() {
	webappCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVar(&host, "host", "0.0.0.0", "The host to bind the HTTP service to")
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "The port for the HTTP service to listen at")
}

func serve(cmd *cobra.Command, args []string) {
	engine := gin.Default()
	engine.LoadHTMLGlob("webapp/templates/*.tpl")
	engine.GET("/", webapp.Home)

	listenAddress := fmt.Sprintf("%s:%d", host, port)
	log.Fatal(engine.Run(listenAddress))
}
