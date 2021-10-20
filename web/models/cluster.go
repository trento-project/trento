package models

import "github.com/lib/pq"

const (
	ClusterTypeScaleUp  = "HANA scale-up"
	ClusterTypeScaleOut = "HANA scale-out"
	ClusterTypeUnknown  = "Unknown"
)

type Cluster struct {
	ID                string
	Name              string
	ClusterType       string
	SIDs              pq.StringArray `gorm:"column:sids; type:text[]"`
	ResourcesNumber   int
	HostsNumber       int
	Health            string   `gorm:"-"`
	Tags              []string `gorm:"-"`
	HasDuplicatedName bool     `gorm:"-"`
}

type ClusterList []*Cluster

// GetAllSIDs returns all the deduplicated
// and non empty sids in the cluster used by the frontend to show selectable filters
func (clusterList ClusterList) GetAllSIDs() []string {
	var sids []string
	set := make(map[string]struct{})

	for _, c := range clusterList {
		for _, sid := range c.SIDs {
			if sid == "" {
				continue
			}

			_, ok := set[sid]
			if !ok {
				set[sid] = struct{}{}
				sids = append(sids, sid)
			}
		}
	}

	return sids
}

func (clusterList ClusterList) GetAllTags() []string {
	var tags []string
	set := make(map[string]struct{})

	for _, c := range clusterList {
		for _, tag := range c.Tags {
			_, ok := set[tag]
			if !ok {
				set[tag] = struct{}{}
				tags = append(tags, tag)
			}
		}
	}

	return tags
}

func (clusterList ClusterList) GetAllClusterTypes() []string {
	var clusterTypes []string
	set := make(map[string]struct{})

	for _, c := range clusterList {
		_, ok := set[c.ClusterType]
		if !ok {
			set[c.ClusterType] = struct{}{}
			clusterTypes = append(clusterTypes, c.ClusterType)
		}
	}

	return clusterTypes
}
