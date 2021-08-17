package hosts

import (
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/internal/consul/mocks"
)

func TestCreateFilterMetaQuery(t *testing.T) {

	queryMap := map[string][]string{
		"trento-meta1": {"val1", "val2"},
		"trento-meta2": {"val3"},
		"trento-meta3": {"val4", "val5", "val6"},
		"meta4":        {"val7"},
	}

	query := CreateFilterMetaQuery(queryMap)

	expectedQuery := "(Meta[\"trento-meta1\"] == \"val1\" or Meta[\"trento-meta1\"] == \"val2\") and (Meta[\"trento-meta2\"] == \"val3\") and (Meta[\"trento-meta3\"] == \"val4\" or Meta[\"trento-meta3\"] == \"val5\" or Meta[\"trento-meta3\"] == \"val6\")"

	assert.Equal(t, expectedQuery, query)
}

func TestLoad(t *testing.T) {
	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	health := new(mocks.Health)

	consulInst.On("Catalog").Return(catalog)

	node1HealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthPassing,
		},
	}

	node2HealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthWarning,
		},
	}

	nodes := []*consulApi.Node{
		{
			Node: "node1",
		},
		{
			Node: "node2",
		},
	}

	consulInst.On("Health").Return(health)
	catalog.On("Nodes", &consulApi.QueryOptions{Filter: "query"}).Return(nodes, nil, nil)
	health.On("Node", "node1", (*consulApi.QueryOptions)(nil)).Return(node1HealthChecks, nil, nil)
	health.On("Node", "node2", (*consulApi.QueryOptions)(nil)).Return(node2HealthChecks, nil, nil)

	host1 := NewHost(
		consulApi.Node{
			Node: "node1",
		},
		consulInst,
	)

	expectedHosts := HostList{
		&host1,
	}

	h, _ := Load(consulInst, "query", []string{"passing"}, nil)

	assert.Equal(t, expectedHosts, h)
}
