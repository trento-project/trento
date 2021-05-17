package consul

import (
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

func TestKVListMap(t *testing.T) {

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
	result, _ := kv.ListMap("/trento", "")

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

func TestKVListMapOffset(t *testing.T) {

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
	result, _ := kv.ListMap("/trento", "/trento/")

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

func TestKVListMapOffsetNoBackslash(t *testing.T) {

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
	result, _ := kv.ListMap("/trento", "/trento")

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

func TestKVListMapSingleEntries(t *testing.T) {

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
	result, _ := kv.ListMap("/trento", "/trento/")

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

// Tests for PutMap

func TestKVPutMap(t *testing.T) {
	testMap := map[string]interface{}{
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

	expectedPut := []*consulApi.KVPair{
		&consulApi.KVPair{Key: "/trento/env1/landscapes/land1/sapsystems/sys1/name", Value: []byte("name")},
		&consulApi.KVPair{Key: "/trento/env1/landscapes/land1/sapsystems/sys1/description", Value: []byte("desc")},
		&consulApi.KVPair{Key: "/trento/env1/landscapes/land2/description", Value: []byte("land_desc")},
		//TODO: Fix the code to make these options work
		//&consulApi.KVPair{Key: "/trento/env1/landscapes/land2/sapsystems/sys2/", Value: []byte("")},
		//&consulApi.KVPair{Key: "/trento/env2/landscapes/land3/sapsystems/sys3/", Value: []byte("")},
	}

	resultPut := []*consulApi.KVPair{}

	mockedFunc := func(p *consulApi.KVPair, q *consulApi.WriteOptions) (*consulApi.WriteMeta, error) {
		resultPut = append(resultPut, p)
		return nil, nil
	}

	kv := &kv{put: mockedFunc}

	kv.PutMap("/trento", testMap)

	assert.ElementsMatch(t, expectedPut, resultPut)
}
