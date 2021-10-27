package helpers

import (
	"strings"

	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("db-integration-tests", true)
	viper.SetDefault("db-host", "localhost")
	viper.SetDefault("db-port", "32432")
	viper.SetDefault("db-user", "postgres")
	viper.SetDefault("db-password", "postgres")
	viper.SetDefault("db-name", "trento_test")

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.SetEnvPrefix("TRENTO")
	viper.AutomaticEnv()
}
