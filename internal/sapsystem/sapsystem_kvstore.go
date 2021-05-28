package sapsystem

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/trento-project/trento/internal/environments"
	"github.com/trento-project/trento/internal/hosts"

	"github.com/trento-project/trento/internal/consul"
)

func (s *SAPSystem) getKVPath() string {
	host, _ := os.Hostname()
	key := fmt.Sprintf(consul.KvHostsSAPSystemPath, host)
	name := s.Properties["SAPSYSTEMNAME"].Value
	kvPath := fmt.Sprintf("%s/%s", key, name)

	return kvPath
}

func (s *SAPSystem) getKVMetadataPath() string {
	host, _ := os.Hostname()
	kvPath := fmt.Sprintf(consul.KvHostsMetadataPath, host)

	return kvPath
}

func (s *SAPSystem) Store(client consul.Client) error {
	kvPath := s.getKVPath()

	// Clean the current data before storing the new values
	_, err := client.KV().DeleteTree(kvPath, nil)
	if err != nil {
		return errors.Wrap(err, "Error deleting SAP system content")
	}

	systemMap := make(map[string]interface{})
	mapstructure.Decode(s, &systemMap)

	err = client.KV().PutMap(kvPath, systemMap)
	if err != nil {
		return errors.Wrap(err, "Error storing a SAP instance")
	}

	return nil
}

// Load from KV storage

func Load(client consul.Client, host string) (map[string]*SAPSystem, error) {
	var sapSystems = map[string]*SAPSystem{}

	kvPath := fmt.Sprintf(consul.KvHostsSAPSystemPath, host)

	entries, err := client.KV().ListMap(kvPath, kvPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for SAP systems KV values")
	}

	for sys, sysValue := range entries {
		system := &SAPSystem{}
		mapstructure.Decode(sysValue, &system)
		sapSystems[sys] = system
	}

	return sapSystems, nil
}

func (s *SAPSystem) StoreSAPSystemTags(client consul.Client) error {
	sid := s.GetSID()

	envName, landName, sysName, err := loadSAPSystemTags(client, sid)
	if err != nil {
		return err
	}

	// If we didn't find any environment, we create a new default one
	if envName == "" {
		land := environments.NewDefaultLandscape()
		land.AddSystem(environments.NewSystem(sysName, s.Type))
		env := environments.NewDefaultEnvironment()
		env.AddLandscape(land)

		err := env.Store(client)
		if err != nil {
			return err
		}
		envName = env.Name
		landName = land.Name
	}

	// Store host metadata
	metadata := hosts.Metadata{
		Environment: envName,
		Landscape:   landName,
		SAPSystem:   sysName,
	}

	err = metadata.Store(client)
	if err != nil {
		return err
	}

	return nil
}

// These methods must go here. We cannot put them in the internal/sapsystem.go package
// as this creates potential cyclical imports
func loadSAPSystemTags(client consul.Client, sid string) (string, string, string, error) {
	var env, land string
	sys := sid

	envs, err := environments.Load(client)
	if err != nil {
		return env, land, sys, err
	}
	for envKey, envValue := range envs {
		for landKey, landValue := range envValue.Landscapes {
			for sysKey := range landValue.SAPSystems {
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
