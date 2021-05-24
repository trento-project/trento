package cluster

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	SBDPath            = "/usr/sbin/sbd"
	SBDConfigPath      = "/etc/sysconfig/sbd"
	SBDStatusUnknown   = "unknown"
	SBDStatusUnhealthy = "unhealthy"
	SBDStatusHealthy   = "healthy"
)

type SBD struct {
	cluster string
	Devices []*SBDDevice      `mapstructure:"devices,omitempty"`
	Config  map[string]string `mapstructure:"config,omitempty"`
}

type SBDDevice struct {
	sbdPath string
	Device  string     `mapstructure:"device,omitempty"`
	Status  string     `mapstructure:"status,omitempty"`
	Dump    SBDDump    `mapstructure:"dump,omitempty"`
	List    []*SBDNode `mapstructure:"list,omitempty"`
}

type SBDDump struct {
	Header          string `mapstructure:"header,omitempty"`
	Uuid            string `mapstructure:"uuid,omitempty"`
	Slots           int    `mapstructure:"slots,omitempty"`
	SectorSize      int    `mapstructure:"sectorsize,omitempty"`
	TimeoutWatchdog int    `mapstructure:"timeoutwatchdog,omitempty"`
	TimeoutAllocate int    `mapstructure:"timeoutallocate,omitempty"`
	TimeoutLoop     int    `mapstructure:"timeoutloop,omitempty"`
	TimeoutMsgwait  int    `mapstructure:"timeoutmsgwait,omitempty"`
}

type SBDNode struct {
	Id     int    `mapstructure:"id,omitempty"`
	Name   string `mapstructure:"name,omitempty"`
	Status string `mapstructure:"status,omitempty"`
}

func NewSBD(cluster string) (SBD, error) {
	var s = SBD{cluster: cluster}

	c, err := getSBDConfig(SBDConfigPath)
	if err != nil {
		return s, err
	} else if _, ok := c["SBD_DEVICE"]; !ok {
		return s, fmt.Errorf("could not find SBD_DEVICE entry in sbd config file")
	}
	s.Config = c

	for _, device := range strings.Split(c["SBD_DEVICE"], ";") {
		sbdDevice := NewSBDDevice(SBDPath, device)
		err := sbdDevice.LoadDeviceData()
		if err != nil {
			log.Printf("Error getting sbd information: %s", err)
			continue
		}
		s.Devices = append(s.Devices, &sbdDevice)
	}

	return s, nil
}

func getSBDConfig(sbdConfigPath string) (map[string]string, error) {
	configMap := make(map[string]string)

	sbdConfFile, err := os.Open(sbdConfigPath)
	if err != nil {
		return nil, fmt.Errorf("could not open sbd config file %s", err)
	}

	defer sbdConfFile.Close()

	sbdConfigRaw, err := ioutil.ReadAll(sbdConfFile)

	if err != nil {
		return nil, fmt.Errorf("could not read sbd config file %s", err)
	}

	// Loop through sbd list output and find for matches
	r := regexp.MustCompile(`(\S+)\s*=(\S+)`)
	values := r.FindAllStringSubmatch(string(sbdConfigRaw), -1)
	for _, match := range values {
		// I was not able to create a regular expression which excludes lines starting with #...
		if strings.HasPrefix(match[1], "#") {
			continue
		}
		configMap[match[1]] = match[2]
	}
	return configMap, nil
}

func NewSBDDevice(sbdPath string, device string) SBDDevice {
	return SBDDevice{
		sbdPath: sbdPath,
		Device:  device,
		Status:  SBDStatusUnknown,
	}
}

func (s *SBDDevice) LoadDeviceData() error {
	dump, err := sbdDump(s.sbdPath, s.Device)
	if err != nil {
		s.Status = SBDStatusUnhealthy
		return err
	}
	s.Dump = dump
	s.Status = SBDStatusHealthy

	list, err := sbdList(s.sbdPath, s.Device)
	if err != nil {
		return err
	}
	s.List = list

	return nil
}

func assignPatternResult(text string, pattern string) []string {
	r := regexp.MustCompile(pattern)
	match := r.FindAllStringSubmatch(text, -1)
	if len(match) > 0 {
		return match[0]
	} else {
		// Retrun empty information if pattern is not found
		return []string{"", ""}
	}
}

// Possible output
//==Dumping header on disk /dev/vdc
//Header version     : 2.1
//UUID               : 541bdcea-16af-44a4-8ab9-6a98602e65ca
//Number of slots    : 255
//Sector size        : 512
//Timeout (watchdog) : 5
//Timeout (allocate) : 2
//Timeout (loop)     : 1
//Timeout (msgwait)  : 10
//==Header on disk /dev/vdc is dumped
func sbdDump(sbdPath string, device string) (SBDDump, error) {
	var dump = SBDDump{}

	sbdDump, err := exec.Command(sbdPath, "-d", device, "dump").Output()
	if err != nil {
		return dump, errors.Wrap(err, "sbd dump command error")
	}
	sbdDumpStr := string(sbdDump)

	dump.Header = assignPatternResult(sbdDumpStr, `Header version *: (.*)`)[1]
	dump.Uuid = assignPatternResult(sbdDumpStr, `UUID *: (.*)`)[1]
	dump.Slots, _ = strconv.Atoi(assignPatternResult(sbdDumpStr, `Number of slots *: (.*)`)[1])
	dump.SectorSize, _ = strconv.Atoi(assignPatternResult(sbdDumpStr, `Number of slots *: (.*)`)[1])
	dump.TimeoutWatchdog, _ = strconv.Atoi(assignPatternResult(sbdDumpStr, `Timeout \(watchdog\) *: (.*)`)[1])
	dump.TimeoutAllocate, _ = strconv.Atoi(assignPatternResult(sbdDumpStr, `Timeout \(allocate\) *: (.*)`)[1])
	dump.TimeoutLoop, _ = strconv.Atoi(assignPatternResult(sbdDumpStr, `Timeout \(loop\) *: (.*)`)[1])
	dump.TimeoutMsgwait, _ = strconv.Atoi(assignPatternResult(sbdDumpStr, `Timeout \(msgwait\) *: (.*)`)[1])

	return dump, nil
}

// Possible output
//0	hana01	clear
//1	hana02	clear
func sbdList(sbdPath string, device string) ([]*SBDNode, error) {
	var list = []*SBDNode{}

	output, err := exec.Command(sbdPath, "-d", device, "list").Output()
	if err != nil {
		return list, errors.Wrap(err, "sbd list command error")
	}

	// Loop through sbd list output and find for matches
	r := regexp.MustCompile(`(\d+)\s+(\S+)\s+(\S+)`)
	values := r.FindAllStringSubmatch(string(output), -1)
	for _, match := range values {
		// Continue loop if all the groups are not found
		if len(match) != 4 {
			continue
		}

		id, _ := strconv.Atoi(match[1])
		node := &SBDNode{
			Id:     id,
			Name:   match[2],
			Status: match[3],
		}
		list = append(list, node)
	}

	return list, nil
}
