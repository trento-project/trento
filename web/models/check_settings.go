package models

import (
	"github.com/lib/pq"
)

type SelectedChecks struct {
	ID             string         `gorm:"primaryKey"`
	SelectedChecks pq.StringArray `gorm:"type:text[]"`
}

type ConnectionSettings struct {
	ID   string `gorm:"primaryKey"`
	Node string `gorm:"primaryKey"`
	User string
}
