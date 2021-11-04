package datapipeline

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	ClusterDiscovery      = "ha_cluster_discovery"
	SAPsystemDiscovery    = "sap_system_discovery"
	HostDiscovery         = "host_discovery"
	SubscriptionDiscovery = "subscription_discovery"
	CloudDiscovery        = "cloud_discovery"
)

type DataCollectedEvent struct {
	ID            int64
	CreatedAt     time.Time
	AgentID       string         `json:"agent_id" binding:"required"`
	DiscoveryType string         `json:"discovery_type" binding:"required"`
	Payload       datatypes.JSON `json:"payload" binding:"required"`
}

func PruneEvents(olderThan time.Duration, db *gorm.DB) error {
	prunedEvents := db.Delete(DataCollectedEvent{}, "created_at < ?", time.Now().Add(-olderThan))
	log.Debugf("Pruned %d events", prunedEvents.RowsAffected)

	return prunedEvents.Error
}
