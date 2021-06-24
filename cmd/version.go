package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Trento",
	Long:  `All software has versions. This is Trento's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Trento version %s\nbuilt with %s %s/%s\n", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	},
}
