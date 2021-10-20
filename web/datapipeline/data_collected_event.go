package datapipeline

import (
	"time"

	"gorm.io/datatypes"
)

const (
	ClusterDiscovery   = "clusterDiscovery"
	SAPsystemDiscovery = "sapsystemDiscovery"
)

type DataCollectedEvent struct {
	ID            int64
	CreatedAt     time.Time
	AgentID       string         `json:"agent_id" binding:"required"`
	DiscoveryType string         `json:"discovery_type" binding:"required"`
	Payload       datatypes.JSON `json:"payload" binding:"required"`
}
