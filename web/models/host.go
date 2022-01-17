package models

import (
	"github.com/trento-project/trento/internal/cloud"
)

const (
	HostHealthPassing  = "passing"
	HostHealthWarning  = "warning"
	HostHealthCritical = "critical"
	HostHealthUnknown  = ""

	Azure = "Azure"
	Aws   = "AWS"
	Gcp   = "GCP"
)

type Host struct {
	ID            string
	Name          string
	Health        string
	IPAddresses   []string
	CloudProvider string
	ClusterID     string
	ClusterName   string
	ClusterType   string
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

func (h *Host) PrettyProvider() string {
	switch h.CloudProvider {
	case cloud.Azure:
		return Azure
	case cloud.Aws:
		return Aws
	case cloud.Gcp:
		return Gcp
	default:
		return ""
	}
}
