package models

const (
	HostHealthPassing  = "passing"
	HostHealthWarning  = "warning"
	HostHealthCritical = "critical"
	HostHealthUnknown  = ""
)

type Host struct {
	ID            string
	Name          string
	Health        string
	IPAddresses   []string
	CloudProvider string
	ClusterID     string
	ClusterName   string
	SAPSystems    []*SAPSystem
	AgentVersion  string
	Tags          []string
	CloudData     interface{}
}

type AzureCloudData struct {
	VMName          string `json:"vmname"`
	ResourceGroup   string `json:"resourceGroup"`
	Location        string `json:"location"`
	VMSize          string `json:"vmsize"`
	DataDisksNumber int    `json:"dataDisksNumber"`
	Offer           string `json:"offer"`
	SKU             string `json:"sku"`
	AdminUsername   string `json:"adminUsername"`
}

type HostList []*Host
