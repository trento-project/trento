package agent

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/SUSE/console-for-sap-applications/agent"
)

func NewAgentCmd() *cobra.Command {

	cmdRegister := &cobra.Command{
		Use:   "register",
		Short: "Register the agent in the system",
	}

	cmdStart := &cobra.Command{
		Use:   "start",
		Short: "Start the agent",
		Run: start,
	}

	cmdAgent := &cobra.Command{
		Use:   "agent",
		Short: "Command tree related to the agent component",
	}

	cmdAgent.AddCommand(cmdRegister)
	cmdAgent.AddCommand(cmdStart)

	return cmdAgent
}

func start(cmd *cobra.Command, args []string) {
	var err error

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	a, err := agent.New()
	if err != nil {
		log.Fatal("Failed to create the agent instance: ", err)
	}

	defer func() {
		quit := <-signals
		log.Printf("Caught %s signal! Stopping the agent...", quit)

		a.Stop()

		log.Println("Agent stopped.")
	}()

	go func() {
		err = a.Start()
		if err != nil {
			log.Fatal("Failed to start the agent: ", err)
		}
	}()
}
