package agent

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/trento-project/trento/agent"
	"github.com/trento-project/trento/internal"
)

func NewAgentCmd() *cobra.Command {
	var bindIP string
	var consulConfigDir string
	var discoveryPeriod int

	var collectorHost string
	var collectorPort int

	var enablemTLS bool
	var cert string
	var key string
	var ca string

	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Command tree related to the agent component",
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the agent",
		Run:   start,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return internal.InitConfig("agent")
		},
	}

	startCmd.Flags().StringVar(&bindIP, "bind-ip", "", "Private address to which the trento-agent should be bound for internal communications. Reachable by the other hosts, including the trento server.")
	startCmd.MarkFlagRequired("bind-ip")

	startCmd.Flags().StringVarP(&consulConfigDir, "consul-config-dir", "", "consul.d", "Consul configuration directory used to store node meta-data")
	startCmd.Flags().IntVarP(&discoveryPeriod, "discovery-period", "", 2, "Discovery mechanism loop period on minutes")

	startCmd.Flags().StringVar(&collectorHost, "collector-host", "localhost", "Data Collector host")
	startCmd.Flags().IntVar(&collectorPort, "collector-port", 8081, "Data Collector port")

	startCmd.Flags().BoolVar(&enablemTLS, "enable-mtls", false, "Enable mTLS authentication between server and agent")
	startCmd.Flags().StringVar(&cert, "cert", "", "mTLS client certificate")
	startCmd.Flags().StringVar(&key, "key", "", "mTLS client key")
	startCmd.Flags().StringVar(&ca, "ca", "", "mTLS Certificate Authority")

	// Bind the flags to viper and make them available to the application
	startCmd.Flags().VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag(f.Name, f)
	})

	agentCmd.AddCommand(startCmd)

	return agentCmd
}

func start(cmd *cobra.Command, args []string) {
	var err error

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Failed to create the agent configuration: ", err)
	}

	a, err := agent.NewAgent(config)
	if err != nil {
		log.Fatal("Failed to create the agent: ", err)
	}

	go func() {
		quit := <-signals
		log.Printf("Caught %s signal!", quit)

		log.Println("Stopping the agent...")
		a.Stop()
	}()

	log.Println("Starting the Console Agent...")
	err = a.Start()
	if err != nil {
		log.Fatal("Failed to start the agent: ", err)
	}
}
