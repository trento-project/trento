package entities

import (
	"gorm.io/datatypes"
)

type HealthState struct {
	ID             string `gorm:"primaryKey"`
	Health         string
	PartialHealths datatypes.JSON
}

// PartialHealths is something like
// {"config_checks": "passing", "hana_sr_health": "passing"}

// The Health and PartialHealths are changed upon events:
// config_checks when we receive new check results
// hana_sr_health when we discover new cluster data
