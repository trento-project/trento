package consul

import (
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

func TestKVMaps(t *testing.T) {

	kvPairs := consulApi.KVPairs{
		&consulApi.KVPair{Key: "/trento/", Value: []byte("foo")},
		&consulApi.KVPair{Key: "/trento/env1/", Value: []byte("bar")},
		&consulApi.KVPair{Key: "/trento/env1/landscapes/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/env1/landscapes/land1/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/env1/landscapes/land1/sapsystems/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/env1/landscapes/land1/sapsystems/sys1/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/env1/landscapes/land1/sapsystems/sys1/name", Value: []byte("name")},
		&consulApi.KVPair{Key: "/trento/env1/landscapes/land2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/env1/landscapes/land2/sapsystems/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/env1/landscapes/land2/sapsystems/sys2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/env2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/env2/landscapes/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/env2/landscapes/land3/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/env2/landscapes/land3/sapsystems/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/env2/landscapes/land3/sapsystems/sys3/", Value: []byte("")},
	}

	kv := &kv{list: func(prefix string, q *consulApi.QueryOptions) (consulApi.KVPairs, *consulApi.QueryMeta, error) {
		return kvPairs, nil, nil
	}}
	result, _ := kv.Maps("", "")

	expectedMap := map[string]interface{}{"/trento/": "foo"}

	assert.Equal(t, expectedMap, result)
}
