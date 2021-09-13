package discovery

import (
	"fmt"

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

func NewSAPSystemsDiscovery(client consul.Client) SAPSystemsDiscovery {
	r := SAPSystemsDiscovery{}
	r.id = SAPDiscoveryId
	r.discovery = NewDiscovery(client)
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
		err := s.Store(d.discovery.client)
		if err != nil {
			return "", err
		}
	}

	// Store SAP System on hosts metadata
	err = storeSAPSystemTags(d.discovery.client, d.SAPSystems)
	if err != nil {
		return "", err
	}

	sysNames := systems.GetSIDsString()
	if sysNames != "" {
		return fmt.Sprintf("SAP system(s) with ID: %s discovered", sysNames), nil
	}
	output := "No SAP systems were found (you possibly need to run the trento agent with root privileges)"
	return output, nil
}

func storeSAPSystemTags(client consul.Client, systems sapsystem.SAPSystemsList) error {
	sysNames := systems.GetSIDsString()

	// Store host metadata
	metadata := hosts.Metadata{
		SAPSystems: sysNames,
	}

	if err := metadata.Store(client); err != nil {
		return err
	}

	return nil
}
