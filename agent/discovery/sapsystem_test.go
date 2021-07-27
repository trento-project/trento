package discovery

import (
	"fmt"
	"os"
	"path"
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trento-project/trento/internal/consul"
	consulMocks "github.com/trento-project/trento/internal/consul/mocks"
	"github.com/trento-project/trento/internal/environments"
	"github.com/trento-project/trento/internal/sapsystem"
	"github.com/trento-project/trento/internal/sapsystem/sapcontrol"
	sapcontrolMocks "github.com/trento-project/trento/internal/sapsystem/sapcontrol/mocks"
)

func TestStoreSAPSystemTags(t *testing.T) {
	kv := new(consulMocks.KV)
	catalog := new(consulMocks.Catalog)
	client := new(consulMocks.Client)
	client.On("Catalog").Return(catalog)
	client.On("KV").Return(kv)
	host, _ := os.Hostname()

	mockWebService := new(sapcontrolMocks.WebService)

	mockWebService.On("GetInstanceProperties").Return(&sapcontrol.GetInstancePropertiesResponse{
		Properties: []*sapcontrol.InstanceProperty{
			{
				Property:     "SAPSYSTEMNAME",
				Propertytype: "string",
				Value:        "sys1",
			},
		},
	}, nil)

	mockWebService.On("GetProcessList").Return(&sapcontrol.GetProcessListResponse{}, nil)
	mockWebService.On("GetSystemInstanceList").Return(&sapcontrol.GetSystemInstanceListResponse{}, nil)

	env := map[string]interface{}{
		"env1": map[string]interface{}{
			"name": "env1",
			"landscapes": map[string]interface{}{
				"land1": map[string]interface{}{
					"name": "land1",
					"sapsystems": map[string]interface{}{
						"DEV": map[string]interface{}{
							"name": "DEV",
							"type": sapsystem.Database,
						},
					},
				},
			},
		},
	}

	expectedHostMetadata := map[string]interface{}{
		"sap-environment": "env1",
		"sap-landscape":   "land1",
		"sap-systems":     "DEV",
	}

	kv.On("ListMap", consul.KvEnvironmentsPath, consul.KvEnvironmentsPath).Return(env, nil)
	catalog.On("Nodes", mock.Anything).Return([]*consulApi.Node{}, nil, nil)
	kv.On("PutMap", fmt.Sprintf(consul.KvHostsMetadataPath, host), expectedHostMetadata).Return(nil, nil)

	sapSystem := &sapsystem.SAPSystem{SID: "DEV", Type: sapsystem.Database}

	err := storeSAPSystemTags(client, []*sapsystem.SAPSystem{sapSystem})
	assert.NoError(t, err)

	kv.AssertExpectations(t)
}

func TestStoreSAPSystemTagsWithNoEnvs(t *testing.T) {
	kv := new(consulMocks.KV)
	catalog := new(consulMocks.Catalog)
	client := new(consulMocks.Client)
	client.On("Catalog").Return(catalog)
	client.On("KV").Return(kv)
	host, _ := os.Hostname()

	mockWebService := new(sapcontrolMocks.WebService)

	mockWebService.On("GetInstanceProperties").Return(&sapcontrol.GetInstancePropertiesResponse{
		Properties: []*sapcontrol.InstanceProperty{
			{
				Property:     "SAPSYSTEMNAME",
				Propertytype: "string",
				Value:        "sys1",
			},
		},
	}, nil)

	mockWebService.On("GetProcessList").Return(&sapcontrol.GetProcessListResponse{}, nil)
	mockWebService.On("GetSystemInstanceList").Return(&sapcontrol.GetSystemInstanceListResponse{}, nil)

	kv.On("ListMap", consul.KvEnvironmentsPath, consul.KvEnvironmentsPath).Return(nil, nil)
	catalog.On("Nodes", mock.Anything).Return([]*consulApi.Node{}, nil, nil)

	expectedNewEnv := map[string]interface{}{
		"name": "default",
		"landscapes": map[string]*environments.Landscape{
			"default": {
				Name: "default",
				SAPSystems: map[string]*environments.SAPSystem{
					"DEV": {
						Name: "DEV",
						Type: sapsystem.Database,
					},
				},
			},
		},
	}
	kv.On("PutMap", path.Join(consul.KvEnvironmentsPath, "default"), expectedNewEnv).Return(nil, nil)

	expectedHostMetadata := map[string]interface{}{
		"sap-environment": "default",
		"sap-landscape":   "default",
		"sap-systems":     "DEV",
	}
	kv.On("PutMap", fmt.Sprintf(consul.KvHostsMetadataPath, host), expectedHostMetadata).Return(nil, nil)

	sapSystem := &sapsystem.SAPSystem{SID: "DEV", Type: sapsystem.Database}

	err := storeSAPSystemTags(client, []*sapsystem.SAPSystem{sapSystem})
	assert.NoError(t, err)

	kv.AssertExpectations(t)
}
