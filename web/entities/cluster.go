package entities

import (
	"time"

	"github.com/lib/pq"
	"github.com/trento-project/trento/web/models"
)

type Cluster struct {
	ID              string `gorm:"primaryKey"`
	Name            string
	ClusterType     string
	SIDs            pq.StringArray `gorm:"column:sids; type:text[]"`
	ResourcesNumber int
	HostsNumber     int
	Tags            []models.Tag `gorm:"polymorphic:Resource;polymorphicValue:clusters"`
	UpdatedAt       time.Time
}

func (h *Cluster) ToModel() *models.Cluster {
	// TODO: move to Tags entity when we will have it
	var tags []string
	for _, tag := range h.Tags {
		tags = append(tags, tag.Value)
	}

	return &models.Cluster{
		ID:              h.ID,
		Name:            h.Name,
		ClusterType:     h.ClusterType,
		SIDs:            h.SIDs,
		ResourcesNumber: h.ResourcesNumber,
		HostsNumber:     h.HostsNumber,
		Tags:            tags,
	}
}
