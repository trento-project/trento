package entities

import (
	"time"

	"github.com/lib/pq"
	"github.com/trento-project/trento/web/models"
)

type Host struct {
	AgentID       string `gorm:"primaryKey"`
	Name          string
	IPAddresses   pq.StringArray `gorm:"type:text[]"`
	CloudProvider string
	ClusterID     string
	ClusterName   string
	SIDs          pq.StringArray `gorm:"column:sids; type:text[]"`
	AgentVersion  string
	Tags          []models.Tag `gorm:"polymorphic:Resource;polymorphicValue:hosts"`
	UpdatedAt     time.Time
}

func (h *Host) ToModel() *models.Host {
	// TODO: move to Tags entity when we will have it
	var tags []string
	for _, tag := range h.Tags {
		tags = append(tags, tag.Value)
	}

	return &models.Host{
		ID:            h.AgentID,
		Name:          h.Name,
		IPAddresses:   h.IPAddresses,
		CloudProvider: h.CloudProvider,
		ClusterID:     h.ClusterID,
		ClusterName:   h.ClusterName,
		SIDs:          h.SIDs,
		AgentVersion:  h.AgentVersion,
		Tags:          tags,
	}
}
