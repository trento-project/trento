package discover

import (
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/sapsystem"
	"github.com/trento-project/trento/internal/environments"
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

		// Store sap instance name on hosts metadata
		err = storeSAPSystemTag(discover.host.client, s.Properties["SAPSYSTEMNAME"].Value)
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

// These methods must go here. We cannot put them in the internal/sapsystem.go package
// as this creates potential cyclical imports
func getCurrentEnvironment(client consul.Client, sid string) (string, string, string, error) {
	var env string = consul.KvUngrouped
	var land string = consul.KvUngrouped
	var sys string = sid

	envs, err := environments.Load(client)
	if err != nil {
		return env, land, sys, err
	}
	for envKey, envValue := range envs {
		for landKey, landValue := range envValue.Landscapes {
			for sysKey, _ := range landValue.SAPSystems {
				if sysKey == sys {
					env = envKey
					land = landKey
					break
				}
			}
		}
	}
	return env, land, sys, err
}

func storeSAPSystemTag(client consul.Client, sid string) error {
	env, land, sys, err := getCurrentEnvironment(client, sid)
	if err != nil {
		return err
	}

	// Create a new ungrouped entry
	if env == consul.KvUngrouped {
		newEnv := environments.NewEnvironment(env, land, sys)
		err := newEnv.Store(client)
		if err != nil {
			return err
		}
	}

	// Store host metadata
	metadata := environments.NewMetadata()
	metadata.Environment = env
	metadata.Landscape = land
	metadata.SAPSystem = sys
	err = metadata.Store(client)
	if err != nil {
		return err
	}

	return nil
}
