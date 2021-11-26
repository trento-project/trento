package entities

import (
	"time"
)

type SAPSystemInstance struct {
	ID                string `gorm:"primaryKey"`
	AgentID           string `gorm:"primaryKey"`
	Type              string
	SID               string `gorm:"column:sid"`
	InstanceNumber    string `gorm:"primaryKey"`
	Features          string
	Description       string
	SystemReplication string
	DBHost            string
	Host              *Host `gorm:"foreignKey:AgentID"`
	UpdatedAt         time.Time
}
