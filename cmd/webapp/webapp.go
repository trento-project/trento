package webapp

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SUSE/console-for-sap-applications/webapp"
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

	r := webapp.InitRouter()

	listenAddress := fmt.Sprintf("%s:%d", host, port)
	err := http.ListenAndServe(listenAddress, r)
	if err != nil {
		log.Println("Error while serving HTTP:", err)
	}
	log.Printf("serving on port %s", listenAddress)

}
