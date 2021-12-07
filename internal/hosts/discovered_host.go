package hosts

type DiscoveredHost struct {
	AgentBindIP     string   `json:"agent_bind_ip"`
	OSVersion       string   `json:"os_version"`
	HostIpAddresses []string `json:"ip_addresses"`
	HostName        string   `json:"hostname"`
	CPUCount        int      `json:"cpu_count"`
	SocketCount     int      `json:"socket_count"`
	TotalMemoryMB   int      `json:"total_memory_mb"`
	AgentVersion    string   `json:"agent_version"`
}
