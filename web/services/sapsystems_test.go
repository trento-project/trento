package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/consul/mocks"
	"github.com/trento-project/trento/internal/sapsystem"
)

func mockTests(consulInst *mocks.Client, catalog *mocks.Catalog, kv *mocks.KV) {
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
	catalog.On("Nodes", &consulApi.QueryOptions{Filter: ""}).Return(nodes, nil, nil)

	kvPath1 := fmt.Sprintf(consul.KvHostsSAPSystemPath, "node1")
	listMap1 := map[string]interface{}{
		"PRD": map[string]interface{}{
			"id":   "systemId1",
			"sid":  "PRD",
			"type": sapsystem.Database,
			"databases": []interface{}{
				map[string]interface{}{
					"database": "PRD",
					"host":     "node1",
				},
			},
		},
		"DEV": map[string]interface{}{
			"id":   "systemId2",
			"sid":  "DEV",
			"type": sapsystem.Database,
		},
	}

	kv.On("ListMap", kvPath1, kvPath1).Return(listMap1, nil)
	consulInst.On("WaitLock", fmt.Sprintf(consul.KvHostsSAPSystemPath, "node1")).Return(nil)

	kvPath2 := fmt.Sprintf(consul.KvHostsSAPSystemPath, "node2")
	listMap2 := map[string]interface{}{
		"QAS": map[string]interface{}{
			"id":   "systemId3",
			"sid":  "QAS",
			"type": sapsystem.Application,
			"profile": map[string]interface{}{
				"dbs": map[string]interface{}{
					"hdb": map[string]interface{}{
						"dbname": "QAS",
					},
				},
				"SAPDBHOST": "192.168.1.6",
			},
		},
		"HA1": map[string]interface{}{
			"id":   "systemId4",
			"sid":  "HA1",
			"type": sapsystem.Application,
			"profile": map[string]interface{}{
				"dbs": map[string]interface{}{
					"hdb": map[string]interface{}{
						"dbname": "PRD",
					},
				},
				"SAPDBHOST": "192.168.1.5",
			},
		},
	}

	kv.On("ListMap", kvPath2, kvPath2).Return(listMap2, nil)
	consulInst.On("WaitLock", fmt.Sprintf(consul.KvHostsSAPSystemPath, "node2")).Return(nil)

	consulInst.On("KV").Return(kv)
}

func TestGetSAPSystems(t *testing.T) {
	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	kv := new(mocks.KV)

	mockTests(consulInst, catalog, kv)

	sapSystemsService := NewSAPSystemsService(consulInst)
	sapSystems, err := sapSystemsService.GetSAPSystems()

	expectedSystems := sapsystem.SAPSystemsList{
		&sapsystem.SAPSystem{
			Id:   "systemId1",
			SID:  "PRD",
			Type: sapsystem.Database,
			Databases: []*sapsystem.DatabaseData{
				&sapsystem.DatabaseData{
					Database: "PRD",
					Host:     "node1",
				},
			},
		},
		&sapsystem.SAPSystem{
			Id:   "systemId2",
			SID:  "DEV",
			Type: sapsystem.Database,
		},
		&sapsystem.SAPSystem{
			Id:   "systemId3",
			SID:  "QAS",
			Type: sapsystem.Application,
			Profile: map[string]interface{}{
				"dbs": map[string]interface{}{
					"hdb": map[string]interface{}{
						"dbname": "QAS",
					},
				},
				"SAPDBHOST": "192.168.1.6",
			},
		},
		&sapsystem.SAPSystem{
			Id:   "systemId4",
			SID:  "HA1",
			Type: sapsystem.Application,
			Profile: map[string]interface{}{
				"dbs": map[string]interface{}{
					"hdb": map[string]interface{}{
						"dbname": "PRD",
					},
				},
				"SAPDBHOST": "192.168.1.5",
			},
		},
	}

	catalog.AssertExpectations(t)
	consulInst.AssertExpectations(t)
	kv.AssertExpectations(t)

	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedSystems, sapSystems)
}

func TestGetSAPSystemsById(t *testing.T) {
	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	kv := new(mocks.KV)

	mockTests(consulInst, catalog, kv)

	sapSystemsService := NewSAPSystemsService(consulInst)
	sapSystems, err := sapSystemsService.GetSAPSystemsById("systemId1")

	expectedSystems := sapsystem.SAPSystemsList{
		&sapsystem.SAPSystem{
			Id:   "systemId1",
			SID:  "PRD",
			Type: sapsystem.Database,
			Databases: []*sapsystem.DatabaseData{
				&sapsystem.DatabaseData{
					Database: "PRD",
					Host:     "node1",
				},
			},
		},
	}

	catalog.AssertExpectations(t)
	consulInst.AssertExpectations(t)
	kv.AssertExpectations(t)

	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedSystems, sapSystems)
}

