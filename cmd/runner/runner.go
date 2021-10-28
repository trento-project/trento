package runner

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/trento-project/trento/runner"
)

var araServer string
var apiHost string
var apiPort int
var interval int
var ansibleFolder string

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
	startCmd.Flags().StringVar(&apiHost, "api-host", "0.0.0.0", "Trento web server API host")
	startCmd.Flags().IntVar(&apiPort, "api-port", 8080, "Trento web server API port")
	startCmd.Flags().IntVarP(&interval, "interval", "i", 5, "Interval in minutes to run the checks")
	startCmd.Flags().StringVar(&ansibleFolder, "ansible-folder", "/tmp/trento", "Folder where the ansible file structure will be created")

	// Bind the flags to viper and make them available to the application
	startCmd.Flags().VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag(f.Name, f)
	})

	runnerCmd.AddCommand(startCmd)

	return runnerCmd
}

func start(cmd *cobra.Command, args []string) {
	var err error

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	config := LoadConfig()

	runner, err := runner.NewRunner(config)
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
