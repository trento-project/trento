package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/trento-project/trento/cmd/agent"
	"github.com/trento-project/trento/cmd/web"
	"github.com/trento-project/trento/internal"
)

var cfgFile string
var logLevel string

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
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.trento.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "then minimum severity (error, warn, info, debug) of logs to output")
	rootCmd.AddCommand(web.NewWebCmd())
	rootCmd.AddCommand(agent.NewAgentCmd())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".trento" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".trento")
	}

	if logLevel == "" {
		logLevel = "info"
	}
	internal.SetLogLevel(logLevel)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
