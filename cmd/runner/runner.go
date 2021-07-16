package runner

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trento-project/trento/runner"
)

var araServer string
var interval int
var ansibleFolder string
var consulAddr string

func NewRunnerCmd() *cobra.Command {
	runnerCmd := &cobra.Command{
		Use:   "runner",
		Short: "Command tree related to the runner component",
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Starts the runner process. This process takes care of running the checks",
		Run:   start,
	}

	startCmd.Flags().StringVar(&araServer, "ara-server", "http://127.0.0.1:8000", "ARA server url (ex: http://localhost:8000)")
	startCmd.Flags().StringVar(&consulAddr, "consul-addr", "127.0.0.1:8500", "Consul host address (ex: localhost:8500)")
	startCmd.Flags().IntVarP(&interval, "interval", "i", 5, "Interval in minutes to run the checks")
	startCmd.Flags().StringVar(&ansibleFolder, "ansible-folder", "/usr/etc/trento", "Folder where the ansible file structure will be created")

	runnerCmd.AddCommand(startCmd)

	return runnerCmd
}

func start(cmd *cobra.Command, args []string) {
	var err error

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := runner.DefaultConfig()
	if err != nil {
		log.Fatal("Failed to create the runner configuration: ", err)
	}

	cfg.AraServer = araServer
	cfg.ConsulAddr = consulAddr
	cfg.Interval = time.Duration(interval) * time.Minute
	cfg.AnsibleFolder = ansibleFolder
	cfg.ConsulTemplateLogLevel = viper.GetString("log-level")

	runner, err := runner.NewWithConfig(cfg)
	if err != nil {
		log.Fatal("Failed to create the runner instance: ", err)
	}

	go func() {
		quit := <-signals
		log.Printf("Caught %s signal!", quit)

		log.Println("Stopping the runner...")
		runner.Stop()
	}()

	err = runner.Start()
	if err != nil {
		log.Fatal("Failed to start the runner service: ", err)
	}
}
