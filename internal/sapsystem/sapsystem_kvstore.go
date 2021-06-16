package sapsystem

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"
)

func (s *SAPSystem) getKVPath() string {
	host, _ := os.Hostname()
	key := fmt.Sprintf(consul.KvHostsSAPSystemPath, host)
	kvPath := path.Join(key, s.SID)

	return kvPath
}

func (s *SAPSystem) getKVMetadataPath() string {
	host, _ := os.Hostname()
	kvPath := fmt.Sprintf(consul.KvHostsMetadataPath, host)

	return kvPath
}

func (s *SAPSystem) Store(client consul.Client) error {
	host, _ := os.Hostname()
	l, err := client.AcquireLockKey(fmt.Sprintf(consul.KvHostsSAPSystemPath, host))
	if err != nil {
		return errors.Wrap(err, "could not lock the kv for sapsystem")
	}
	defer l.Unlock()

	kvPath := s.getKVPath()
	// Clean the current data before storing the new values
	_, err = client.KV().DeleteTree(kvPath, nil)
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

	err := client.WaitLock(kvPath)
	if err != nil {
		return nil, errors.Wrap(err, "error waiting for the lock for sapsystems")
	}

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
