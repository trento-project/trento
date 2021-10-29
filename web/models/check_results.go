package models

import (
	"time"

	"gorm.io/datatypes"
)

const (
	CheckPassing   string = "passing"
	CheckWarning   string = "warning"
	CheckCritical  string = "critical"
	CheckSkipped   string = "skipped"
	CheckUndefined string = "undefined"
)

type CheckResultsRaw struct {
	ID        int64
	CreatedAt time.Time
	GroupID   string
	Payload   datatypes.JSON
}

// TableName overrides the table name used by CheckResultsRaw to `check_results`
func (CheckResultsRaw) TableName() string {
	return "check_results"
}

type Results struct {
	Hosts  map[string]*Host         `json:"hosts,omitempty"`
	Checks map[string]*ChecksByHost `json:"checks,omitempty"`
}

// Simplifed stuct consumed by the frontend
type ResultsAsList struct {
	Hosts  map[string]*Host `json:"hosts,omitempty"`
	Checks []*ChecksByHost  `json:"checks,omitempty"`
}

// The ChecksByHost struct stores the checks list, but the results are grouped by hosts
type ChecksByHost struct {
	Hosts       map[string]*Check `json:"hosts,omitempty"`
	ID          string            `json:"id,omitempty"`
	Group       string            `json:"group,omitempty"`
	Description string            `json:"description,omitempty"`
}

type Host struct {
	Reachable bool   `json:"reachable"`
	Msg       string `json:"msg"`
}
