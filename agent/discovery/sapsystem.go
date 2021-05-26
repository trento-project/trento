package discovery

import (
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/environments"
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

func (s SAPSystemsDiscovery) GetId() string {
	return s.id
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

		// Store sap instance name on hosts metadata
		err = storeSAPSystemTag(
			d.discovery.client,
			s.Properties["SAPSYSTEMNAME"].Value,
			s.Type)
		if err != nil {
			return err
		}
	}

	return nil
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

func storeSAPSystemTag(client consul.Client, sid, systemType string) error {
	env, land, sys, err := getCurrentEnvironment(client, sid)
	if err != nil {
		return err
	}

	// Create a new ungrouped entry
	if env == consul.KvUngrouped {
		newEnv := environments.NewEnvironment(env, land, sys)
		newEnv.Landscapes[land].SAPSystems[sys].Type = systemType
		err := newEnv.Store(client)
		if err != nil {
			return err
		}
	}

	// Store host metadata
	metadata := hosts.Metadata{
		Environment: env,
		Landscape:   land,
		SAPSystem:   sys,
	}

	err = metadata.Store(client)
	if err != nil {
		return err
	}

	return nil
}
