package hosts

type Metadata struct {
	Cluster        string `mapstructure:"ha-cluster,omitempty"`
	ClusterId      string `mapstructure:"ha-cluster-id,omitempty"`
	SAPSystems     string `mapstructure:"sap-systems,omitempty"`
	SAPSystemsId   string `mapstructure:"sap-systems-id,omitempty"`
	SAPSystemsType string `mapstructure:"sap-systems-type,omitempty"`
	CloudProvider  string `mapstructure:"cloud-provider,omitempty"`
	AgentVersion   string `mapstructure:"agent-version,omitempty"`
}
