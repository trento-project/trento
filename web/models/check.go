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
}

func (c *Check) NormalizeID() string {
	return strings.Replace(c.ID, ".", "-", -1)
}

func (c *Check) NormalizeGroup() string {
	item := strings.Split(c.ID, ".")
	return fmt.Sprintf("%s.%s - %s", item[0], item[1], c.Group)
}
