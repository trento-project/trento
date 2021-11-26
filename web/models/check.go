package models

import (
	"sort"
)

const (
	CheckPassing   string = "passing"
	CheckWarning   string = "warning"
	CheckCritical  string = "critical"
	CheckSkipped   string = "skipped"
	CheckUndefined string = "undefined"
)

type CheckData struct {
	Metadata Metadata            `json:"metadata,omitempty" mapstructure:"metadata,omitempty"`
	Groups   map[string]*Results `json:"results,omitempty" mapstructure:"results,omitempty"`
}

// List is used instead of a map as it guarantees order
type ChecksCatalog []*Check

type Check struct {
	ID             string `json:"id,omitempty" mapstructure:"id,omitempty"`
	Name           string `json:"name,omitempty" mapstructure:"name,omitempty"`
	Group          string `json:"group,omitempty" mapstructure:"group,omitempty"`
	Description    string `json:"description,omitempty" mapstructure:"description,omitempty"`
	Remediation    string `json:"remediation,omitempty" mapstructure:"remediation,omitempty"`
	Implementation string `json:"implementation,omitempty" mapstructure:"implementation,omitempty"`
	Labels         string `json:"labels,omitempty" mapstructure:"labels,omitempty"`
	Selected       bool   `json:"selected,omitempty" mapstructure:"selected,omitempty"`
	Result         string `json:"result,omitempty" mapstructure:"result,omitempty"`
}

type GroupedChecks struct {
	Group  string
	Checks ChecksCatalog
}

type GroupedCheckList []*GroupedChecks

type Metadata struct {
	Checks ChecksCatalog `json:"checks,omitempty" mapstructure:"checks,omitempty"`
}

type Results struct {
	Hosts  map[string]*CheckHost    `json:"hosts,omitempty" mapstructure:"hosts,omitempty"`
	Checks map[string]*ChecksByHost `json:"checks,omitempty" mapstructure:"checks,omitempty"`
}

// The ChecksByHost struct stores the checks list, but the results are grouped by hosts
type ChecksByHost struct {
	Hosts map[string]*Check `json:"hosts,omitempty" mapstructure:"hosts,omitempty"`
}

type CheckHost struct {
	Reachable bool   `json:"reachable" mapstructure:"reachable"`
	Msg       string `json:"msg" mapstructure:"msg"`
}

// Simplified models for the frontend
type ClusterCheckResults struct {
	Hosts  map[string]*CheckHost `json:"hosts" mapstructure:"hosts,omitempty"`
	Checks []ClusterCheckResult  `json:"checks" mapstructure:"checks,omitempty"`
}

type ClusterCheckResult struct {
	ID          string            `json:"id,omitempty" mapstructure:"id,omitempty"`
	Hosts       map[string]*Check `json:"hosts,omitempty" mapstructure:"hosts,omitempty"`
	Group       string            `json:"group,omitempty" mapstructure:"group,omitempty"`
	Description string            `json:"description,omitempty" mapstructure:"description,omitempty"`
}

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
