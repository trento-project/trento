package discovery

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/agent/discovery/collector"
	"github.com/trento-project/trento/internal/sapsystem"
)

const SAPDiscoveryId string = "sap_system_discovery"

type SAPSystemsDiscovery struct {
	id              string
	collectorClient collector.Client
	host            string
	interval        time.Duration
}

func NewSAPSystemsDiscovery(collectorClient collector.Client, config DiscoveriesConfig) (Discovery, error) {
	if config.DiscoveriesPeriodsConfig.SAPSystem < 1 {
		return nil, fmt.Errorf("invalid interval %s", config.DiscoveriesPeriodsConfig.SAPSystem)
	}

	r := SAPSystemsDiscovery{}
	r.id = SAPDiscoveryId
	r.collectorClient = collectorClient
	r.host, _ = os.Hostname()
	r.interval = config.DiscoveriesPeriodsConfig.SAPSystem

	return r, nil
}

func (d SAPSystemsDiscovery) GetId() string {
	return d.id
}

func (d SAPSystemsDiscovery) GetInterval() time.Duration {
	return d.interval
}

func (d SAPSystemsDiscovery) Discover() (string, error) {
	systems, err := sapsystem.NewSAPSystemsList()

	if err != nil {
		return "", err
	}

	err = d.collectorClient.Publish(d.id, systems)
	if err != nil {
		log.Debugf("Error while sending sapsystem discovery to data collector: %s", err)
		return "", err
	}

	sysNames := systems.GetSIDsString()
	if sysNames != "" {

		return fmt.Sprintf("SAP system(s) with ID: %s discovered", sysNames), nil
	}

	return "No SAP system discovered on this host", nil
}
