package discovery

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/agent/collector"
	"github.com/trento-project/trento/internal/cloud"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
)

const CloudDiscoveryId string = "cloud_discovery"

type CloudDiscovery struct {
	id        string
	discovery BaseDiscovery
}

func NewCloudDiscovery(consulClient consul.Client, collectorClient collector.Client) CloudDiscovery {
	r := CloudDiscovery{}
	r.id = CloudDiscoveryId
	r.discovery = NewDiscovery(consulClient, collectorClient)
	return r
}

func (d CloudDiscovery) GetId() string {
	return d.id
}

func (d CloudDiscovery) Discover() (string, error) {
	cloudData, err := cloud.NewCloudInstance()
	if err != nil {
		return "", err
	}

	err = cloudData.Store(d.discovery.consulClient)
	if err != nil {
		return "", err
	}

	err = storeCloudMetadata(d.discovery.consulClient, cloudData.Provider)
	if err != nil {
		return "", err
	}

	err = d.discovery.collectorClient.Publish(d.id, cloudData)
	if err != nil {
		log.Debugf("Error while sending cloud discovery to data collector: %s", err)
		return "", err
	}

	if cloudData.Provider == "" {
		return "No cloud provider discovered on this host", nil
	}

	return fmt.Sprintf("Cloud provider %s discovered", cloudData.Provider), nil
}

func storeCloudMetadata(client consul.Client, cloudProvider string) error {
	metadata := hosts.Metadata{
		CloudProvider: cloudProvider,
	}
	err := metadata.Store(client)
	if err != nil {
		return err
	}

	return nil
}
