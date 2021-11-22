package entities

import (
	"time"
)

type SAPSystemInstance struct {
	AgentID           string `gorm:"primaryKey"`
	Type              string
	SystemID          string `gorm:"primaryKey"`
	SID               string `gorm:"column:sid"`
	InstanceNumber    string `gorm:"primaryKey"`
	Features          string
	Description       string
	SystemReplication string
	DBHost            string
	Host              *Host `gorm:"foreignKey:AgentID"`
	UpdatedAt         time.Time
}
