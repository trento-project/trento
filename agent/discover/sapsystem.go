package discover

import (
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/sapsystem"
)

const SAPDiscoverId string = "discover_sap"

type SAPSystemsDiscover struct {
	id         string
	host       Discover
	SAPSystems sapsystem.SAPSystemsList
}

func (s SAPSystemsDiscover) GetId() string {
	return s.id
}

func (s SAPSystemsDiscover) ShouldDiscover(client consul.Client) bool {
	return true
}

func (s SAPSystemsDiscover) storeDiscovery(cStorePath, cStoreValue string) error {
	return nil
}

func (discover SAPSystemsDiscover) Discover() error {
	systems, err := sapsystem.NewSAPSystemsList()

	if err != nil {
		return err
	}

	discover.SAPSystems = systems
	for _, s := range discover.SAPSystems {
		err := s.Store(discover.host.client)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewSAPSystemsDiscover(client consul.Client) SAPSystemsDiscover {
	r := SAPSystemsDiscover{}
	r.id = SAPDiscoverId
	r.host = NewDiscover(client)
	return r
}
