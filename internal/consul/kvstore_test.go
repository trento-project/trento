package consul

import (
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/internal/consul/mocks"
)

func TestKVMaps(t *testing.T) {

  kvpairs := consulApi.KVPairs{
		&consulApi.KVPair{Key: "/trento/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/env1/", Value: []byte("")},
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

	kv := new(KV)

	consulApi.On("KV").Return(kv)

	kv.On("List", "/trento", (*consulApi.QueryOptions)(nil)).Return(kvpairs, nil, nil)

  result, _ := Maps(kv, "/trento/", "/trento")

  expectedMap := map[string]interace{"a": 1}

  assert.Equal(t, result, expectedMap)
}
