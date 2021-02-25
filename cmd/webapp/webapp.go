package webapp

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/SUSE/console-for-sap/webapp"
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
	engine := webapp.NewEngine()

	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", host, port),
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
