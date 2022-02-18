package cluster

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockSbdDump(command string, args ...string) *exec.Cmd {
	cmd := `==Dumping header on disk /dev/vdc
Header version     : 2.1
UUID               : 541bdcea-16af-44a4-8ab9-6a98602e65ca
Number of slots    : 255
Sector size        : 512
Timeout (watchdog) : 5
Timeout (allocate) : 2
Timeout (loop)     : 1
Timeout (msgwait)  : 10
==Header on disk /dev/vdc is dumped`
	return exec.Command("echo", cmd)
}

func mockSbdDumpErr(command string, args ...string) *exec.Cmd {
	cmd := `==Dumping header on disk /dev/vdc
Header version     : 2.1
UUID               : 541bdcea-16af-44a4-8ab9-6a98602e65ca
==Number of slots on disk /dev/vdb NOT dumped
sbd failed; please check the logs.`

	script := fmt.Sprintf("echo \"%s\" && exit 1", cmd)

	return exec.Command("bash", "-c", script)
}

func mockSbdList(command string, args ...string) *exec.Cmd {
	cmd := `0	hana01	clear
1	hana02	clear`
	return exec.Command("echo", cmd)
}

func mockSbdListErr(command string, args ...string) *exec.Cmd {
	cmd := `== disk /dev/vdxx unreadable!
sbd failed; please check the logs.`

	script := fmt.Sprintf("echo \"%s\" && exit 1", cmd)

	return exec.Command("bash", "-c", script)
}

func TestSbdDump(t *testing.T) {
	sbdDumpExecCommand = mockSbdDump

	dump, err := sbdDump("/bin/sbd", "/dev/vdc")

	expectedDump := SBDDump{
		Header:          "2.1",
		Uuid:            "541bdcea-16af-44a4-8ab9-6a98602e65ca",
		Slots:           255,
		SectorSize:      512,
		TimeoutWatchdog: 5,
		TimeoutAllocate: 2,
		TimeoutLoop:     1,
		TimeoutMsgwait:  10,
	}

	assert.Equal(t, expectedDump, dump)
	assert.NoError(t, err)
}

func TestSbdDumpError(t *testing.T) {
	sbdDumpExecCommand = mockSbdDumpErr

	dump, err := sbdDump("/bin/sbd", "/dev/vdc")

	expectedDump := SBDDump{
		Header:          "2.1",
		Uuid:            "541bdcea-16af-44a4-8ab9-6a98602e65ca",
		Slots:           0,
		SectorSize:      0,
		TimeoutWatchdog: 0,
		TimeoutAllocate: 0,
		TimeoutLoop:     0,
		TimeoutMsgwait:  0,
	}

	assert.Equal(t, expectedDump, dump)
	assert.EqualError(t, err, "sbd dump command error: exit status 1")
}

func TestSbdList(t *testing.T) {
	sbdListExecCommand = mockSbdList

	list, err := sbdList("/bin/sbd", "/dev/vdc")

	expectedList := []*SBDNode{
		&SBDNode{
			Id:     0,
			Name:   "hana01",
			Status: "clear",
		},
		&SBDNode{
			Id:     1,
			Name:   "hana02",
			Status: "clear",
		},
	}

	assert.Equal(t, expectedList, list)
	assert.NoError(t, err)
}

func TestSbdListError(t *testing.T) {
	sbdListExecCommand = mockSbdListErr

	list, err := sbdList("/bin/sbd", "/dev/vdc")

	expectedList := []*SBDNode{}

	assert.Equal(t, expectedList, list)
	assert.EqualError(t, err, "sbd list command error: exit status 1")
}

func TestLoadDeviceData(t *testing.T) {
	s := NewSBDDevice("/bin/sbd", "/dev/vdc")

	sbdDumpExecCommand = mockSbdDump
	sbdListExecCommand = mockSbdList

	err := s.LoadDeviceData()

	expectedDevice := NewSBDDevice("/bin/sbd", "/dev/vdc")
	expectedDevice.Status = "healthy"
	expectedDevice.Dump = SBDDump{
		Header:          "2.1",
		Uuid:            "541bdcea-16af-44a4-8ab9-6a98602e65ca",
		Slots:           255,
		SectorSize:      512,
		TimeoutWatchdog: 5,
		TimeoutAllocate: 2,
		TimeoutLoop:     1,
		TimeoutMsgwait:  10,
	}
	expectedDevice.List = []*SBDNode{
		&SBDNode{
			Id:     0,
			Name:   "hana01",
			Status: "clear",
		},
		&SBDNode{
			Id:     1,
			Name:   "hana02",
			Status: "clear",
		},
	}

	assert.Equal(t, expectedDevice, s)
	assert.NoError(t, err)
}

