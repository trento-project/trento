package datapipeline

import "gorm.io/gorm"

type ProjectorRegistry []Projector

// InitInitProjectorsRegistry initialize the ProjectorRegistry
func InitProjectorsRegistry(db *gorm.DB) ProjectorRegistry {
	clusterListProjector := NewProjector("cluster_list", db)
	clusterListProjector.AddHandler(ClusterDiscovery, ClusterListHandler)

	return ProjectorRegistry{
		clusterListProjector,
	}
}
