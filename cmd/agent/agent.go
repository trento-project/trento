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
var rules []string

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
		Use:   "start",
		Short: "Start the agent",
		Run:   start,
	}
	startCmd.Flags().DurationVar(&TTL, "ttl", time.Second*10, "Duration of Consul TTL checks")
	startCmd.Flags().IntVarP(&port, "port", "p", 8700, "The TCP port to use for the web service")
	startCmd.Flags().StringVarP(&consulConfigDir, "consul-config-dir", "", "consul.d", "Consul configuration directory used to store node meta-data")

	runOnceCmd.Flags().StringSliceVar(&rules, "rulesets", []string{}, "User defined rulesets. This flag can be used multiple times")

	agentCmd.AddCommand(startCmd)
	agentCmd.AddCommand(runOnceCmd)

	return agentCmd
}

func runOnce(cmd *cobra.Command, args []string) {
	var err error

	ruleSet, err := ruleset.NewRuleSets(rules)
	if err != nil {
		log.Fatal("could not load embedded rulesets", err)
	}

	rulesetsYaml, err := ruleSet.GetRulesetsYaml()
	if err != nil {
		log.Println("An error occurred while generating the rulesets:", err)
		return
	}

	res, err := agent.NewCheckResult(rulesetsYaml)
	if err != nil {
		log.Println("An error occurred while running health checks:", err)
		return
	}
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
	for _, definitionsPath := range rules {
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
