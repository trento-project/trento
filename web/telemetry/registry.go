package telemetry

import "gorm.io/gorm"

func NewTelemetryRegistry(db *gorm.DB) *TelemetryRegistry {
	return &TelemetryRegistry{
		"host_telemetry": NewHostTelemetryExtractor(db),
		//"cluster": NewClusterTelemetryExtractor(),
	}
}
