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
	ResourceGroup   string `json:"resource_group"`
	Location        string `json:"location"`
	VMSize          string `json:"vmsize"`
	DataDisksNumber int    `json:"data_disks_number"`
	Offer           string `json:"offer"`
	SKU             string `json:"sku"`
	AdminUsername   string `json:"admin_username"`
}

type HostList []*Host
