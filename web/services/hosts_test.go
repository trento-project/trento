package services

import (
	"testing"

	"github.com/stretchr/testify/assert"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/trento-project/trento/internal/consul/mocks"
	"github.com/trento-project/trento/internal/hosts"
)

func TestGetHostMetadata(t *testing.T) {
	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)

	nodes := []*consulApi.Node{
		{
			Node: "node1",
			Meta: map[string]string{
				"trento-meta1": "value1",
				"trento-meta2": "value2",
			},
		},
		{
			Node: "node2",
			Meta: map[string]string{
				"trento-meta3": "value3",
				"trento-meta4": "value4",
			},
		},
	}

	consulInst.On("Catalog").Return(catalog)
	catalog.On("Nodes", &consulApi.QueryOptions{Filter: "Node == node1"}).Return(nodes, nil, nil)

	hostsService := NewHostsService(consulInst)
	meta, err := hostsService.GetHostMetadata("node1")

	expectedMeta := map[string]string{
		"trento-meta1": "value1",
		"trento-meta2": "value2",
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedMeta, meta)
}

func TestGetHostsBySid(t *testing.T) {
	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)

	nodes := []*consulApi.Node{
		{
			Node: "node1",
			Meta: map[string]string{
				"trento-sap-systems": "PRD",
			},
		},
		{
			Node: "node2",
			Meta: map[string]string{
				"trento-sap-systems": "PRD",
			},
		},
	}

	consulInst.On("Catalog").Return(catalog)
	catalog.On("Nodes", &consulApi.QueryOptions{Filter: "Meta[\"trento-sap-systems\"] == PRD"}).Return(nodes, nil, nil)

	hostsService := NewHostsService(consulInst)
	hostsBySid, err := hostsService.GetHostsBySid("PRD")

	host1 := hosts.NewHost(
		consulApi.Node{
			Node: "node1",
			Meta: map[string]string{
				"trento-sap-systems": "PRD",
			},
		},
		consulInst,
	)

	host2 := hosts.NewHost(
		consulApi.Node{
			Node: "node2",
			Meta: map[string]string{
				"trento-sap-systems": "PRD",
			},
		},
		consulInst,
	)

	expectedHosts := hosts.HostList{
		&host1,
		&host2,
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedHosts, hostsBySid)
}
