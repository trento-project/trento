package discovery

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/agent/discovery/collector"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/internal/sapsystem"
)

const SAPDiscoveryId string = "sap_system_discovery"

type SAPSystemsDiscovery struct {
	id         string
	discovery  BaseDiscovery
	SAPSystems sapsystem.SAPSystemsList
}

func NewSAPSystemsDiscovery(consulClient consul.Client, collectorClient collector.Client) SAPSystemsDiscovery {
	r := SAPSystemsDiscovery{}
	r.id = SAPDiscoveryId
	r.discovery = NewDiscovery(consulClient, collectorClient)
	return r
}

func (d SAPSystemsDiscovery) GetId() string {
	return d.id
}

func (d SAPSystemsDiscovery) Discover() (string, error) {
	systems, err := sapsystem.NewSAPSystemsList()

	if err != nil {
		return "", err
	}

	d.SAPSystems = systems
	for _, s := range d.SAPSystems {
		err := s.Store(d.discovery.consulClient)
		if err != nil {
			return "", err
		}
	}

	// Store SAP System on hosts metadata
	err = storeSAPSystemTags(d.discovery.consulClient, d.SAPSystems)
	if err != nil {
		return "", err
	}

	err = d.discovery.collectorClient.Publish(d.id, systems)
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

func storeSAPSystemTags(client consul.Client, systems sapsystem.SAPSystemsList) error {
	sysNames := systems.GetSIDsString()
	sysIds := systems.GetIDsString()
	sysTypes := systems.GetTypesString()

	// Store host metadata
	metadata := hosts.Metadata{
		SAPSystems:     sysNames,
		SAPSystemsId:   sysIds,
		SAPSystemsType: sysTypes,
	}

	if err := metadata.Store(client); err != nil {
		return err
	}

	return nil
}
