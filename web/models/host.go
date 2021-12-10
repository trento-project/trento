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
	ClusterType   string
	SAPSystems    []*SAPSystem
	AgentVersion  string
	Tags          []string
}

type HostList []*Host
