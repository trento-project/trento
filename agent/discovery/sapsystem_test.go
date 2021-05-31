package discovery

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/SUSE/sap_host_exporter/lib/sapcontrol"
	"github.com/SUSE/sap_host_exporter/test/mock_sapcontrol"
	"github.com/golang/mock/gomock"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/consul/mocks"
	"github.com/trento-project/trento/internal/environments"
	"github.com/trento-project/trento/internal/sapsystem"
)

func TestStoreSAPSystemTags(t *testing.T) {
	kv := new(mocks.KV)
	catalog := new(mocks.Catalog)
	client := new(mocks.Client)
	client.On("Catalog").Return(catalog)
	client.On("KV").Return(kv)
	host, _ := os.Hostname()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWebService := mock_sapcontrol.NewMockWebService(ctrl)

	mockWebService.EXPECT().GetInstanceProperties().Return(&sapcontrol.GetInstancePropertiesResponse{
		Properties: []*sapcontrol.InstanceProperty{
			{
				Property:     "SAPSYSTEMNAME",
				Propertytype: "string",
				Value:        "sys1",
			},
		},
	}, nil)

	mockWebService.EXPECT().GetProcessList().Return(&sapcontrol.GetProcessListResponse{}, nil)
	mockWebService.EXPECT().GetSystemInstanceList().Return(&sapcontrol.GetSystemInstanceListResponse{}, nil)

	environment := map[string]interface{}{
		"env1": map[string]interface{}{
			"name": "env1",
			"landscapes": map[string]interface{}{
				"land1": map[string]interface{}{
					"name": "land1",
					"sapsystems": map[string]interface{}{
						"sys1": map[string]interface{}{
							"name": "sys1",
							"type": "type1",
						},
					},
				},
			},
		},
	}
	expectedHostMetadata := map[string]interface{}{
		"sap-environment": "env1",
		"sap-landscape":   "land1",
		"sap-system":      "sys1",
	}

	kv.On("ListMap", consul.KvEnvironmentsPath, consul.KvEnvironmentsPath).Return(environment, nil)
	catalog.On("Nodes", mock.Anything).Return([]*consulApi.Node{}, nil, nil)
	kv.On("PutMap", fmt.Sprintf(consul.KvHostsMetadataPath, host), expectedHostMetadata).Return(nil, nil)

	var err error

	sapSystem, err := sapsystem.NewSAPSystem(mockWebService)
	assert.NoError(t, err)

	err = storeSAPSystemTags(client, sapSystem)
	assert.NoError(t, err)

	kv.AssertExpectations(t)
}

func TestStoreSAPSystemTagsWithNoEnvs(t *testing.T) {
	kv := new(mocks.KV)
	catalog := new(mocks.Catalog)
	client := new(mocks.Client)
	client.On("Catalog").Return(catalog)
	client.On("KV").Return(kv)
	host, _ := os.Hostname()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWebService := mock_sapcontrol.NewMockWebService(ctrl)

	mockWebService.EXPECT().GetInstanceProperties().Return(&sapcontrol.GetInstancePropertiesResponse{
		Properties: []*sapcontrol.InstanceProperty{
			{
				Property:     "SAPSYSTEMNAME",
				Propertytype: "string",
				Value:        "sys1",
			},
		},
	}, nil)

	mockWebService.EXPECT().GetProcessList().Return(&sapcontrol.GetProcessListResponse{}, nil)
	mockWebService.EXPECT().GetSystemInstanceList().Return(&sapcontrol.GetSystemInstanceListResponse{}, nil)

	kv.On("ListMap", consul.KvEnvironmentsPath, consul.KvEnvironmentsPath).Return(nil, nil)
	catalog.On("Nodes", mock.Anything).Return([]*consulApi.Node{}, nil, nil)

	expectedNewEnv := map[string]interface{}{
		"name": "default",
		"landscapes": map[string]*environments.Landscape{
			"default": {
				Name: "default",
				SAPSystems: map[string]*environments.SAPSystem{
					"sys1": {
						Name: "sys1",
						Type: "APP",
					},
				},
			},
		},
	}
	kv.On("PutMap", path.Join(consul.KvEnvironmentsPath, "default"), expectedNewEnv).Return(nil, nil)

	expectedHostMetadata := map[string]interface{}{
		"sap-environment": "default",
		"sap-landscape":   "default",
		"sap-system":      "sys1",
	}
	kv.On("PutMap", fmt.Sprintf(consul.KvHostsMetadataPath, host), expectedHostMetadata).Return(nil, nil)

	var err error

	sapSystem, err := sapsystem.NewSAPSystem(mockWebService)
	assert.NoError(t, err)

	err = storeSAPSystemTags(client, sapSystem)
	assert.NoError(t, err)

	kv.AssertExpectations(t)
}
