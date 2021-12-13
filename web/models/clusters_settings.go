package models

type ClusterSettings struct {
	ID             string            `json:"id"`
	SelectedChecks []string          `json:"selected_checks"`
	Hosts          []*HostConnection `json:"hosts"`
}

type HostConnection struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	User    string `json:"user"`
}

type ClustersSettings []*ClusterSettings
