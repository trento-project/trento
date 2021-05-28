package sapsystem

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
)

func TestStore(t *testing.T) {
	host, _ := os.Hostname()
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)

	consulInst.On("KV").Return(kv)
	kvPath := fmt.Sprintf("%s/%s", fmt.Sprintf(consul.KvHostsSAPSystemPath, host), "PRD")

	expectedPutMap := map[string]interface{}{
		"type": "HANA",
		"processes": map[string]*sapcontrol.OSProcess{
			"enserver": {
				Name:        "enserver",
				Description: "foobar",
				Dispstatus:  sapcontrol.STATECOLOR_GREEN,
				Textstatus:  "Running",
				Starttime:   "1",
				Elapsedtime: "2",
				Pid:         30787,
			},
			"msg_server": {
				Name:        "msg_server",
				Description: "foobar2",
				Dispstatus:  sapcontrol.STATECOLOR_YELLOW,
				Textstatus:  "Stopping",
				Starttime:   "43",
				Elapsedtime: "",
				Pid:         30786,
			},
		},
		"properties": map[string]*sapcontrol.InstanceProperty{
			"INSTANCE_NAME": {
				Property:     "INSTANCE_NAME",
				Propertytype: "string",
				Value:        "HDB00",
			},
			"SAPSYSTEMNAME": {
				Property:     "SAPSYSTEMNAME",
				Propertytype: "string",
				Value:        "PRD",
			},
			"HANA Roles": {
				Property:     "HANA Roles",
				Propertytype: "type3",
				Value:        "some hana value",
			},
		},
		"instances": map[string]*sapcontrol.SAPInstance{
			"host1": {
				Hostname:      "host1",
				InstanceNr:    0,
				HttpPort:      50013,
				HttpsPort:     50014,
				StartPriority: "0.3",
				Features:      "some features",
				Dispstatus:    sapcontrol.STATECOLOR_GREEN,
			},
			"host2": {
				Hostname:      "host2",
				InstanceNr:    1,
				HttpPort:      50113,
				HttpsPort:     50114,
				StartPriority: "0.3",
				Features:      "some other features",
				Dispstatus:    sapcontrol.STATECOLOR_YELLOW,
			},
		},
	}

	kv.On("DeleteTree", kvPath, (*consulApi.WriteOptions)(nil)).Return(nil, nil)
	kv.On("PutMap", kvPath, expectedPutMap).Return(nil, nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWebService := mock_sapcontrol.NewMockWebService(ctrl)

	s := SAPSystem{
		webService: mockWebService,
		Id:         "",
		Type:       "HANA",
		Processes: map[string]*sapcontrol.OSProcess{
			"enserver": {
				Name:        "enserver",
				Description: "foobar",
				Dispstatus:  sapcontrol.STATECOLOR_GREEN,
				Textstatus:  "Running",
				Starttime:   "1",
				Elapsedtime: "2",
				Pid:         30787,
			},
			"msg_server": {
				Name:        "msg_server",
				Description: "foobar2",
				Dispstatus:  sapcontrol.STATECOLOR_YELLOW,
				Textstatus:  "Stopping",
				Starttime:   "43",
				Elapsedtime: "",
				Pid:         30786,
			},
		},
		Properties: map[string]*sapcontrol.InstanceProperty{
			"INSTANCE_NAME": {
				Property:     "INSTANCE_NAME",
				Propertytype: "string",
				Value:        "HDB00",
			},
			"SAPSYSTEMNAME": {
				Property:     "SAPSYSTEMNAME",
				Propertytype: "string",
				Value:        "PRD",
			},
			"HANA Roles": {
				Property:     "HANA Roles",
				Propertytype: "type3",
				Value:        "some hana value",
			},
		},
		Instances: map[string]*sapcontrol.SAPInstance{
			"host1": {
				Hostname:      "host1",
				InstanceNr:    0,
				HttpPort:      50013,
				HttpsPort:     50014,
				StartPriority: "0.3",
				Features:      "some features",
				Dispstatus:    sapcontrol.STATECOLOR_GREEN,
			},
			"host2": {
				Hostname:      "host2",
				InstanceNr:    1,
				HttpPort:      50113,
				HttpsPort:     50114,
				StartPriority: "0.3",
				Features:      "some other features",
				Dispstatus:    sapcontrol.STATECOLOR_YELLOW,
			},
		},
	}

	err := s.Store(consulInst)

	assert.NoError(t, err)

	kv.AssertExpectations(t)
}

