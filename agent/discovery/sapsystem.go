package discovery

import (
	"github.com/trento-project/trento/internal/consul"
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

func (d SAPSystemsDiscovery) Discover() error {
	systems, err := sapsystem.NewSAPSystemsList()

	if err != nil {
		return err
	}

	d.SAPSystems = systems
	for _, s := range d.SAPSystems {
		err := s.Store(d.discovery.client)
		if err != nil {
			return err
		}

		// Store SAP System, Landscape and Environment names on hosts metadata
		err = s.StoreSAPSystemTags(d.discovery.client)
		if err != nil {
			return err
		}
	}

	return nil
}
