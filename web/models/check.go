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
type CheckList []*Check

type GroupedChecks struct {
	Group  string
	Checks CheckList
}

type GroupedCheckList []*GroupedChecks

type Metadata struct {
	Checks CheckList `json:"checks,omitempty" mapstructure:"checks,omitempty"`
}

type Results struct {
	Checks map[string]*ChecksByHost `json:"checks,omitempty" mapstructure:"checks,omitempty"`
}

// The ChecksByHost struct stores the checks list, but the results are grouped by hosts
type ChecksByHost struct {
	Hosts map[string]*Check `json:"hosts,omitempty" mapstructure:"hosts,omitempty"`
}

type ClusterCheckResults struct {
	ID          string            `json:"id,omitempty" mapstructure:"id,omitempty"`
	Hosts       map[string]*Check `json:"hosts,omitempty" mapstructure:"hosts,omitempty"`
	Group       string            `json:"group,omitempty" mapstructure:"group,omitempty"`
	Description string            `json:"description,omitempty" mapstructure:"description,omitempty"`
}

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

func (c *Results) GetHostNames() []string {
	var hostNames []string
	for _, cList := range c.Checks {
		for host, _ := range cList.Hosts {
			hostNames = append(hostNames, host)
		}
		break
	}
	return hostNames
}

func (c *Results) HostResultPresent(host string) bool {
	hostList := c.GetHostNames()
	for _, v := range hostList {
		if v == host {
			return true
		}
	}

	return false
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
