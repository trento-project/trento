package hosts

type DiscoveredHost struct {
	HostIpAddresses []string `json:"ip_addresses"`
	HostName        string   `json:"hostname"`
	CPUCount        int      `json:"cpu_count"`
	SocketCount     int      `json:"socket_count"`
	TotalMemoryMB   int      `json:"total_memory_mb"`
	AgentVersion    string   `json:"agent_version"`
}
