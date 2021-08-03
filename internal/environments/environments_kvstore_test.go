package environments

import (
	"fmt"
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/consul/mocks"
	"github.com/trento-project/trento/internal/hosts"
)

func TestEnvironmentStore(t *testing.T) {
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)

	consulInst.On("KV").Return(kv)

	e := Environment{
		Name: "env1",
		Landscapes: map[string]*Landscape{
			"land1": &Landscape{
				Name: "land1",
				SAPSystems: map[string]*SAPSystem{
					"sys1": &SAPSystem{Name: "sys1", Type: 1},
				},
			},
		},
	}

	expectedPutMap := map[string]interface{}{
		"name": "env1",
		"landscapes": map[string]*Landscape{
			"land1": &Landscape{
				Name: "land1",
				SAPSystems: map[string]*SAPSystem{
					"sys1": &SAPSystem{
						Name: "sys1",
						Type: 1,
					},
				},
			},
		},
	}

	kvPath := fmt.Sprintf("%s%s", consul.KvEnvironmentsPath, "env1")
	kv.On("PutMap", kvPath, expectedPutMap).Return(nil)

	result := e.Store(consulInst)

	assert.Equal(t, result, nil)
}

func TestEnvironmentLoad(t *testing.T) {
	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	kv := new(mocks.KV)

	consulInst.On("Catalog").Return(catalog)
	consulInst.On("KV").Return(kv)

	nodes := []*consulApi.Node{
		{
			Node: "node1",
			Meta: map[string]string{
				"trento-sap-environment": "env1",
				"trento-sap-landscape":   "land1",
				"trento-sap-systems":     "sys1,sys3",
			},
		},
		{
			Node: "node2",
			Meta: map[string]string{
				"trento-sap-environment": "env2",
				"trento-sap-landscape":   "land2",
				"trento-sap-systems":     "sys2,sys4",
			},
		},
	}

	filter := &consulApi.QueryOptions{
		Filter: "(Meta[\"trento-sap-environment\"] == \"env1\") and (Meta[\"trento-sap-landscape\"] == \"land1\") and (Meta[\"trento-sap-systems\"] contains \"sys1\")"}
	catalog.On("Nodes", (filter)).Return(nodes, nil, nil)

	returnedMap := map[string]interface{}{
		"env1": map[string]interface{}{
			"name": "env1",
			"landscapes": map[string]interface{}{
				"land1": map[string]interface{}{
					"name": "land1",
					"sapsystems": map[string]interface{}{
						"sys1": map[string]interface{}{
							"name": "sys1",
							"type": 1,
						},
					},
				},
			},
		},
	}

	kv.On("ListMap", consul.KvEnvironmentsPath, consul.KvEnvironmentsPath).Return(returnedMap, nil)

	host1 := hosts.NewHost(*nodes[0], consulInst)
	host2 := hosts.NewHost(*nodes[1], consulInst)

	expectedEnv := map[string]*Environment{
		"env1": &Environment{
			Name: "env1",
			Landscapes: map[string]*Landscape{
				"land1": &Landscape{
					Name: "land1",
					SAPSystems: map[string]*SAPSystem{
						"sys1": &SAPSystem{
							Name: "sys1",
							Type: 1,
							Hosts: hosts.HostList{
								&host1,
								&host2,
							},
						},
					},
				},
			},
		},
	}

	e, _ := Load(consulInst)

	assert.Equal(t, e, expectedEnv)
}
