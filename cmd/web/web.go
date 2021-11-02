package web

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/trento-project/trento/web"
)

func NewWebCmd() *cobra.Command {
	var dbHost string
	var dbPort string
	var dbUser string
	var dbPassword string
	var dbName string

	webCmd := &cobra.Command{
		Use:   "web",
		Short: "Command tree related to the web application component",
	}

	webCmd.PersistentFlags().StringVar(&dbHost, "db-host", "localhost", "The database host")
	webCmd.PersistentFlags().StringVar(&dbPort, "db-port", "5432", "The database port to connect to")
	webCmd.PersistentFlags().StringVar(&dbUser, "db-user", "postgres", "The database user")
	webCmd.PersistentFlags().StringVar(&dbPassword, "db-password", "postgres", "The database password")
	webCmd.PersistentFlags().StringVar(&dbName, "db-name", "trento", "The database name that the application will use")

	webCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag(f.Name, f)
	})

	addServeCmd(webCmd)
	addPruneCmd(webCmd)

	return webCmd
}

func addServeCmd(webCmd *cobra.Command) {
	var host string
	var port int
	var araAddr string

	var collectorPort int
	var enablemTLS bool
	var cert string
	var key string
	var ca string

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Starts the web application",
		Run:   serve,
	}

	serveCmd.Flags().StringVar(&host, "host", "0.0.0.0", "The host to bind the HTTP service to")
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "The port for the HTTP service to listen on")
	serveCmd.Flags().StringVar(&araAddr, "ara-addr", "127.0.0.1:8000", "Address where ARA is running (ex: localhost:80)")

	serveCmd.Flags().IntVar(&collectorPort, "collector-port", 8081, "The port for the data collector service to listen on")
	serveCmd.Flags().BoolVar(&enablemTLS, "enable-mtls", false, "Enable mTLS authentication between server and agents")
	serveCmd.Flags().StringVar(&cert, "cert", "", "mTLS server certificate")
	serveCmd.Flags().StringVar(&key, "key", "", "mTLS server key")
	serveCmd.Flags().StringVar(&ca, "ca", "", "mTLS Certificate Authority")

	// Bind the flags to viper and make them available to the application
	serveCmd.Flags().VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag(f.Name, f)
	})

	webCmd.AddCommand(serveCmd)
}

func addPruneCmd(webCmd *cobra.Command) {
	var olderThan uint

	pruneCmd := &cobra.Command{
		Use:   "prune-events",
		Short: "Prune events older than",
		Run:   prune,
	}

	pruneCmd.Flags().UintVar(&olderThan, "older-than", 10, "Prune data discoveryu events older than <value> days.")

	webCmd.AddCommand(pruneCmd)
}

func serve(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		log.Infof("Bye!")
		cancel()
	}()

	app := initApp()
	if err := app.Start(ctx); err != nil {
		log.Fatal("Error while running the web application server: ", err)
	}
}

func prune(cmd *cobra.Command, args []string) {
	olderThan := time.Duration(viper.GetInt("older-than")) * 24 * time.Hour

	app := initApp()
	log.Infof("Pruning events older than %d days.", olderThan)
	if err := app.PruneEvents(olderThan); err != nil {
		log.Fatalf("Error while pruning older events: %s", err)
	}
	log.New().Infof("Events older than %d days pruned.", olderThan)
}

func initApp() *web.App {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Failed to configure the web application instance: ", err)
	}

	app, err := web.NewApp(config)
	if err != nil {
		log.Fatal("Failed to create the web application instance: ", err)
	}

	return app
}
