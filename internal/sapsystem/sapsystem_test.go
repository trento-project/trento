package sapsystem

import (
	"sort"
	"testing"

	"github.com/SUSE/sap_host_exporter/lib/sapcontrol"
	"github.com/SUSE/sap_host_exporter/test/mock_sapcontrol"
	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNewSAPSystem(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWebService := mock_sapcontrol.NewMockWebService(ctrl)

	mockWebService.EXPECT().GetInstanceProperties().Return(&sapcontrol.GetInstancePropertiesResponse{
		Properties: []*sapcontrol.InstanceProperty{
			{
				Property:     "prop1",
				Propertytype: "type1",
				Value:        "value1",
			},
			{
				Property:     "SAPSYSTEMNAME",
				Propertytype: "string",
				Value:        "PRD",
			},
			{
				Property:     "HANA Roles",
				Propertytype: "type3",
				Value:        "some hana value",
			},
		},
	}, nil)

	mockWebService.EXPECT().GetProcessList().Return(&sapcontrol.GetProcessListResponse{
		Processes: []*sapcontrol.OSProcess{
			{
				Name:        "enserver",
				Description: "foobar",
				Dispstatus:  sapcontrol.STATECOLOR_GREEN,
				Textstatus:  "Running",
				Starttime:   "",
				Elapsedtime: "",
				Pid:         30787,
			},
			{
				Name:        "msg_server",
				Description: "foobar2",
				Dispstatus:  sapcontrol.STATECOLOR_YELLOW,
				Textstatus:  "Stopping",
				Starttime:   "",
				Elapsedtime: "",
				Pid:         30786,
			},
		},
	}, nil)

	mockWebService.EXPECT().GetSystemInstanceList().Return(&sapcontrol.GetSystemInstanceListResponse{
		Instances: []*sapcontrol.SAPInstance{
			{
				Hostname:      "host1",
				InstanceNr:    0,
				HttpPort:      50013,
				HttpsPort:     50014,
				StartPriority: "0.3",
				Features:      "some features",
				Dispstatus:    sapcontrol.STATECOLOR_GREEN,
			},
			{
				Hostname:      "host2",
				InstanceNr:    1,
				HttpPort:      50113,
				HttpsPort:     50114,
				StartPriority: "0.3",
				Features:      "some other features",
				Dispstatus:    sapcontrol.STATECOLOR_YELLOW,
			},
		},
	}, nil)

	sapSystem, _ := NewSAPSystem(mockWebService)

	expectedSystem := SAPSystem{
		webService: mockWebService,
		Id:         "",
		Type:       "HANA",
		Processes: map[string]*sapcontrol.OSProcess{
			"enserver": &sapcontrol.OSProcess{
				Name:        "enserver",
				Description: "foobar",
				Dispstatus:  sapcontrol.STATECOLOR_GREEN,
				Textstatus:  "Running",
				Starttime:   "",
				Elapsedtime: "",
				Pid:         30787,
			},
			"msg_server": &sapcontrol.OSProcess{
				Name:        "msg_server",
				Description: "foobar2",
				Dispstatus:  sapcontrol.STATECOLOR_YELLOW,
				Textstatus:  "Stopping",
				Starttime:   "",
				Elapsedtime: "",
				Pid:         30786,
			},
		},
		Properties: map[string]*sapcontrol.InstanceProperty{
			"prop1": &sapcontrol.InstanceProperty{
				Property:     "prop1",
				Propertytype: "type1",
				Value:        "value1",
			},
			"SAPSYSTEMNAME": &sapcontrol.InstanceProperty{
				Property:     "SAPSYSTEMNAME",
				Propertytype: "string",
				Value:        "PRD",
			},
			"HANA Roles": &sapcontrol.InstanceProperty{
				Property:     "HANA Roles",
				Propertytype: "type3",
				Value:        "some hana value",
			},
		},
		Instances: map[string]*sapcontrol.SAPInstance{
			"host1": &sapcontrol.SAPInstance{
				Hostname:      "host1",
				InstanceNr:    0,
				HttpPort:      50013,
				HttpsPort:     50014,
				StartPriority: "0.3",
				Features:      "some features",
				Dispstatus:    sapcontrol.STATECOLOR_GREEN,
			},
			"host2": &sapcontrol.SAPInstance{
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

	assert.Equal(t, expectedSystem, sapSystem)
}

func TestFindSystemsNotFound(t *testing.T) {
	appFS := afero.NewMemMapFs()
	// create test files and directories
	appFS.MkdirAll("/usr/sap/", 0755)

	systems, _ := findSystems(appFS)

	assert.Equal(t, []string{}, systems)
}

func TestFindSystems(t *testing.T) {
	appFS := afero.NewMemMapFs()
	// create test files and directories
	appFS.MkdirAll("/usr/sap/PRD/HDB00", 0755)
	appFS.MkdirAll("/usr/sap/PRD/HDB01", 0755)
	appFS.MkdirAll("/usr/sap/DEV/ASCS02", 0755)
	appFS.MkdirAll("/usr/sap/DEV1/ASCS02", 0755)
	appFS.MkdirAll("/usr/sap/DEV/SYS/BLA12", 0755)
	appFS.MkdirAll("/usr/sap/DEV/PRD0", 0755)

	systems, _ := findSystems(appFS)
	sort.Strings(systems)
	assert.Equal(t, []string{"00", "01", "02"}, systems)
}
