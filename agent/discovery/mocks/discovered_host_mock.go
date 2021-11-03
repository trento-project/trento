package mocks

import "github.com/trento-project/trento/internal/hosts"

func NewDiscoveredHostMock() hosts.DiscoveredHost {
	return hosts.DiscoveredHost{
		HostIpAddresses: []string{"10.1.1.4", "10.1.1.5", "10.1.1.6"},
		HostName:        "thehostnamewherethediscoveryhappened",
		CPUCount:        64,
		SocketCount:     16,
		TotalMemoryMB:   4096,
		AgentVersion:    "trento-agent-version",
	}
}
