package models

type Check struct {
	ID             string `json:"id,omitempty" mapstructure:"id,omitempty"`
	Name           string `json:"name,omitempty" mapstructure:"name,omitempty"`
	Description    string `json:"description,omitempty" mapstructure:"description,omitempty"`
	Remediation    string `json:"remediation,omitempty" mapstructure:"remediation,omitempty"`
	Implementation string `json:"implementation,omitempty" mapstructure:"implementation,omitempty"`
	Labels         string `json:"labels,omitempty" mapstructure:"labels,omitempty"`
}
