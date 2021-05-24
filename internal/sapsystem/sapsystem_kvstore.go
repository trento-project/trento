package sapsystem

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure" //MIT license, is this a problem?
	"github.com/pkg/errors"

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
