package consul

import (
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

func TestKVListMap(t *testing.T) {

	kvPairs := consulApi.KVPairs{
		&consulApi.KVPair{Key: "/trento/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/resource2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/resource2/sys1/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/resource2/sys1/name", Value: []byte("name")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/resource2/sys1/int32", Value: []byte("1"), Flags: int32Flag},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath2/resource2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath2/resource2/sys2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path2/resource1/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path2/resource1/subpath3/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path2/resource1/subpath3/resource2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path2/resource1/subpath3/resource2/sys3/", Value: []byte("")},
	}

	kv := &kv{list: func(prefix string, q *consulApi.QueryOptions) (consulApi.KVPairs, *consulApi.QueryMeta, error) {
		return kvPairs, nil, nil
	}}
	result, _ := kv.ListMap("/trento", "")

	expectedMap := map[string]interface{}{
		"trento": map[string]interface{}{
			"path1": map[string]interface{}{
				"resource1": map[string]interface{}{
					"subpath1": map[string]interface{}{
						"resource2": map[string]interface{}{
							"sys1": map[string]interface{}{
								"name":  "name",
								"int32": int32(1),
							},
						},
					},
					"subpath2": map[string]interface{}{
						"resource2": map[string]interface{}{
							"sys2": map[string]interface{}{},
						},
					},
				},
			},
			"path2": map[string]interface{}{
				"resource1": map[string]interface{}{
					"subpath3": map[string]interface{}{
						"resource2": map[string]interface{}{
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
		&consulApi.KVPair{Key: "/trento/path1/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/resource2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/resource2/sys1/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/resource2/sys1/name", Value: []byte("name")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath2/resource2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath2/resource2/sys2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path2/resource1/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path2/resource1/subpath3/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path2/resource1/subpath3/resource2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path2/resource1/subpath3/resource2/sys3/", Value: []byte("")},
	}

	kv := &kv{list: func(prefix string, q *consulApi.QueryOptions) (consulApi.KVPairs, *consulApi.QueryMeta, error) {
		return kvPairs, nil, nil
	}}
	result, _ := kv.ListMap("/trento", "/trento/")

	expectedMap := map[string]interface{}{
		"path1": map[string]interface{}{
			"resource1": map[string]interface{}{
				"subpath1": map[string]interface{}{
					"resource2": map[string]interface{}{
						"sys1": map[string]interface{}{
							"name": "name",
						},
					},
				},
				"subpath2": map[string]interface{}{
					"resource2": map[string]interface{}{
						"sys2": map[string]interface{}{},
					},
				},
			},
		},
		"path2": map[string]interface{}{
			"resource1": map[string]interface{}{
				"subpath3": map[string]interface{}{
					"resource2": map[string]interface{}{
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
		&consulApi.KVPair{Key: "/trento/path1/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/resource2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/resource2/sys1/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/resource2/sys1/name", Value: []byte("name")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath2/resource2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath2/resource2/sys2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path2/resource1/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path2/resource1/subpath3/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path2/resource1/subpath3/resource2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path2/resource1/subpath3/resource2/sys3/", Value: []byte("")},
	}

	kv := &kv{list: func(prefix string, q *consulApi.QueryOptions) (consulApi.KVPairs, *consulApi.QueryMeta, error) {
		return kvPairs, nil, nil
	}}
	result, _ := kv.ListMap("/trento", "/trento")

	expectedMap := map[string]interface{}{
		"path1": map[string]interface{}{
			"resource1": map[string]interface{}{
				"subpath1": map[string]interface{}{
					"resource2": map[string]interface{}{
						"sys1": map[string]interface{}{
							"name": "name",
						},
					},
				},
				"subpath2": map[string]interface{}{
					"resource2": map[string]interface{}{
						"sys2": map[string]interface{}{},
					},
				},
			},
		},
		"path2": map[string]interface{}{
			"resource1": map[string]interface{}{
				"subpath3": map[string]interface{}{
					"resource2": map[string]interface{}{
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
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/resource2/sys1/name", Value: []byte("name")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/resource2/sys1/description", Value: []byte("desc")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath2/description", Value: []byte("subpath_desc")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath2/resource2/sys2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path2/resource1/subpath3/resource2/sys3/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/list/", Value: []byte(""), Flags: sliceFlag},
		&consulApi.KVPair{Key: "/trento/list/0000/data1", Value: []byte("false"), Flags: boolFlag},
		&consulApi.KVPair{Key: "/trento/list/0000/data2", Value: []byte("3"), Flags: int32Flag},
		&consulApi.KVPair{Key: "/trento/list/0001/data3/dataa", Value: []byte("4"), Flags: intFlag},
		&consulApi.KVPair{Key: "/trento/list/0001/data3/datab", Value: []byte("5"), Flags: intFlag},
		&consulApi.KVPair{Key: "/trento/list/0002/otherlist/", Value: []byte(""), Flags: sliceFlag},
		&consulApi.KVPair{Key: "/trento/list/0002/otherlist/0000/datac", Value: []byte("c"), Flags: stringFlag},
		&consulApi.KVPair{Key: "/trento/list/0002/otherlist/0001/datad", Value: []byte("d"), Flags: stringFlag},
	}

	kv := &kv{list: func(prefix string, q *consulApi.QueryOptions) (consulApi.KVPairs, *consulApi.QueryMeta, error) {
		return kvPairs, nil, nil
	}}
	result, _ := kv.ListMap("/trento", "/trento/")

	expectedMap := map[string]interface{}{
		"path1": map[string]interface{}{
			"resource1": map[string]interface{}{
				"subpath1": map[string]interface{}{
					"resource2": map[string]interface{}{
						"sys1": map[string]interface{}{
							"name":        "name",
							"description": "desc",
						},
					},
				},
				"subpath2": map[string]interface{}{
					"description": "subpath_desc",
					"resource2": map[string]interface{}{
						"sys2": map[string]interface{}{},
					},
				},
			},
		},
		"path2": map[string]interface{}{
			"resource1": map[string]interface{}{
				"subpath3": map[string]interface{}{
					"resource2": map[string]interface{}{
						"sys3": map[string]interface{}{},
					},
				},
			},
		},
		"list": []interface{}{
			map[string]interface{}{
				"data1": false,
				"data2": int32(3),
			},
			map[string]interface{}{
				"data3": map[string]interface{}{
					"dataa": 4,
					"datab": 5,
				},
			},
			map[string]interface{}{
				"otherlist": []interface{}{
					map[string]interface{}{
						"datac": "c",
					},
					map[string]interface{}{
						"datad": "d",
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
		"path1": map[string]interface{}{
			"resource1": map[string]interface{}{
				"subpath1": map[string]interface{}{
					"resource2": map[string]interface{}{
						"sys1": map[string]interface{}{
							"name":        "name",
							"description": "desc",
							"int32":       int32(1),
						},
					},
				},
				"subpath2": map[string]interface{}{
					"description": "subpath_desc",
					"resource2": map[string]interface{}{
						"sys2": map[string]interface{}{},
					},
				},
			},
		},
		"path2": map[string]interface{}{
			"resource1": map[string]interface{}{
				"subpath3": map[string]interface{}{
					"resource2": map[string]interface{}{
						"sys3": map[string]interface{}{},
					},
				},
			},
		},
		"list": []interface{}{
			map[string]interface{}{
				"item1": int32(1),
				"other_list": []interface{}{
					map[string]interface{}{
						"itema": "a",
						"itemb": "b",
					},
					map[string]interface{}{
						"itemc": "c",
						"itemd": "d",
					},
				},
				"simple_list": []interface{}{
					"value1",
					"value2",
				},
			},
			map[string]interface{}{
				"item2": true,
			},
			map[string]interface{}{
				"item3": 3,
			},
			map[string]interface{}{
				"item4": 4,
			},
			map[string]interface{}{
				"item5": 5,
			},
			map[string]interface{}{
				"item6": 6,
			},
			map[string]interface{}{
				"item7": 7,
			},
			map[string]interface{}{
				"item8": 8,
			},
			map[string]interface{}{
				"item9": 9,
			},
			map[string]interface{}{
				"item10": 10,
			},
			map[string]interface{}{
				"item11": 11,
			},
		},
	}

	expectedPut := []*consulApi.KVPair{
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/resource2/sys1/name", Value: []byte("name")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/resource2/sys1/description", Value: []byte("desc")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath1/resource2/sys1/int32", Value: []byte("1"), Flags: int32Flag},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath2/description", Value: []byte("subpath_desc")},
		&consulApi.KVPair{Key: "/trento/path1/resource1/subpath2/resource2/sys2/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/path2/resource1/subpath3/resource2/sys3/", Value: []byte("")},
		&consulApi.KVPair{Key: "/trento/list/", Value: []byte(""), Flags: sliceFlag},
		&consulApi.KVPair{Key: "/trento/list/0000/item1", Value: []byte("1"), Flags: int32Flag},
		&consulApi.KVPair{Key: "/trento/list/0000/other_list/", Value: []byte(""), Flags: sliceFlag},
		&consulApi.KVPair{Key: "/trento/list/0000/other_list/0000/itema", Value: []byte("a"), Flags: stringFlag},
		&consulApi.KVPair{Key: "/trento/list/0000/other_list/0000/itemb", Value: []byte("b"), Flags: stringFlag},
		&consulApi.KVPair{Key: "/trento/list/0000/other_list/0001/itemc", Value: []byte("c"), Flags: stringFlag},
		&consulApi.KVPair{Key: "/trento/list/0000/other_list/0001/itemd", Value: []byte("d"), Flags: stringFlag},
		&consulApi.KVPair{Key: "/trento/list/0000/simple_list/", Value: []byte(""), Flags: sliceFlag},
		&consulApi.KVPair{Key: "/trento/list/0000/simple_list/0000", Value: []byte("value1"), Flags: stringFlag},
		&consulApi.KVPair{Key: "/trento/list/0000/simple_list/0001", Value: []byte("value2"), Flags: stringFlag},
		&consulApi.KVPair{Key: "/trento/list/0001/item2", Value: []byte("true"), Flags: boolFlag},
		&consulApi.KVPair{Key: "/trento/list/0002/item3", Value: []byte("3"), Flags: intFlag},
		&consulApi.KVPair{Key: "/trento/list/0003/item4", Value: []byte("4"), Flags: intFlag},
		&consulApi.KVPair{Key: "/trento/list/0004/item5", Value: []byte("5"), Flags: intFlag},
		&consulApi.KVPair{Key: "/trento/list/0005/item6", Value: []byte("6"), Flags: intFlag},
		&consulApi.KVPair{Key: "/trento/list/0006/item7", Value: []byte("7"), Flags: intFlag},
		&consulApi.KVPair{Key: "/trento/list/0007/item8", Value: []byte("8"), Flags: intFlag},
		&consulApi.KVPair{Key: "/trento/list/0008/item9", Value: []byte("9"), Flags: intFlag},
		&consulApi.KVPair{Key: "/trento/list/0009/item10", Value: []byte("10"), Flags: intFlag},
		&consulApi.KVPair{Key: "/trento/list/0010/item11", Value: []byte("11"), Flags: intFlag},
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
