package agent

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/trento-project/trento/agent"
	"github.com/trento-project/trento/internal/ruleset"
)

var TTL time.Duration
var port int
var consulConfigDir string

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
	startCmd.Flags().IntVarP(&port, "port", "p", 8700, "The TCP port to use for the web service")
	startCmd.Flags().StringVarP(&consulConfigDir, "consul-config-dir", "", "consul.d", "Consul configuration directory used to store node meta-data")

	agentCmd.AddCommand(startCmd)
	agentCmd.AddCommand(runOnceCmd)

	return agentCmd
}

func runOnce(cmd *cobra.Command, args []string) {
	var err error

	ruleSet, err := ruleset.NewRuleSet(args)
	if err != nil {
		log.Fatal("could not load embedded rulesets", err)
	}

	ruleSetsData, err := ruleSet.GetRulesets()
	if err != nil {
		log.Fatal("could not get rulesets data", err)
	}

	checker, err := agent.NewChecker(ruleSetsData)
	if err != nil {
		log.Fatal("Failed to create a Checker instance: ", err)
	}

	res, err := checker()
	if err != nil {
		log.Fatal("Failed to do checks: ", err)
	}

	res.CheckPrettyPrint(os.Stdout)
}

func start(cmd *cobra.Command, args []string) {
	var err error

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := agent.DefaultConfig()
	if err != nil {
		log.Fatal("Failed to create the agent configuration: ", err)
	}

	cfg.DefinitionsPaths = args
	cfg.WebPort = port
	cfg.CheckerTTL = TTL
	cfg.ConsulConfigDir = consulConfigDir

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
	for _, definitionsPath := range args {
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
	}

	return nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
