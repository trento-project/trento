package discovery

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/agent/discovery/collector"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/version"
)

const HostDiscoveryId string = "host_discovery"
const HostDiscoveryMinPeriod time.Duration = 1 * time.Second

type HostDiscovery struct {
	id              string
	sshAddress      string
	collectorClient collector.Client
	host            string
	interval        time.Duration
}

func NewHostDiscovery(collectorClient collector.Client, config DiscoveriesConfig) Discovery {
	d := HostDiscovery{}
	d.id = HostDiscoveryId
	d.collectorClient = collectorClient
	d.host, _ = os.Hostname()
	d.interval = config.DiscoveriesPeriodsConfig.Host
	d.sshAddress = config.SSHAddress
	return d
}

func (d HostDiscovery) GetId() string {
	return d.id
}

func (d HostDiscovery) GetInterval() time.Duration {
	return d.interval
}

// Execute one iteration of a discovery and publish to the collector
func (d HostDiscovery) Discover() (string, error) {
	ipAddresses, err := getHostIpAddresses()
	if err != nil {
		return "", err
	}

	host := hosts.DiscoveredHost{
		SSHAddress:      d.sshAddress,
		OSVersion:       getOSVersion(),
		HostIpAddresses: ipAddresses,
		HostName:        d.host,
		CPUCount:        getLogicalCPUs(),
		SocketCount:     getCPUSocketCount(),
		TotalMemoryMB:   getTotalMemoryMB(),
		AgentVersion:    version.Version,
	}

	err = d.collectorClient.Publish(d.id, host)
	if err != nil {
		log.Debugf("Error while sending host discovery to data collector: %s", err)
		return "", err
	}

	return fmt.Sprintf("Host with name: %s successfully discovered", d.host), nil
}

func getHostIpAddresses() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ipAddrList := make([]string, 0)

	for _, inter := range interfaces {
		addrs, err := inter.Addrs()
		if err != nil {
			continue
		}

		for _, ipaddr := range addrs {
			ipv4Addr, _, _ := net.ParseCIDR(ipaddr.String())
			ipAddrList = append(ipAddrList, ipv4Addr.String())
		}
	}

	return ipAddrList, nil
}

func getOSVersion() string {
	infoStat, err := host.Info()
	if err != nil {
		log.Errorf("Error while getting host info: %s", err)
	}
	return infoStat.PlatformVersion
}

func getTotalMemoryMB() int {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Errorf("Error while getting memory info: %s", err)
	}
	return int(v.Total) / 1024 / 1024
}

func getLogicalCPUs() int {
	logical, err := cpu.Counts(true)
	if err != nil {
		log.Errorf("Error while getting logical CPU count: %s", err)
	}
	return logical
}

func getCPUSocketCount() int {
	info, err := cpu.Info()

	if err != nil {
		log.Errorf("Error while getting CPU info: %s", err)
		return 0
	}

	// Get the last CPU info and get the physical ID of it
	lastCpuInfo := info[len(info)-1]

	physicalID, err := strconv.Atoi(lastCpuInfo.PhysicalID)

	if err != nil {
		log.Errorf("Unable to convert CPU socket count: %s", err)
		return 0
	}

	// Increase by one as physicalIDs start in zero
	return physicalID + 1
}
