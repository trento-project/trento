package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/spf13/viper"

	"github.com/trento-project/trento/cmd/agent"
	"github.com/trento-project/trento/cmd/runner"
	"github.com/trento-project/trento/cmd/web"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "trento",
	Short: "An open cloud-native web console improving on the life of SAP Applications administrators.",
	Long: `Trento is a web-based graphical user interface
that can help you deploy, provision and operate infrastructure for SAP Applications`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	var cfgFile string
	var logLevel string

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.trento.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "then minimum severity (error, warn, info, debug) of logs to output")

	// Make global flags available in the children commands
	rootCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag(f.Name, f)
	})

	rootCmd.AddCommand(web.NewWebCmd())
	rootCmd.AddCommand(agent.NewAgentCmd())
	rootCmd.AddCommand(runner.NewRunnerCmd())
}
