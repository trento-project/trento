package web

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/trento-project/trento/web"
)

var host string
var port int
var araAddr string

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
	serveCmd.Flags().StringVar(&araAddr, "ara-addr", "127.0.0.1:8000", "Address where ARA is running (ex: localhost:80)")

	webCmd.AddCommand(serveCmd)

	return webCmd
}

func serve(cmd *cobra.Command, args []string) {
	var err error

	deps := web.DefaultDependencies()
	deps.SetAraAddr(araAddr)

	app, err := web.NewAppWithDeps(host, port, deps)
	if err != nil {
		log.Fatal("Failed to create the web application instance: ", err)
	}

	err = app.Start()
	if err != nil {
		log.Fatal("Failed to start the web application service: ", err)
	}
}
