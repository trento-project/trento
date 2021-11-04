package datapipeline

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	ClusterDiscovery   = "ha_cluster_discovery"
	SAPsystemDiscovery = "sap_system_discovery"
)

type DataCollectedEvent struct {
	ID            int64
	CreatedAt     time.Time
	AgentID       string         `json:"agent_id" binding:"required"`
	DiscoveryType string         `json:"discovery_type" binding:"required"`
	Payload       datatypes.JSON `json:"payload" binding:"required"`
}

func PruneEvents(olderThan time.Duration, db *gorm.DB) error {
	return db.Delete(DataCollectedEvent{}, "created_at < ?", time.Now().Add(-olderThan)).Error
}
