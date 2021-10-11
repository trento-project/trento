package discovery

import (
	"fmt"
	"os"
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trento-project/trento/internal/consul"
	consulMocks "github.com/trento-project/trento/internal/consul/mocks"
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

	expectedHostMetadata := map[string]interface{}{
		"sap-systems":      "DEV",
		"sap-systems-type": "Database",
		"sap-systems-id":   "systemId",
	}

	catalog.On("Nodes", mock.Anything).Return([]*consulApi.Node{}, nil, nil)
	kv.On("PutMap", fmt.Sprintf(consul.KvHostsMetadataPath, host), expectedHostMetadata).Return(nil, nil)

	sapSystem := &sapsystem.SAPSystem{Id: "systemId", SID: "DEV", Type: sapsystem.Database}

	err := storeSAPSystemTags(client, []*sapsystem.SAPSystem{sapSystem})
	assert.NoError(t, err)

	kv.AssertExpectations(t)
}
