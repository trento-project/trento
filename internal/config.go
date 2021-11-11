package internal

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// InitConfig intializes the config for the application
// If no config file is provided with the --config flag
// it will look for a config in following locations:
//
// ${context} being one of the supported components: agent|web|runner
//
// Order represents priority
// /etc/trento/${context}.yaml
// /usr/etc/trento/${context}.yaml
// $HOME/.config/trento/${context}.yaml
func InitConfig(configName string) error {
	viper.SetConfigType("yaml")
	SetLogLevel(viper.GetString("log-level"))
	SetLogFormatter("2006-01-02 15:04:05")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.SetEnvPrefix("TRENTO")
	viper.AutomaticEnv() // read in environment variables that match

	cfgFile := viper.GetString("config")
	if cfgFile != "" {
		_, err := os.Stat(cfgFile)

		if err != nil {
			// if a config file has been explicitly provided by --config flag,
			// then we should break if that file does not exist
			return fmt.Errorf("cannot load configuration file: %s %s", cfgFile, err)
		}

		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()

		if err != nil {
			return err
		}

		// if no configuration file was explicitly provided,
		// we should look for a config in the expected locations
		viper.AddConfigPath("/etc/trento/")
		viper.AddConfigPath("/usr/etc/trento/")
		viper.AddConfigPath(path.Join(home, ".config", "trento"))
		viper.SetConfigName(configName)
	}

	err := viper.ReadInConfig()

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return nil
}
