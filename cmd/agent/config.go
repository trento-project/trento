package agent

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/trento-project/trento/agent"
	"github.com/trento-project/trento/agent/discovery"
	"github.com/trento-project/trento/agent/discovery/collector"
)

func validatePeriod(durationFlag string, minValue time.Duration) error {
	period := viper.GetDuration(durationFlag)
	if period < minValue {
		return errors.Errorf("%s: invalid interval %s, should be at least %s", durationFlag, period, minValue)
	}

	return nil
}

func LoadConfig() (*agent.Config, error) {
	enablemTLS := viper.GetBool("enable-mtls")
	cert := viper.GetString("cert")
	key := viper.GetString("key")
	ca := viper.GetString("ca")

	minPeriodValues := map[string]time.Duration{
		"cluster-discovery-period":      discovery.ClusterDiscoveryMinPeriod,
		"sapsystem-discovery-period":    discovery.SAPDiscoveryMinPeriod,
		"cloud-discovery-period":        discovery.CloudDiscoveryMinPeriod,
		"host-discovery-period":         discovery.HostDiscoveryMinPeriod,
		"subscription-discovery-period": discovery.SubscriptionDiscoveryMinPeriod,
	}

	for flagName, minPeriodValue := range minPeriodValues {
		err := validatePeriod(flagName, minPeriodValue)
		if err != nil {
			return nil, err
		}
	}

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

	sshAddress := viper.GetString("ssh-address")
	if sshAddress == "" {
		return nil, errors.New("ssh-address is required, cannot start agent")
	}

	collectorConfig := &collector.Config{
		CollectorHost: viper.GetString("collector-host"),
		CollectorPort: viper.GetInt("collector-port"),
		EnablemTLS:    enablemTLS,
		Cert:          cert,
		Key:           key,
		CA:            ca,
	}

	discoveryPeriodsConfig := &discovery.DiscoveriesPeriodConfig{
		Cluster:      viper.GetDuration("cluster-discovery-period"),
		SAPSystem:    viper.GetDuration("sapsystem-discovery-period"),
		Cloud:        viper.GetDuration("cloud-discovery-period"),
		Host:         viper.GetDuration("host-discovery-period"),
		Subscription: viper.GetDuration("subscription-discovery-period"),
	}

	discoveriesConfig := &discovery.DiscoveriesConfig{
		SSHAddress:               sshAddress,
		CollectorConfig:          collectorConfig,
		DiscoveriesPeriodsConfig: discoveryPeriodsConfig,
	}

	return &agent.Config{
		InstanceName:      hostname,
		DiscoveriesConfig: discoveriesConfig,
	}, nil
}
