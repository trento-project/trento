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
	"github.com/trento-project/trento/internal"
	"github.com/trento-project/trento/internal/db"
	"github.com/trento-project/trento/web"
	"github.com/trento-project/trento/web/datapipeline"
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
		PersistentPreRunE: func(webCmd *cobra.Command, _ []string) error {
			webCmd.Flags().VisitAll(func(f *pflag.Flag) {
				viper.BindPFlag(f.Name, f)
			})
			webCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
				viper.BindPFlag(f.Name, f)
			})
			return internal.InitConfig("web")
		},
	}

	webCmd.PersistentFlags().StringVar(&dbHost, "db-host", "localhost", "The database host")
	webCmd.PersistentFlags().StringVar(&dbPort, "db-port", "5432", "The database port to connect to")
	webCmd.PersistentFlags().StringVar(&dbUser, "db-user", "postgres", "The database user")
	webCmd.PersistentFlags().StringVar(&dbPassword, "db-password", "postgres", "The database password")
	webCmd.PersistentFlags().StringVar(&dbName, "db-name", "trento", "The database name that the application will use")

	addServeCmd(webCmd)
	addPruneCmd(webCmd)

	return webCmd
}

func addServeCmd(webCmd *cobra.Command) {
	var host string
	var port int

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

	serveCmd.Flags().IntVar(&collectorPort, "collector-port", 8081, "The port for the data collector service to listen on")
	serveCmd.Flags().BoolVar(&enablemTLS, "enable-mtls", false, "Enable mTLS authentication between server and agents")
	serveCmd.Flags().StringVar(&cert, "cert", "", "mTLS server certificate")
	serveCmd.Flags().StringVar(&key, "key", "", "mTLS server key")
	serveCmd.Flags().StringVar(&ca, "ca", "", "mTLS Certificate Authority")

	webCmd.AddCommand(serveCmd)
}

func addPruneCmd(webCmd *cobra.Command) {
	var olderThan uint

	pruneCmd := &cobra.Command{
		Use:   "prune-events",
		Short: "Prune events older than",
		Run:   prune,
	}

	pruneCmd.Flags().UintVar(&olderThan, "older-than", 10, "Prune data discovery events older than <value> days.")

	webCmd.AddCommand(pruneCmd)
}

func serve(*cobra.Command, []string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		log.Infof("Bye!")
		cancel()
	}()

	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Failed to configure the web application instance: ", err)
	}

	app, err := web.NewApp(config)
	if err != nil {
		log.Fatal("Failed to create the web application instance: ", err)
	}

	if err := app.Start(ctx); err != nil {
		log.Fatal("Error while running the web application server: ", err)
	}
}

func prune(_ *cobra.Command, _ []string) {
	olderThan := viper.GetUint("older-than")
	olderThanDuration := time.Duration(olderThan) * 24 * time.Hour

	dbConfig := LoadDBConfig()
	db, err := db.InitDB(dbConfig)
	if err != nil {
		log.Fatal("Error while initializing the database: ", err)
	}

	log.Infof("Pruning events older than %d days.", olderThan)
	if err := datapipeline.PruneEvents(olderThanDuration, db); err != nil {
		log.Fatalf("Error while pruning older events: %s", err)
	}
	log.Infof("Events older than %d days pruned.", olderThan)
}
