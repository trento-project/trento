package models

import (
	"fmt"
	"strings"
)

type Check struct {
	ID             string `json:"id,omitempty" mapstructure:"id,omitempty"`
	Name           string `json:"name,omitempty" mapstructure:"name,omitempty"`
	Group          string `json:"group,omitempty" mapstructure:"group,omitempty"`
	Description    string `json:"description,omitempty" mapstructure:"description,omitempty"`
	Remediation    string `json:"remediation,omitempty" mapstructure:"remediation,omitempty"`
	Implementation string `json:"implementation,omitempty" mapstructure:"implementation,omitempty"`
	Labels         string `json:"labels,omitempty" mapstructure:"labels,omitempty"`
	Selected       bool   `json:"selected,omitempty" mapstructure:"selected,omitempty"`
	Result         bool   `json:"result,omitempty" mapstructure:"result,omitempty"`
}

type ChecksResult map[string]ChecksResultByCheck

type ChecksResultByCheck map[string]ChecksResultByHost

type ChecksResultByHost map[string]*Check

func (c *ChecksResultByCheck) GetHostNames() []string {
	var hostNames []string
	for _, rByHost := range (*c) {
		for h, _ := range rByHost {
			hostNames = append(hostNames, h)
		}
		break
	}

	return hostNames
}

func (c *Check) NormalizeID() string {
	return strings.Replace(c.ID, ".", "-", -1)
}

func (c *Check) ExtendedGroupName() string {
	item := strings.Split(c.ID, ".")
	return fmt.Sprintf("%s.%s - %s", item[0], item[1], c.Group)
}
