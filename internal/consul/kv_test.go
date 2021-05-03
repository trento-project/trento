package consul

import (
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

func TestKVMaps(t *testing.T) {

	kvPairs := consulApi.KVPairs{
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

	kv := &kv{list: func(prefix string, q *consulApi.QueryOptions) (consulApi.KVPairs, *consulApi.QueryMeta, error) {
		return kvPairs, nil, nil
	}}
	result, _ := kv.Maps("/trento", "")

	expectedMap := map[string]interface{}{
		"trento": map[string]interface{}{
			"env1": map[string]interface{}{
				"landscapes": map[string]interface{}{
					"land1": map[string]interface{}{
						"sapsystems": map[string]interface{}{
							"sys1": map[string]interface{}{
								"name": "name",
							},
						},
					},
					"land2": map[string]interface{}{
						"sapsystems": map[string]interface{}{
							"sys2": map[string]interface{}{},
						},
					},
				},
			},
			"env2": map[string]interface{}{
				"landscapes": map[string]interface{}{
					"land3": map[string]interface{}{
						"sapsystems": map[string]interface{}{
							"sys3": map[string]interface{}{},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expectedMap, result)
}

func TestKVMapsOffset(t *testing.T) {

	kvPairs := consulApi.KVPairs{
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

	kv := &kv{list: func(prefix string, q *consulApi.QueryOptions) (consulApi.KVPairs, *consulApi.QueryMeta, error) {
		return kvPairs, nil, nil
	}}
	result, _ := kv.Maps("/trento", "/trento/")

	expectedMap := map[string]interface{}{
		"env1": map[string]interface{}{
			"landscapes": map[string]interface{}{
				"land1": map[string]interface{}{
					"sapsystems": map[string]interface{}{
						"sys1": map[string]interface{}{
							"name": "name",
						},
					},
				},
				"land2": map[string]interface{}{
					"sapsystems": map[string]interface{}{
						"sys2": map[string]interface{}{},
					},
				},
			},
		},
		"env2": map[string]interface{}{
			"landscapes": map[string]interface{}{
				"land3": map[string]interface{}{
					"sapsystems": map[string]interface{}{
						"sys3": map[string]interface{}{},
					},
				},
			},
		},
	}

	assert.Equal(t, expectedMap, result)
}

func TestKVMapsOffsetNoBackslash(t *testing.T) {

	kvPairs := consulApi.KVPairs{
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

	kv := &kv{list: func(prefix string, q *consulApi.QueryOptions) (consulApi.KVPairs, *consulApi.QueryMeta, error) {
		return kvPairs, nil, nil
	}}
	result, _ := kv.Maps("/trento", "/trento")

	expectedMap := map[string]interface{}{
		"env1": map[string]interface{}{
			"landscapes": map[string]interface{}{
				"land1": map[string]interface{}{
					"sapsystems": map[string]interface{}{
						"sys1": map[string]interface{}{
							"name": "name",
						},
					},
				},
				"land2": map[string]interface{}{
					"sapsystems": map[string]interface{}{
						"sys2": map[string]interface{}{},
					},
				},
			},
		},
		"env2": map[string]interface{}{
			"landscapes": map[string]interface{}{
				"land3": map[string]interface{}{
					"sapsystems": map[string]interface{}{
						"sys3": map[string]interface{}{},
					},
				},
			},
		},
	}

	assert.Equal(t, expectedMap, result)
}

func TestKVMapsSingleEntries(t *testing.T) {

	kvPairs := consulApi.KVPairs{
		&consulApi.KVPair{Key: "/trento/env1/landscapes/land1/sapsystems/sys1/name", Value: []byte("name")},
		&consulApi.KVPair{Key: "/trento/env1/landscapes/land1/sapsystems/sys1/description", Value: []byte("desc")},
		&consulApi.KVPair{Key: "/trento/env1/landscapes/land2/description", Value: []byte("land_desc")},
		&consulApi.KVPair{Key: "/trento/env1/landscapes/land2/sapsystems/sys2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/env2/landscapes/land3/sapsystems/sys3/", Value: []byte("")},
	}

	kv := &kv{list: func(prefix string, q *consulApi.QueryOptions) (consulApi.KVPairs, *consulApi.QueryMeta, error) {
		return kvPairs, nil, nil
	}}
	result, _ := kv.Maps("/trento", "/trento/")

	expectedMap := map[string]interface{}{
		"env1": map[string]interface{}{
			"landscapes": map[string]interface{}{
				"land1": map[string]interface{}{
					"sapsystems": map[string]interface{}{
						"sys1": map[string]interface{}{
							"name":        "name",
							"description": "desc",
						},
					},
				},
				"land2": map[string]interface{}{
					"description": "land_desc",
					"sapsystems": map[string]interface{}{
						"sys2": map[string]interface{}{},
					},
				},
			},
		},
		"env2": map[string]interface{}{
			"landscapes": map[string]interface{}{
				"land3": map[string]interface{}{
					"sapsystems": map[string]interface{}{
						"sys3": map[string]interface{}{},
					},
				},
			},
		},
	}

	assert.Equal(t, expectedMap, result)
}