func TestGetSAPSystemsByType(t *testing.T) {
	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	kv := new(mocks.KV)

	mockTests(consulInst, catalog, kv)

	sapSystemsService := NewSAPSystemsService(consulInst)
	sapSystems, err := sapSystemsService.GetSAPSystemsByType(sapsystem.Database)

	expectedSystems := sapsystem.SAPSystemsList{
		&sapsystem.SAPSystem{
			Id:   "systemId1",
			SID:  "PRD",
			Type: sapsystem.Database,
			Databases: []*sapsystem.DatabaseData{
				&sapsystem.DatabaseData{
					Database: "PRD",
					Host:     "node1",
				},
			},
		},
		&sapsystem.SAPSystem{
			Id:   "systemId2",
			SID:  "DEV",
			Type: sapsystem.Database,
		},
	}

	catalog.AssertExpectations(t)
	consulInst.AssertExpectations(t)
	kv.AssertExpectations(t)

	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedSystems, sapSystems)
}

func TestGetAttachedDatabasesByIdSystemNotFound(t *testing.T) {
	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	kv := new(mocks.KV)

	mockTests(consulInst, catalog, kv)

	sapSystemsService := NewSAPSystemsService(consulInst)
	_, err := sapSystemsService.GetAttachedDatabasesById("systemNotFound")

	catalog.AssertExpectations(t)
	consulInst.AssertExpectations(t)
	kv.AssertExpectations(t)

	assert.Error(t, err, fmt.Errorf("system with systemNotFound not found"))
}

func TestGetAttachedDatabasesByIdIpNotFound(t *testing.T) {
	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	kv := new(mocks.KV)

	mockTests(consulInst, catalog, kv)

	catalog.On("Nodes",
		&consulApi.QueryOptions{Filter: "Meta[\"trento-host-ip-addresses\"] contains \"192.168.1.5\""}).Return(
		nil, nil, nil)

	sapSystemsService := NewSAPSystemsService(consulInst)
	sapAttachedDbs, err := sapSystemsService.GetAttachedDatabasesById("systemId4")

	catalog.AssertExpectations(t)
	consulInst.AssertExpectations(t)
	kv.AssertExpectations(t)

	assert.NoError(t, err)
	assert.ElementsMatch(t, sapsystem.SAPSystemsList{}, sapAttachedDbs)
}

func TestGetAttachedDatabasesByIdDbNotFound(t *testing.T) {
	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	kv := new(mocks.KV)

	mockTests(consulInst, catalog, kv)

	nodes := []*consulApi.Node{
		{
			Node: "node1",
		},
	}

	catalog.On("Nodes",
		&consulApi.QueryOptions{Filter: "Meta[\"trento-host-ip-addresses\"] contains \"192.168.1.6\""}).Return(
		nodes, nil, nil)

	sapSystemsService := NewSAPSystemsService(consulInst)
	sapAttachedDbs, err := sapSystemsService.GetAttachedDatabasesById("systemId3")

	catalog.AssertExpectations(t)
	consulInst.AssertExpectations(t)
	kv.AssertExpectations(t)

	assert.NoError(t, err)
	assert.ElementsMatch(t, sapsystem.SAPSystemsList{}, sapAttachedDbs)
}

func TestGetAttachedDatabasesById(t *testing.T) {
	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	kv := new(mocks.KV)

	mockTests(consulInst, catalog, kv)

	nodes := []*consulApi.Node{
		{
			Node: "node1",
		},
	}

	catalog.On("Nodes",
		&consulApi.QueryOptions{Filter: "Meta[\"trento-host-ip-addresses\"] contains \"192.168.1.5\""}).Return(
		nodes, nil, nil)

	sapSystemsService := NewSAPSystemsService(consulInst)
	sapAttachedDbs, err := sapSystemsService.GetAttachedDatabasesById("systemId4")

	expectedDbs := sapsystem.SAPSystemsList{
		&sapsystem.SAPSystem{
			Id:   "systemId1",
			SID:  "PRD",
			Type: sapsystem.Database,
			Databases: []*sapsystem.DatabaseData{
				&sapsystem.DatabaseData{
					Database: "PRD",
					Host:     "node1",
				},
			},
		},
	}

	catalog.AssertExpectations(t)
	consulInst.AssertExpectations(t)
	kv.AssertExpectations(t)

	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedDbs, sapAttachedDbs)
}
