package models

import "time"

type HostTelemetry struct {
	AgentID       string    `gorm:"column:agent_id; primaryKey"`
	HostName      string    `gorm:"column:host_name"`
	SLESVersion   string    `gorm:"column:sles_version"`
	CPUCount      int       `gorm:"column:cpu_count"`
	SocketCount   int       `gorm:"column:socket_count"`
	TotalMemoryMB int       `gorm:"column:total_memory_mb"`
	CloudProvider string    `gorm:"column:cloud_provider"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (HostTelemetry) TableName() string {
	return "host_telemetry"
}
