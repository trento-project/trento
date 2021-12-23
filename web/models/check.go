package models

import (
	"sort"
)

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
	Premium        bool   `json:"premium" mapstructure:"premium"`
	Selected       bool   `json:"selected,omitempty" mapstructure:"selected,omitempty"`
	Result         string `json:"result,omitempty" mapstructure:"result,omitempty"`
	Msg            string `json:"msg,omitempty" mapstructure:"msg,omitempty"`
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
