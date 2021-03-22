package agent

import (
	"github.com/spf13/cobra"
)

func NewAgentCmd() *cobra.Command {

	cmdRegister := &cobra.Command{
		Use:   "register",
		Short: "Register the agent in the system",
	}

	cmdStart := &cobra.Command{
		Use:   "start",
		Short: "Start the agent",
	}

	cmdAgent := &cobra.Command{
		Use:   "agent",
		Short: "Intermediate agent command",
	}

	cmdAgent.AddCommand(cmdRegister)
	cmdAgent.AddCommand(cmdStart)

	return cmdAgent
}