func TestLoad(t *testing.T) {
	host, _ := os.Hostname()
	kvPath := fmt.Sprintf(consul.KvHostsSAPSystemPath, host)
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)

	listMap := map[string]interface{}{
		"PRD": map[string]interface{}{
			"type": "HANA",
			"processes": map[string]interface{}{
				"enserver": map[string]interface{}{
					"Name":        "enserver",
					"Description": "foobar",
					"Dispstatus":  sapcontrol.STATECOLOR_GREEN,
					"Textstatus":  "Running",
					"Starttime":   "1",
					"Elapsedtime": "2",
					"Pid":         30787,
				},
				"msg_server": map[string]interface{}{
					"Name":        "msg_server",
					"Description": "foobar2",
					"Dispstatus":  sapcontrol.STATECOLOR_YELLOW,
					"Textstatus":  "Stopping",
					"Starttime":   "43",
					"Elapsedtime": "",
					"Pid":         30786,
				},
			},
			"properties": map[string]interface{}{
				"INSTANCE_NAME": map[string]interface{}{
					"Property":     "INSTANCE_NAME",
					"Propertytype": "string",
					"Value":        "HDB00",
				},
				"SAPSYSTEMNAME": map[string]interface{}{
					"Property":     "SAPSYSTEMNAME",
					"Propertytype": "string",
					"Value":        "PRD",
				},
				"HANA Roles": map[string]interface{}{
					"Property":     "HANA Roles",
					"Propertytype": "type3",
					"Value":        "some hana value",
				},
			},
			"instances": map[string]interface{}{
				"host1": map[string]interface{}{
					"Hostname":      "host1",
					"InstanceNr":    0,
					"HttpPort":      50013,
					"HttpsPort":     50014,
					"StartPriority": "0.3",
					"Features":      "some features",
					"Dispstatus":    sapcontrol.STATECOLOR_GREEN,
				},
				"host2": map[string]interface{}{
					"Hostname":      "host2",
					"InstanceNr":    1,
					"HttpPort":      50113,
					"HttpsPort":     50114,
					"StartPriority": "0.3",
					"Features":      "some other features",
					"Dispstatus":    sapcontrol.STATECOLOR_YELLOW,
				},
			},
		},
	}

	kv.On("ListMap", kvPath, kvPath).Return(listMap, nil)

	consulInst.On("KV").Return(kv)

	systems, _ := Load(consulInst, host)

	expectedSystems := map[string]*SAPSystem{
		"PRD": {
			Id:   "",
			Type: "HANA",
			Processes: map[string]*sapcontrol.OSProcess{
				"enserver": {
					Name:        "enserver",
					Description: "foobar",
					Dispstatus:  sapcontrol.STATECOLOR_GREEN,
					Textstatus:  "Running",
					Starttime:   "1",
					Elapsedtime: "2",
					Pid:         30787,
				},
				"msg_server": {
					Name:        "msg_server",
					Description: "foobar2",
					Dispstatus:  sapcontrol.STATECOLOR_YELLOW,
					Textstatus:  "Stopping",
					Starttime:   "43",
					Elapsedtime: "",
					Pid:         30786,
				},
			},
			Properties: map[string]*sapcontrol.InstanceProperty{
				"INSTANCE_NAME": {
					Property:     "INSTANCE_NAME",
					Propertytype: "string",
					Value:        "HDB00",
				},
				"SAPSYSTEMNAME": {
					Property:     "SAPSYSTEMNAME",
					Propertytype: "string",
					Value:        "PRD",
				},
				"HANA Roles": {
					Property:     "HANA Roles",
					Propertytype: "type3",
					Value:        "some hana value",
				},
			},
			Instances: map[string]*sapcontrol.SAPInstance{
				"host1": {
					Hostname:      "host1",
					InstanceNr:    0,
					HttpPort:      50013,
					HttpsPort:     50014,
					StartPriority: "0.3",
					Features:      "some features",
					Dispstatus:    sapcontrol.STATECOLOR_GREEN,
				},
				"host2": {
					Hostname:      "host2",
					InstanceNr:    1,
					HttpPort:      50113,
					HttpsPort:     50114,
					StartPriority: "0.3",
					Features:      "some other features",
					Dispstatus:    sapcontrol.STATECOLOR_YELLOW,
				},
			},
		},
	}

	assert.Equal(t, expectedSystems, systems)
}

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

	sapSystem, err := NewSAPSystem(mockWebService)
	assert.NoError(t, err)

	err = sapSystem.StoreSAPSystemTags(client)
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

	sapSystem, err := NewSAPSystem(mockWebService)
	assert.NoError(t, err)

	err = sapSystem.StoreSAPSystemTags(client)
	assert.NoError(t, err)

	kv.AssertExpectations(t)
}
