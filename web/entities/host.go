package entities

import (
	"time"

	"github.com/lib/pq"
	"github.com/trento-project/trento/web/models"
	"gorm.io/datatypes"
)

type Host struct {
	AgentID            string `gorm:"primaryKey"`
	AgentBindIP        string
	Name               string
	IPAddresses        pq.StringArray `gorm:"type:text[]"`
	CloudProvider      string
	ClusterID          string
	ClusterName        string
	ClusterType        string
	SAPSystemInstances SAPSystemInstances `gorm:"foreignkey:AgentID"`
	AgentVersion       string
	Heartbeat          *HostHeartbeat    `gorm:"foreignKey:AgentID"`
	Subscription       *SlesSubscription `gorm:"foreignKey:AgentID"`
	Tags               []*models.Tag     `gorm:"polymorphic:Resource;polymorphicValue:hosts"`
	UpdatedAt          time.Time
	CloudData          datatypes.JSON
}

type HostHeartbeat struct {
	AgentID   string `gorm:"primaryKey"`
	UpdatedAt time.Time
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

func (h *Host) ToModel() *models.Host {
	// TODO: move to Tags entity when we will have it
	var tags []string
	for _, tag := range h.Tags {
		tags = append(tags, tag.Value)
	}

	return &models.Host{
		ID:            h.AgentID,
		Name:          h.Name,
		IPAddresses:   h.IPAddresses,
		CloudProvider: h.CloudProvider,
		ClusterID:     h.ClusterID,
		ClusterName:   h.ClusterName,
		ClusterType:   h.ClusterType,
		AgentVersion:  h.AgentVersion,
		Tags:          tags,
		SAPSystems:    h.SAPSystemInstances.ToModel(),
	}
}
