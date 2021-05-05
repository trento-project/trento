package sapsystem

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure" //MIT license, is this a problem?
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"
)

func (s *SAPSystem) getKVPath() string {
	host, _ := os.Hostname()
	key := fmt.Sprintf(consul.KvHostsSAPSystemPath, host)
	name := s.Properties["INSTANCE_NAME"].Value
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

func (s *SAPSystem) storeSAPSystemTag(client consul.Client) error {
	// This should be done with the unique ID rather than the SID, as this is not unique
	sid := s.Properties["SAPSYSTEMNAME"].Value

	metadata := fmt.Sprintf("%s/%s", s.getKVMetadataPath(), consul.KvMetadataSAPSystem)
	err := client.KV().PutStr(metadata, sid)
	if err != nil {
		return err
	}

	envs, _, err := client.KV().Keys(consul.KvEnvironmentsPath, "", nil)
	if err != nil {
		return err
	}

	for _, env := range envs {
		if strings.HasSuffix(env, fmt.Sprintf("sapsystems/%s/", sid)) {
			return nil
		}
	}

	err = client.KV().PutStr(
		fmt.Sprintf(consul.KvEnvironmentsSAPSystemPath, consul.KvUngrouped, consul.KvUngrouped, sid), "")
	if err != nil {
		return err
	}

	//This 2 next operations should be done most probably somewhere else
	envMetadata := fmt.Sprintf("%s/%s", s.getKVMetadataPath(), consul.KvMetadataSAPEnvironment)
	err = client.KV().PutStr(envMetadata, consul.KvUngrouped)
	if err != nil {
		return err
	}

	landMetadata := fmt.Sprintf("%s/%s", s.getKVMetadataPath(), consul.KvMetadataSAPLandscape)
	err = client.KV().PutStr(landMetadata, consul.KvUngrouped)
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
