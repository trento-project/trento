package mocks

// TODO: this needs to be redesigned to capture what we are discovering about a Host
func NewDiscoveredHostMock() struct {
	HostIpAddresses string
	HostName        string
} {
	return struct {
		HostIpAddresses string
		HostName        string
	}{
		"10.1.1.4,10.1.1.5,10.1.1.6",
		"thehostnamewherethediscoveryhappened",
	}
}
