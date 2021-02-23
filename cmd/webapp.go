package cmd

import (
	"github.com/spf13/cobra"
)

// webappCmd represents the webapp command
var webappCmd = &cobra.Command{
	Use:   "webapp",
	Short: "Command tree related to the web application component",
}

func init() {
	rootCmd.AddCommand(webappCmd)
}
