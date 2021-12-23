package runner

import (
	"time"

	"github.com/spf13/viper"
	"github.com/trento-project/trento/runner"
)

func LoadConfig() *runner.Config {
	return &runner.Config{
		ApiHost:       viper.GetString("api-host"),
		ApiPort:       viper.GetInt("api-port"),
		Interval:      time.Duration(interval) * time.Minute,
		AnsibleFolder: viper.GetString("ansible-folder"),
	}
}
