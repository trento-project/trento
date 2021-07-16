package checkrunner

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trento-project/trento/checkrunner"
)

var araServer string
var interval int
var ansibleFolder string
var consulAddr string

func NewCheckRunnerCmd() *cobra.Command {
	checkRunnerCmd := &cobra.Command{
		Use:   "checkrunner",
		Short: "Command tree related to the check runner component",
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Starts the check runner",
		Run:   start,
	}

	startCmd.Flags().StringVar(&araServer, "ara-server", "http://127.0.0.1:8000", "ARA server url (ex: http://localhost:8000)")
	startCmd.Flags().StringVar(&consulAddr, "consul-addr", "127.0.0.1:8500", "Consul host address (ex: localhost:8500)")
	startCmd.Flags().IntVarP(&interval, "interval", "i", 5, "Interval in minutes to run the checks")
	startCmd.Flags().StringVar(&ansibleFolder, "ansible-folder", "/usr/etc/trento", "Folder where the ansible file structure will be created")

	checkRunnerCmd.AddCommand(startCmd)

	return checkRunnerCmd
}

func start(cmd *cobra.Command, args []string) {
	var err error

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := checkrunner.DefaultConfig()
	if err != nil {
		log.Fatal("Failed to create the check runner configuration: ", err)
	}

	cfg.AraServer = araServer
	cfg.ConsulAddr = consulAddr
	cfg.Interval = time.Duration(interval) * time.Minute
	cfg.AnsibleFolder = ansibleFolder
	cfg.ConsulTemplateLogLevel = viper.GetString("log-level")

	runner, err := checkrunner.NewWithConfig(cfg)
	if err != nil {
		log.Fatal("Failed to create the check runner instance: ", err)
	}

	go func() {
		quit := <-signals
		log.Printf("Caught %s signal!", quit)

		log.Println("Stopping the checker...")
		runner.Stop()
	}()

	err = runner.Start()
	if err != nil {
		log.Fatal("Failed to start the check runner service: ", err)
	}
}
