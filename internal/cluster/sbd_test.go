package cluster

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockCmdError(command string, args ...string) *exec.Cmd {
	return exec.Command(command, "error")
}

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

func mockSbdList(command string, args ...string) *exec.Cmd {
	cmd := `0	hana01	clear
1	hana02	clear`
	return exec.Command("echo", cmd)
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
	sbdDumpExecCommand = mockCmdError

	dump, err := sbdDump("/bin/sbdErr", "/dev/vdc")

	expectedDump := SBDDump{}

	assert.Equal(t, expectedDump, dump)
	assert.EqualError(t, err, "sbd dump command error: fork/exec /bin/sbdErr: no such file or directory")
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
	sbdListExecCommand = mockCmdError

	list, err := sbdList("/bin/sbdErr", "/dev/vdc")

	expectedList := []*SBDNode{}

	assert.Equal(t, expectedList, list)
	assert.EqualError(t, err, "sbd list command error: fork/exec /bin/sbdErr: no such file or directory")
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

	sbdDumpExecCommand = mockCmdError

	err := s.LoadDeviceData()

	expectedDevice := NewSBDDevice("/bin/sbdErr", "/dev/vdc")
	expectedDevice.Status = "unhealthy"

	assert.Equal(t, expectedDevice, s)
	assert.EqualError(t, err, "sbd dump command error: fork/exec /bin/sbdErr: no such file or directory")
}

func TestLoadDeviceDataListError(t *testing.T) {
	s := NewSBDDevice("/bin/sbdErr", "/dev/vdc")

	sbdDumpExecCommand = mockSbdDump
	sbdListExecCommand = mockCmdError

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

	assert.Equal(t, expectedDevice, s)
	assert.EqualError(t, err, "sbd list command error: fork/exec /bin/sbdErr: no such file or directory")
}

func TestGetSBDConfig(t *testing.T) {
	sbdConfig, err := getSBDConfig("../../test/sbd_config")

	expectedConfig := map[string]string{
		"SBD_PACEMAKER":           "yes",
		"SBD_STARTMODE":           "always",
		"SBD_DELAY_START":         "no",
		"SBD_WATCHDOG_DEV":        "/dev/watchdog",
		"SBD_WATCHDOG_TIMEOUT":    "5",
		"SBD_TIMEOUT_ACTION":      "flush,reboot",
		"SBD_MOVE_TO_ROOT_CGROUP": "auto",
		"SBD_DEVICE":              "/dev/vdc;/dev/vdb",
	}

	assert.Equal(t, expectedConfig, sbdConfig)
	assert.NoError(t, err)
}

func TestGetSBDConfigError(t *testing.T) {
	sbdConfig, err := getSBDConfig("notexist")

	expectedConfig := map[string]string(nil)

	assert.Equal(t, expectedConfig, sbdConfig)
	assert.EqualError(t, err, "could not open sbd config file open notexist: no such file or directory")
}

func TestNewSBD(t *testing.T) {
	sbdDumpExecCommand = mockSbdDump
	sbdListExecCommand = mockSbdList

	s, err := NewSBD("mycluster", "/bin/sbd", "../../test/sbd_config")

	expectedSbd := SBD{
		cluster: "mycluster",
		Config: map[string]string{
			"SBD_PACEMAKER":           "yes",
			"SBD_STARTMODE":           "always",
			"SBD_DELAY_START":         "no",
			"SBD_WATCHDOG_DEV":        "/dev/watchdog",
			"SBD_WATCHDOG_TIMEOUT":    "5",
			"SBD_TIMEOUT_ACTION":      "flush,reboot",
			"SBD_MOVE_TO_ROOT_CGROUP": "auto",
			"SBD_DEVICE":              "/dev/vdc;/dev/vdb",
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

	assert.Equal(t, SBD{cluster: "mycluster"}, s)
	assert.EqualError(t, err, "could not find SBD_DEVICE entry in sbd config file")
}
