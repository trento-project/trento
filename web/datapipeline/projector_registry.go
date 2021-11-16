package datapipeline

import "gorm.io/gorm"

type ProjectorRegistry []Projector

// InitInitProjectorsRegistry initialize the ProjectorRegistry
func InitProjectorsRegistry(db *gorm.DB) ProjectorRegistry {
	return ProjectorRegistry{
		NewClustersProjector(db),
		NewHostsProjector(db),
		NewHostTelemetryProjector(db),
		NewSlesSubscriptionsProjector(db),
		NewSAPSystemsProjector(db),
	}
}
