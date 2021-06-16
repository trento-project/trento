package sapsystem

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/internal/sapsystem/sapcontrol"
	"github.com/trento-project/trento/internal/sapsystem/sapcontrol/mocks"
)

var attemps int

func increaseAttemps() {
	attemps++
}

func fakeNewWebService(instNumber string) sapcontrol.WebService {
	var instance string

	mockWebService := new(mocks.WebService)

	defer increaseAttemps()

	if attemps == 0 {
		instance = "ASCS01"
	} else if attemps == 1 {
		instance = "ERS02"
	}

	mockWebService.On("GetInstanceProperties").Return(&sapcontrol.GetInstancePropertiesResponse{
		Properties: []*sapcontrol.InstanceProperty{
			{
				Property:     "SAPSYSTEMNAME",
				Propertytype: "string",
				Value:        "DEV",
			},
			{
				Property:     "INSTANCE_NAME",
				Propertytype: "string",
				Value:        instance,
			},
		},
	}, nil)

	mockWebService.On("GetProcessList").Return(&sapcontrol.GetProcessListResponse{
		Processes: []*sapcontrol.OSProcess{},
	}, nil)

	mockWebService.On("GetSystemInstanceList").Return(&sapcontrol.GetSystemInstanceListResponse{
		Instances: []*sapcontrol.SAPInstance{},
	}, nil)

	return mockWebService
}

func TestNewSAPSystem(t *testing.T) {

	newWebService = fakeNewWebService

	appFS := afero.NewMemMapFs()
	appFS.MkdirAll("/usr/sap/DEV/ASCS01", 0755)
	appFS.MkdirAll("/usr/sap/DEV/ERS02", 0755)

	system, err := NewSAPSystem(appFS, "/usr/sap/DEV")

	assert.Equal(t, Application, system.Type)
	assert.Contains(t, system.Instances, "ASCS01")
	assert.Contains(t, system.Instances, "ERS02")
	assert.NoError(t, err)
}

func TestNewSAPInstance(t *testing.T) {
	mockWebService := new(mocks.WebService)

	mockWebService.On("GetInstanceProperties").Return(&sapcontrol.GetInstancePropertiesResponse{
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
			{
				Property:     "INSTANCE_NAME",
				Propertytype: "string",
				Value:        "HDB00",
			},
		},
	}, nil)

	mockWebService.On("GetProcessList").Return(&sapcontrol.GetProcessListResponse{
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

	mockWebService.On("GetSystemInstanceList").Return(&sapcontrol.GetSystemInstanceListResponse{
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

	sapInstance, _ := NewSAPInstance(mockWebService)
	host, _ := os.Hostname()

	expectedInstance := &SAPInstance{
		Name: "HDB00",
		Type: Database,
		Host: host,
		SAPControl: &SAPControl{
			webService: mockWebService,
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
				"INSTANCE_NAME": &sapcontrol.InstanceProperty{
					Property:     "INSTANCE_NAME",
					Propertytype: "string",
					Value:        "HDB00",
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
		},
	}

	assert.Equal(t, expectedInstance, sapInstance)
}

func TestFindSystemsNotFound(t *testing.T) {
	appFS := afero.NewMemMapFs()
	// create test files and directories
	appFS.MkdirAll("/usr/sap/", 0755)
	appFS.MkdirAll("/usr/sap/DEV1/", 0755)

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
	assert.ElementsMatch(t, []string{"/usr/sap/PRD", "/usr/sap/DEV"}, systems)
}

func TestFindInstancesNotFound(t *testing.T) {
	appFS := afero.NewMemMapFs()
	// create test files and directories
	appFS.MkdirAll("/usr/sap/DEV/SYS/BLA12", 0755)

	instances, _ := findInstances(appFS, "/usr/sap/DEV")

	assert.Equal(t, [][]string{}, instances)
}

func TestFindInstances(t *testing.T) {
	appFS := afero.NewMemMapFs()
	// create test files and directories
	appFS.MkdirAll("/usr/sap/DEV/ASCS02", 0755)
	appFS.MkdirAll("/usr/sap/DEV/SYS/BLA12", 0755)
	appFS.MkdirAll("/usr/sap/DEV/PRD0", 0755)
	appFS.MkdirAll("/usr/sap/DEV/ERS10", 0755)

	instances, _ := findInstances(appFS, "/usr/sap/DEV")
	expectedInstance := [][]string{
		{"ASCS02", "02"},
		{"ERS10", "10"},
	}
	assert.ElementsMatch(t, expectedInstance, instances)
}
