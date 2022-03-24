package web

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/trento-project/trento/cmd/db"
	"github.com/trento-project/trento/internal"
	"github.com/trento-project/trento/web"
)

func NewWebCmd() *cobra.Command {
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

	db.AddDBFlags(webCmd)
	addServeCmd(webCmd)

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

	var grafanaPublicURL string
	var grafanaApiURL string
	var grafanaUser string
	var grafanaPassword string

	var prometheusURL string

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

	serveCmd.Flags().StringVar(&grafanaPublicURL, "grafana-public-url", "", "Browsable Grafana URL, if not provided, the API url will be used. This is the base url for iframes embedding.")
	serveCmd.Flags().StringVar(&grafanaApiURL, "grafana-api-url", "http://localhost:3000", "Grafana API URL")
	serveCmd.Flags().StringVar(&grafanaUser, "grafana-user", "admin", "Grafana user")
	serveCmd.Flags().StringVar(&grafanaPassword, "grafana-password", "", "Grafana password")

	serveCmd.Flags().StringVar(&prometheusURL, "prometheus-url", "http://localhost:9090", "Prometheus server URL")

	webCmd.AddCommand(serveCmd)
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

	app, err := web.NewApp(ctx, config)
	if err != nil {
		log.Fatal("Failed to create the web application instance: ", err)
	}

	if err := app.Start(ctx); err != nil {
		log.Fatal("Error while running the web application server: ", err)
	}
}
