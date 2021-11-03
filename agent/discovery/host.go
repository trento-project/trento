package discovery

import (
	"fmt"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/agent/discovery/collector"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/version"
	"github.com/zcalusic/sysinfo"
)

const HostDiscoveryId string = "host_discovery"

type HostDiscovery struct {
	id        string
	discovery BaseDiscovery
}

func NewHostDiscovery(consulClient consul.Client, collectorClient collector.Client) HostDiscovery {
	d := HostDiscovery{}
	d.id = HostDiscoveryId
	d.discovery = NewDiscovery(consulClient, collectorClient)
	return d
}

func (h HostDiscovery) GetId() string {
	return h.id
}

// Execute one iteration of a discovery and store the result in the Consul KVStore.
func (h HostDiscovery) Discover() (string, error) {
	ipAddresses, err := getHostIpAddresses()
	if err != nil {
		return "", err
	}

	metadata := hosts.Metadata{
		HostIpAddresses: strings.Join(ipAddresses, ","),
	}
	err = metadata.Store(h.discovery.consulClient)
	if err != nil {
		return "", err
	}

	var si sysinfo.SysInfo
	si.GetSysInfo()

	host := hosts.DiscoveredHost{
		HostIpAddresses: ipAddresses,
		HostName:        h.discovery.host,
		CPUCount:        int(si.CPU.Cpus) * int(si.CPU.Cores),
		SocketCount:     int(si.CPU.Cpus),
		TotalMemoryMB:   int(si.Memory.Size),
		AgentVersion:    version.Version,
	}

	err = h.discovery.collectorClient.Publish(h.id, host)
	if err != nil {
		log.Debugf("Error while sending host discovery to data collector: %s", err)
		return "", err
	}

	return fmt.Sprintf("Host with name: %s successfully discovered", h.discovery.host), nil
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
