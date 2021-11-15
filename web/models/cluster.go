package models

const (
	ClusterTypeScaleUp  = "HANA scale-up"
	ClusterTypeScaleOut = "HANA scale-out"
	ClusterTypeUnknown  = "Unknown"
)

type Cluster struct {
	ID              string
	Name            string
	ClusterType     string
	SIDs            []string
	ResourcesNumber int
	HostsNumber     int
	Health          string
	Tags            []string
	// TODO: this is frontend specific, should be removed
	HasDuplicatedName bool
}

type ClusterList []*Cluster
