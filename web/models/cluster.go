package models

import "time"

const (
	ClusterTypeHANAScaleUp  = "HANA scale-up"
	ClusterTypeHANAScaleOut = "HANA scale-out"
	ClusterTypeUnknown      = "Unknown"
	HANAStatusPrimary       = "Primary"
	HANAStatusSecondary     = "Secondary"
)

type Cluster struct {
	ID              string
	Name            string
	ClusterType     string
	SID             string
	ResourcesNumber int
	HostsNumber     int
	Health          string
	PassingCount    int
	WarningCount    int
	CriticalCount   int
	Tags            []string
	// TODO: this is frontend specific, should be removed
	HasDuplicatedName bool
	Details           interface{}
}

type ClusterList []*Cluster

type HANAClusterDetails struct {
	SystemReplicationMode          string
	SystemReplicationOperationMode string
	SecondarySyncState             string
	SRHealthState                  string
	CIBLastWritten                 time.Time
	FencingType                    string
	StoppedResources               []*ClusterResource
	Nodes                          ClusterNodes
	SBDDevices                     []*SBDDevice
}

type ClusterResource struct {
	ID        string
	Type      string
	Role      string
	Status    string
	FailCount int
}

type HANAClusterNode struct {
	HostID      string
	Name        string
	Site        string
	IPAddresses []string
	VirtualIPs  []string
	Health      string
	HANAStatus  string
	Attributes  map[string]string
	Resources   []*ClusterResource
}

type SBDDevice struct {
	Device string
}

type ClusterNodes []*HANAClusterNode

func (n ClusterNodes) GroupBySite() map[string]ClusterNodes {
	sites := make(map[string]ClusterNodes)
	for _, node := range n {
		sites[node.Site] = append(sites[node.Site], node)
	}

	return sites
}
