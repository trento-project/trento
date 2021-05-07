package sapsystem

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure" //MIT license, is this a problem?
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/environments"
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

	// Store sap instance name on hosts metadata
	err = s.storeSAPSystemTag(client)
	if err != nil {
		return errors.Wrap(err, "Error storing the SAP system data in the environments tree")
	}

	return nil
}


func (s *SAPSystem) getCurrentEnvironment(client consul.Client) (string, string, string, error) {
	var env string = consul.KvUngrouped
	var land string = consul.KvUngrouped
	var sys string = s.Properties["SAPSYSTEMNAME"].Value

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

func (s *SAPSystem) storeSAPSystemTag(client consul.Client) error {
	env, land, sys, err := s.getCurrentEnvironment(client)
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
