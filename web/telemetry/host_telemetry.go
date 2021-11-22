package telemetry

import (
	"errors"
	"time"

	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
)

type HostTelemetry struct {
	InstallationID string    `json:"installation_id"`
	AgentID        string    `json:"agent_id"`
	SLESVersion    string    `json:"sles_version"`
	CPUCount       int       `json:"cpu_count"`
	SocketCount    int       `json:"socket_count"`
	TotalMemoryMB  int       `json:"total_memory_mb"`
	CloudProvider  string    `json:"cloud_provider"`
	Time           time.Time `json:"time"`
}

type HostTelemetries []*HostTelemetry

type HostTelemetryExtractor struct {
	installationIdAwareExtractor
	db *gorm.DB
}

func (ex *HostTelemetryExtractor) Extract() (interface{}, error) {
	var collectedHostsTelemetry []models.HostTelemetry
	var publishableHostTelemetries HostTelemetries

	if err := ex.db.Find(&collectedHostsTelemetry).Error; err != nil {
		return nil, err
	}

	if len(collectedHostsTelemetry) == 0 {
		return nil, errors.New("no host telemetry found")
	}

	for _, hostTelemetry := range collectedHostsTelemetry {
		publishableHostTelemetries = append(publishableHostTelemetries, &HostTelemetry{
			InstallationID: ex.installationID.String(),
			AgentID:        hostTelemetry.AgentID,
			SLESVersion:    hostTelemetry.SLESVersion,
			CPUCount:       hostTelemetry.CPUCount,
			SocketCount:    hostTelemetry.SocketCount,
			TotalMemoryMB:  hostTelemetry.TotalMemoryMB,
			CloudProvider:  hostTelemetry.CloudProvider,
			Time:           hostTelemetry.UpdatedAt,
		})
	}

	return publishableHostTelemetries, nil
}

func NewHostTelemetryExtractor(db *gorm.DB) Extractor {
	return &HostTelemetryExtractor{
		db: db,
	}
}
