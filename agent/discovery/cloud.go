package discovery

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/agent/discovery/collector"
	"github.com/trento-project/trento/internal/cloud"
)

const CloudDiscoveryId string = "cloud_discovery"

type CloudDiscovery struct {
	id              string
	collectorClient collector.Client
	interval        time.Duration
}

func NewCloudDiscovery(collectorClient collector.Client, config DiscoveriesConfig) (Discovery, error) {
	if config.DiscoveriesPeriodsConfig.Cloud < 1 {
		return nil, fmt.Errorf("invalid interval %s", config.DiscoveriesPeriodsConfig.Cloud)
	}

	d := CloudDiscovery{}
	d.collectorClient = collectorClient
	d.id = CloudDiscoveryId
	d.interval = config.DiscoveriesPeriodsConfig.Cloud

	return d, nil
}

func (d CloudDiscovery) GetId() string {
	return d.id
}

func (d CloudDiscovery) GetInterval() time.Duration {
	return d.interval
}

func (d CloudDiscovery) Discover() (string, error) {
	cloudData, err := cloud.NewCloudInstance()
	if err != nil {
		return "", err
	}

	err = d.collectorClient.Publish(d.id, cloudData)
	if err != nil {
		log.Debugf("Error while sending cloud discovery to data collector: %s", err)
		return "", err
	}

	if cloudData.Provider == "" {
		return "No cloud provider discovered on this host", nil
	}

	return fmt.Sprintf("Cloud provider %s discovered", cloudData.Provider), nil
}