func TestLoadDeviceDataDumpError(t *testing.T) {
	s := NewSBDDevice("/bin/sbdErr", "/dev/vdc")

	sbdDumpExecCommand = mockSbdDumpErr

	err := s.LoadDeviceData()

	expectedDevice := NewSBDDevice("/bin/sbdErr", "/dev/vdc")
	expectedDevice.Status = "unhealthy"

	expectedDevice.Dump = SBDDump{
		Header:          "2.1",
		Uuid:            "541bdcea-16af-44a4-8ab9-6a98602e65ca",
		Slots:           0,
		SectorSize:      0,
		TimeoutWatchdog: 0,
		TimeoutAllocate: 0,
		TimeoutLoop:     0,
		TimeoutMsgwait:  0,
	}

	expectedDevice.List = []*SBDNode{
		&SBDNode{
			Id:     0,
			Name:   "hana01",
			Status: "clear",
		},
		&SBDNode{
			Id:     1,
			Name:   "hana02",
			Status: "clear",
		},
	}

	assert.Equal(t, expectedDevice, s)
	assert.EqualError(t, err, "sbd dump command error: exit status 1")
}

func TestLoadDeviceDataListError(t *testing.T) {
	s := NewSBDDevice("/bin/sbdErr", "/dev/vdc")

	sbdDumpExecCommand = mockSbdDump
	sbdListExecCommand = mockSbdListErr

	err := s.LoadDeviceData()

	expectedDevice := NewSBDDevice("/bin/sbdErr", "/dev/vdc")
	expectedDevice.Status = "healthy"
	expectedDevice.Dump = SBDDump{
		Header:          "2.1",
		Uuid:            "541bdcea-16af-44a4-8ab9-6a98602e65ca",
		Slots:           255,
		SectorSize:      512,
		TimeoutWatchdog: 5,
		TimeoutAllocate: 2,
		TimeoutLoop:     1,
		TimeoutMsgwait:  10,
	}

	expectedDevice.List = []*SBDNode{}

	assert.Equal(t, expectedDevice, s)
	assert.EqualError(t, err, "sbd list command error: exit status 1")
}

func TestLoadDeviceDataError(t *testing.T) {
	s := NewSBDDevice("/bin/sbdErr", "/dev/vdc")

	sbdDumpExecCommand = mockSbdDumpErr
	sbdListExecCommand = mockSbdListErr

	err := s.LoadDeviceData()

	expectedDevice := NewSBDDevice("/bin/sbdErr", "/dev/vdc")
	expectedDevice.Status = "unhealthy"

	expectedDevice.Dump = SBDDump{
		Header:          "2.1",
		Uuid:            "541bdcea-16af-44a4-8ab9-6a98602e65ca",
		Slots:           0,
		SectorSize:      0,
		TimeoutWatchdog: 0,
		TimeoutAllocate: 0,
		TimeoutLoop:     0,
		TimeoutMsgwait:  0,
	}

	expectedDevice.List = []*SBDNode{}

	assert.Equal(t, expectedDevice, s)
	assert.EqualError(t, err, "sbd dump command error: exit status 1;sbd list command error: exit status 1")
}

func TestGetSBDConfig(t *testing.T) {
	sbdConfig, err := getSBDConfig("../../test/sbd_config")

	expectedConfig := map[string]interface{}{
		"SBD_PACEMAKER":           "yes",
		"SBD_STARTMODE":           "always",
		"SBD_DELAY_START":         "no",
		"SBD_WATCHDOG_DEV":        "/dev/watchdog",
		"SBD_WATCHDOG_TIMEOUT":    "5",
		"SBD_TIMEOUT_ACTION":      "flush,reboot",
		"SBD_MOVE_TO_ROOT_CGROUP": "auto",
		"SBD_DEVICE":              "/dev/vdc;/dev/vdb",
		"TEST":                    "Value",
		"TEST2":                   "Value2",
	}

	assert.Equal(t, expectedConfig, sbdConfig)
	assert.NoError(t, err)
}

func TestGetSBDConfigError(t *testing.T) {
	sbdConfig, err := getSBDConfig("notexist")

	expectedConfig := map[string]interface{}(nil)

	assert.Equal(t, expectedConfig, sbdConfig)
	assert.EqualError(t, err, "could not open sbd config file open notexist: no such file or directory")
}

