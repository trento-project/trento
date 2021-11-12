package models

const (
	HostHealthPassing  = "passing"
	HostHealthWarning  = "warning"
	HostHealthCritical = "critical"
	HostHealthUnknown  = ""
)

type Host struct {
	ID            string
	Name          string
	Health        string
	IPAddresses   []string
	CloudProvider string
	ClusterID     string
	ClusterName   string
	SIDs          []string
	AgentVersion  string
	Tags          []string
}

type HostList []*Host
