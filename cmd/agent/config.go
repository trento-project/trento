package agent

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/trento-project/trento/agent"
	"github.com/trento-project/trento/agent/discovery/collector"
)

func LoadConfig() (*agent.Config, error) {
	enablemTLS := viper.GetBool("enable-mtls")
	cert := viper.GetString("cert")
	key := viper.GetString("key")
	ca := viper.GetString("ca")

	if enablemTLS {
		var err error

		if cert == "" {
			err = fmt.Errorf("you must provide a server ssl certificate")
		}
		if key == "" {
			err = errors.Wrap(err, "you must provide a key to enable mTLS")
		}
		if ca == "" {
			err = errors.Wrap(err, "you must provide a CA ssl certificate")
		}
		if err != nil {
			return nil, err
		}
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "could not read the hostname")
	}

	return &agent.Config{
		CollectorConfig: &collector.Config{
			CollectorHost: viper.GetString("collector-host"),
			CollectorPort: viper.GetInt("collector-port"),
			EnablemTLS:    enablemTLS,
			Cert:          cert,
			Key:           key,
			CA:            ca,
		},
		ConsulConfigDir: viper.GetString("consul-config-dir"),
		InstanceName:    hostname,
		SSHAddress:      viper.GetString("ssh-address"),
		DiscoveryPeriod: time.Duration(viper.GetInt("discovery-period")) * time.Minute,
	}, nil
}