func TestNewSBD(t *testing.T) {
	sbdDumpExecCommand = mockSbdDump
	sbdListExecCommand = mockSbdList

	s, err := NewSBD("mycluster", "/bin/sbd", "../../test/sbd_config")

	expectedSbd := SBD{
		cluster: "mycluster",
		Config: map[string]interface{}{
			"SBD_PACEMAKER":           "yes",
			"SBD_STARTMODE":           "always",
			"SBD_DELAY_START":         "no",
			"SBD_WATCHDOG_DEV":        "/dev/watchdog",
			"SBD_WATCHDOG_TIMEOUT":    "5",
			"SBD_TIMEOUT_ACTION":      "flush,reboot",
			"SBD_MOVE_TO_ROOT_CGROUP": "auto",
			"SBD_DEVICE":              "/dev/vdc;/dev/vdb",
			"TEST":                    "Value",
			"TEST2":                   "Value2",
		},
		Devices: []*SBDDevice{
			&SBDDevice{
				sbdPath: "/bin/sbd",
				Device:  "/dev/vdc",
				Status:  "healthy",
				Dump: SBDDump{
					Header:          "2.1",
					Uuid:            "541bdcea-16af-44a4-8ab9-6a98602e65ca",
					Slots:           255,
					SectorSize:      512,
					TimeoutWatchdog: 5,
					TimeoutAllocate: 2,
					TimeoutLoop:     1,
					TimeoutMsgwait:  10,
				},
				List: []*SBDNode{
					&SBDNode{
						Id:     0,
						Name:   "hana01",
						Status: "clear",
					},
					&SBDNode{
						Id:     1,
						Name:   "hana02",
						Status: "clear",
					},
				},
			},
			&SBDDevice{
				sbdPath: "/bin/sbd",
				Device:  "/dev/vdb",
				Status:  "healthy",
				Dump: SBDDump{
					Header:          "2.1",
					Uuid:            "541bdcea-16af-44a4-8ab9-6a98602e65ca",
					Slots:           255,
					SectorSize:      512,
					TimeoutWatchdog: 5,
					TimeoutAllocate: 2,
					TimeoutLoop:     1,
					TimeoutMsgwait:  10,
				},
				List: []*SBDNode{
					&SBDNode{
						Id:     0,
						Name:   "hana01",
						Status: "clear",
					},
					&SBDNode{
						Id:     1,
						Name:   "hana02",
						Status: "clear",
					},
				},
			},
		},
	}

	assert.Equal(t, expectedSbd, s)
	assert.NoError(t, err)
}

func TestNewSBDError(t *testing.T) {
	s, err := NewSBD("mycluster", "/bin/sbd", "../../test/sbd_config_no_device")

	expectedSbd := SBD{
		cluster: "mycluster",
		Config: map[string]interface{}{
			"SBD_PACEMAKER":           "yes",
			"SBD_STARTMODE":           "always",
			"SBD_DELAY_START":         "no",
			"SBD_WATCHDOG_DEV":        "/dev/watchdog",
			"SBD_WATCHDOG_TIMEOUT":    "5",
			"SBD_TIMEOUT_ACTION":      "flush,reboot",
			"SBD_MOVE_TO_ROOT_CGROUP": "auto",
		},
	}

	assert.Equal(t, expectedSbd, s)
	assert.EqualError(t, err, "could not find SBD_DEVICE entry in sbd config file")
}

func TestNewSBDUnhealthyDevices(t *testing.T) {
	sbdDumpExecCommand = mockSbdDumpErr
	sbdListExecCommand = mockSbdListErr

	s, err := NewSBD("mycluster", "/bin/sbd", "../../test/sbd_config")

	expectedSbd := SBD{
		cluster: "mycluster",
		Config: map[string]interface{}{
			"SBD_PACEMAKER":           "yes",
			"SBD_STARTMODE":           "always",
			"SBD_DELAY_START":         "no",
			"SBD_WATCHDOG_DEV":        "/dev/watchdog",
			"SBD_WATCHDOG_TIMEOUT":    "5",
			"SBD_TIMEOUT_ACTION":      "flush,reboot",
			"SBD_MOVE_TO_ROOT_CGROUP": "auto",
			"SBD_DEVICE":              "/dev/vdc;/dev/vdb",
			"TEST":                    "Value",
			"TEST2":                   "Value2",
		},
		Devices: []*SBDDevice{
			&SBDDevice{
				sbdPath: "/bin/sbd",
				Device:  "/dev/vdc",
				Status:  "unhealthy",
				Dump: SBDDump{
					Header:          "2.1",
					Uuid:            "541bdcea-16af-44a4-8ab9-6a98602e65ca",
					Slots:           0,
					SectorSize:      0,
					TimeoutWatchdog: 0,
					TimeoutAllocate: 0,
					TimeoutLoop:     0,
					TimeoutMsgwait:  0,
				},
				List: []*SBDNode{},
			},
			&SBDDevice{
				sbdPath: "/bin/sbd",
				Device:  "/dev/vdb",
				Status:  "unhealthy",
				Dump: SBDDump{
					Header:          "2.1",
					Uuid:            "541bdcea-16af-44a4-8ab9-6a98602e65ca",
					Slots:           0,
					SectorSize:      0,
					TimeoutWatchdog: 0,
					TimeoutAllocate: 0,
					TimeoutLoop:     0,
					TimeoutMsgwait:  0,
				},
				List: []*SBDNode{},
			},
		},
	}

	assert.Equal(t, expectedSbd, s)
	assert.NoError(t, err)
}

func TestNewSBDQuotedDevices(t *testing.T) {
	sbdDumpExecCommand = mockSbdDump
	sbdListExecCommand = mockSbdList

	s, err := NewSBD("mycluster", "/bin/sbd", "../../test/sbd_config_quoted_devices")

	assert.Equal(t, len(s.Devices), 2)
	assert.Equal(t, "/dev/vdc", s.Devices[0].Device)
	assert.Equal(t, "/dev/vdb", s.Devices[1].Device)
	assert.NoError(t, err)
}
