package mocks

import "github.com/trento-project/trento/internal/hosts"

func NewDiscoveredHostMock() hosts.DiscoveredHost {
	return hosts.DiscoveredHost{
		SSHAddress:      "10.2.2.22",
		OSVersion:       "15-SP2",
		HostIpAddresses: []string{"10.1.1.4", "10.1.1.5", "10.1.1.6"},
		HostName:        "thehostnamewherethediscoveryhappened",
		CPUCount:        2,
		SocketCount:     1,
		TotalMemoryMB:   4096,
		AgentVersion:    "trento-agent-version",
	}
}
