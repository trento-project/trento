package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/SUSE/console-for-sap-applications/cmd/webapp"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "console-for-sap",
	Short: "A cloud-native, web application to manage OS-related tasks for SAP Applications.",
	Long: `SUSE Console for SAP Applications is a web-based graphical user interface 
that can help you deploy, provision and operate infrastructure for SAP Applications`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.console-for-sap.yaml)")
	rootCmd.AddCommand(webapp.NewWebappCmd())
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

		// Search config in home directory with name ".console-for-sap" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".console-for-sap")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
