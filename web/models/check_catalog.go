package models

import (
	"sort"
	"time"

	"gorm.io/datatypes"
)

// Store the data as payload. Changes in this struct are expected
type CheckRaw struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	Payload   datatypes.JSON
}

// TableName overrides the table name used by CheckRaw to `checks`
func (CheckRaw) TableName() string {
	return "checks"
}

// List is used instead of a map as it guarantees order
type ChecksCatalog []*Check

type Check struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Group          string `json:"group,omitempty"`
	Description    string `json:"description,omitempty"`
	Remediation    string `json:"remediation,omitempty"`
	Implementation string `json:"implementation,omitempty"`
	Labels         string `json:"labels,omitempty"`
	Selected       bool   `json:"selected,omitempty"`
	Result         string `json:"result,omitempty"`
}

type GroupedChecks struct {
	Group  string
	Checks ChecksCatalog
}

type GroupedCheckList []*GroupedChecks

// Sorting methods for GroupedCheckList

func (g GroupedCheckList) Len() int {
	return len(g)
}

func (g GroupedCheckList) Less(i, j int) bool {
	return g[i].Checks[0].Name < g[j].Checks[0].Name
}

func (g GroupedCheckList) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

func (g GroupedCheckList) OrderByName() GroupedCheckList {
	sort.Sort(g)
	return g
}
