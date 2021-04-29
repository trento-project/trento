package agent

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/trento-project/trento/agent"
)

const metaDataFile string = "trento-config.json"

var TTL time.Duration
var serviceName string
var port int
var templatePath string
var configDir string

func NewAgentCmd() *cobra.Command {

	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Command tree related to the agent component",
	}

	runOnceCmd := &cobra.Command{
		Use:   "run-once",
		Short: "run-once",
		Run:   runOnce,
		Args:  startArgsValidator,
	}

	startCmd := &cobra.Command{
		Use:   "start path/to/definitions.yaml",
		Short: "Start the agent",
		Run:   start,
		Args:  startArgsValidator,
	}
	startCmd.Flags().DurationVar(&TTL, "ttl", time.Second*10, "Duration of Consul TTL checks")
	startCmd.Flags().StringVarP(&serviceName, "service-name", "n", "", "The name of the service this agent will monitor")
	startCmd.Flags().IntVarP(&port, "port", "p", 8700, "The TCP port to use for the web service")
	startCmd.Flags().StringVarP(&templatePath, "consul-template", "", "examples/trento-config.tpl", "consul-template template to populate the node meta-data")
	startCmd.Flags().StringVarP(&configDir, "config-dir", "", "consul.d", "Consul configuration directory used to store the trento meta-data file")

	must(startCmd.MarkFlagRequired("service-name"))

	agentCmd.AddCommand(startCmd)

	agentCmd.AddCommand(runOnceCmd)

	return agentCmd
}

func runOnce(cmd *cobra.Command, args []string) {
	var err error

	checker, err := agent.NewChecker(args[0])
	if err != nil {
		log.Fatal("Failed to create a Checker instance: ", err)
	}

	res, err := checker()
	if err != nil {
		log.Fatal("Failed to do checks: ", err)
	}

	res.PrettyPrint(os.Stdout)
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
	cfg.WebPort = port
	cfg.TTL = TTL
	cfg.TemplateSource = templatePath
	cfg.TemplateDestination = path.Join(configDir, metaDataFile)

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
