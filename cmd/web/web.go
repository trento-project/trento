package web

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/SUSE/console-for-sap-applications/web"
)

var host string
var port int

func NewWebCmd() *cobra.Command {
	webCmd := &cobra.Command{
		Use:   "web",
		Short: "Command tree related to the web application component",
	}

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Starts the web application",
		Run:   serve,
	}

	serveCmd.Flags().StringVar(&host, "host", "0.0.0.0", "The host to bind the HTTP service to")
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "The port for the HTTP service to listen at")

	webCmd.AddCommand(serveCmd)

	return webCmd
}

func serve(cmd *cobra.Command, args []string) {
	app := web.NewApp(host, port)

	err := app.Start()
	if err != nil {
		log.Fatal("Failed to start the web application service: ", err)
	}
}
