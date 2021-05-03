package consul

import (
	"strings"

	"github.com/pkg/errors"

	consulApi "github.com/hashicorp/consul/api"
)

// The Trento Agent is periodically updating data structures in the Consul Key-Value Store
// to be consumed by the Trento Web Console
//
// At some point in time this structure needs to become a versioned protocol between the
// Trento Web Console and the Trento Agent

// For now, until we tighten the structure around versioning, this file is a start
// that records constants that need to be in sync or have compat behavior between
// Console and Agent

const KvClustersPath string = "trento/v0/clusters"
const KvHostsPath string = "trento/v0/hosts"
const KvEnvironmentsPath string = "trento/v0/environments"

type ClusterStonithType int

const (
	ClusterStonithNone ClusterStonithType = iota
	ClusterStonithSBD
	ClusterStonithUnknown
)

// Method to convert the KVs output to a Map
// prefix -> The KV prefix to get the data
// offset -> Remove the offset from the outgoing map
// 	         (used to remove initial key from the map like trento/v0/)
func Maps(k *consulApi.KV, prefix, offset string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	entries, _, err := k.List(prefix, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for KV values")
	}

	currentItem := result

	for _, entry := range entries {
		modEntry := strings.TrimPrefix(strings.Trim(entry.Key, " "), offset)
		if len(modEntry) == 0 {
			continue
		}

		currentItem = result
		keys := strings.Split(modEntry, "/")
		for i, key := range keys {

			if len(key) == 0 {
				break
			} else if i == len(keys)-1 {
				currentItem[key] = string(entry.Value)
				break
			}

			if _, ok := currentItem[key]; !ok {
				currentItem[key] = make(map[string]interface{})
			}

			item, _ := currentItem[key].(map[string]interface{})
			currentItem = item
		}
	}

	return result, nil
}
