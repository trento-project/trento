package consul

import (
	"fmt"
	"strings"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
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

type KV interface {
	Get(key string, q *consulApi.QueryOptions) (*consulApi.KVPair, *consulApi.QueryMeta, error)
	List(prefix string, q *consulApi.QueryOptions) (consulApi.KVPairs, *consulApi.QueryMeta, error)
	Keys(prefix, separator string, q *consulApi.QueryOptions) ([]string, *consulApi.QueryMeta, error)
	Put(p *consulApi.KVPair, q *consulApi.WriteOptions) (*consulApi.WriteMeta, error)
	Maps(prefix, offset string) (map[string]interface{}, error)
}

func newKV(wrapped *consulApi.KV) KV {
	return &kv{
		wrapped,
		wrapped.List,
	}
}

type kv struct {
	wrapped *consulApi.KV
	// we need this dedicated function field because the Maps method depends on it and we want to mock it internally
	list func(prefix string, q *consulApi.QueryOptions) (consulApi.KVPairs, *consulApi.QueryMeta, error)
}

func (k *kv) Get(key string, q *consulApi.QueryOptions) (*consulApi.KVPair, *consulApi.QueryMeta, error) {
	return k.wrapped.Get(key, q)
}

func (k *kv) Keys(prefix, separator string, q *consulApi.QueryOptions) ([]string, *consulApi.QueryMeta, error) {
	return k.wrapped.Keys(prefix, separator, q)
}

func (k *kv) Put(p *consulApi.KVPair, q *consulApi.WriteOptions) (*consulApi.WriteMeta, error) {
	return k.wrapped.Put(p, q)
}

func (k *kv) List(prefix string, q *consulApi.QueryOptions) (consulApi.KVPairs, *consulApi.QueryMeta, error) {
	return k.list(prefix, q)
}

// Maps converts the KVs output to a Map
// prefix -> The KV prefix to get the data
// offset -> Remove the offset from the outgoing map
// 	         (used to remove initial key from the map like trento/v0/)
func (k *kv) Maps(prefix, offset string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	entries, _, err := k.list(prefix, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for KV values")
	}

	if !strings.HasSuffix(offset, "/") {
		offset = fmt.Sprintf("%s/", offset)
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
