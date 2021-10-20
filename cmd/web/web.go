package web

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/trento-project/trento/web"
)

var host string
var port int
var araAddr string

var dbHost string
var dbPort string
var dbUser string
var dbPassword string
var dbName string
var enablemTLS bool
var cert string
var key string
var ca string

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

	serveCmd.Flags().StringVar(&dbHost, "db-host", "localhost", "The database host")
	serveCmd.Flags().StringVar(&dbPort, "db-port", "5432", "The database port to connect to")
	serveCmd.Flags().StringVar(&dbUser, "db-user", "postgres", "The database user")
	serveCmd.Flags().StringVar(&dbPassword, "db-password", "postgres", "The database password")
	serveCmd.Flags().StringVar(&dbName, "db-name", "trento", "The database name that the application will use")

	serveCmd.Flags().BoolVar(&enablemTLS, "enable-mtls", false, "Enable mTLS authentication between server and agents")
	serveCmd.Flags().StringVar(&cert, "cert", "", "mTLS server certificate")
	serveCmd.Flags().StringVar(&key, "key", "", "mTLS server key")
	serveCmd.Flags().StringVar(&ca, "ca", "", "mTLS Certificate Authority")

	// Bind the flags to viper and make them available to the application
	serveCmd.Flags().VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag(f.Name, f)
	})

	webCmd.AddCommand(serveCmd)

	return webCmd
}

func serve(cmd *cobra.Command, args []string) {
	var err error

	deps := web.DefaultDependencies()

	app, err := web.NewAppWithDeps(host, port, deps)
	if err != nil {
		log.Fatal("Failed to create the web application instance: ", err)
	}

	err = app.Start()
	if err != nil {
		log.Fatal("Failed to start the web application service: ", err)
	}
}
