package agent

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/SUSE/console-for-sap-applications/agent"
)

var TTL time.Duration
var serviceName string

func NewAgentCmd() *cobra.Command {

	startCmd := &cobra.Command{
		Use:   "start path/to/definitions.yaml",
		Short: "Start the agent",
		Run:   start,
		Args:  startArgsValidator,
	}
	startCmd.Flags().DurationVar(&TTL, "ttl", time.Second*10, "Duration of Consul TTL checks")
	startCmd.Flags().StringVarP(&serviceName, "service-name", "n", "", "The name of the service this agent will monitor")
	must(startCmd.MarkFlagRequired("service-name"))

	cmdAgent := &cobra.Command{
		Use:   "agent",
		Short: "Command tree related to the agent component",
	}

	cmdAgent.AddCommand(startCmd)

	return cmdAgent
}

func start(cmd *cobra.Command, args []string) {
	var err error

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := agent.DefaultConfig()
	if err != nil {
		log.Fatal("Failed to create the agent configuration: ", err)
	}

	cfg.DefinitionsPath = args[0]
	cfg.ServiceName = serviceName
	cfg.TTL = TTL

	a, err := agent.NewWithConfig(cfg)
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

func startArgsValidator(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("accepts exactly 1 argument, received %d", len(args))
	}

	definitionsPath := args[0]

	info, err := os.Lstat(definitionsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("unable to find file %q", definitionsPath)
		}
		return fmt.Errorf("error when running os.Lstat(%q): %s", definitionsPath, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%q is a directory", definitionsPath)
	}

	return nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
