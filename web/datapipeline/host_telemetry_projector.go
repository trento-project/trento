package datapipeline

import (
	"bytes"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	discoveryModels "github.com/trento-project/trento/agent/discovery/models"
	"github.com/trento-project/trento/internal/cloud"
	"github.com/trento-project/trento/internal/subscription"
	"github.com/trento-project/trento/web/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewHostTelemetryProjector(db *gorm.DB) *projector {
	telemetryProjector := NewProjector("host_telemetry", db)

	telemetryProjector.AddHandler(HostDiscovery, hostTelemetryProjector_HostDiscoveryHandler)
	telemetryProjector.AddHandler(CloudDiscovery, hostTelemetryProjector_CloudDiscoveryHandler)
	telemetryProjector.AddHandler(SubscriptionDiscovery, hostTelemetryProjector_SubscriptionDiscoveryHandler)

	return telemetryProjector
}

func hostTelemetryProjector_HostDiscoveryHandler(dataCollectedEvent *DataCollectedEvent, db *gorm.DB) error {
	decoder := payloadDecoder(dataCollectedEvent.Payload)

	var discoveredHost discoveryModels.DiscoveredHost
	if err := decoder.Decode(&discoveredHost); err != nil {
		log.Errorf("can't decode data: %s", err)
		return err
	}

	telemetryReadModel := models.HostTelemetry{
		AgentID:       dataCollectedEvent.AgentID,
		HostName:      discoveredHost.HostName,
		CPUCount:      discoveredHost.CPUCount,
		SocketCount:   discoveredHost.SocketCount,
		TotalMemoryMB: discoveredHost.TotalMemoryMB,
	}

	return storeHostTelemetry(db, telemetryReadModel,
		"host_name",
		"cpu_count",
		"socket_count",
		"total_memory_mb",
	)
}

func hostTelemetryProjector_CloudDiscoveryHandler(dataCollectedEvent *DataCollectedEvent, db *gorm.DB) error {
	decoder := payloadDecoder(dataCollectedEvent.Payload)

	var discoveredCloud cloud.CloudInstance
	if err := decoder.Decode(&discoveredCloud); err != nil {
		log.Errorf("can't decode data: %s", err)
		return err
	}

	telemetryReadModel := models.HostTelemetry{
		AgentID:       dataCollectedEvent.AgentID,
		CloudProvider: discoveredCloud.Provider,
	}

	return storeHostTelemetry(db, telemetryReadModel, "cloud_provider")
}

func hostTelemetryProjector_SubscriptionDiscoveryHandler(dataCollectedEvent *DataCollectedEvent, db *gorm.DB) error {
	decoder := payloadDecoder(dataCollectedEvent.Payload)

	var discoveredSubscription subscription.Subscriptions

	if err := decoder.Decode(&discoveredSubscription); err != nil {
		log.Errorf("can't decode data: %s", err)
		return err
	}

	if len(discoveredSubscription) == 0 {
		return nil
	}

	telemetryReadModel := models.HostTelemetry{
		AgentID:     dataCollectedEvent.AgentID,
		SLESVersion: discoveredSubscription[0].Version,
	}

	return storeHostTelemetry(db, telemetryReadModel, "sles_version")
}

func payloadDecoder(payload datatypes.JSON) *json.Decoder {
	data, _ := payload.MarshalJSON()
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	return decoder
}

func storeHostTelemetry(db *gorm.DB, telemetryReadModel models.HostTelemetry, updateColumns ...string) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "agent_id"},
		},
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(&telemetryReadModel).Error
}
