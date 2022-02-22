package entities

import (
	"time"

	"github.com/trento-project/trento/web/models"
	"gorm.io/datatypes"
)

type Cluster struct {
	ID              string `gorm:"primaryKey"`
	Name            string
	ClusterType     string
	SID             string `gorm:"column:sid"`
	ResourcesNumber int
	HostsNumber     int
	Health          *HealthState  `gorm:"foreignkey:id"`
	Tags            []*models.Tag `gorm:"polymorphic:Resource;polymorphicValue:clusters"`
	UpdatedAt       time.Time
	Hosts           []*Host        `gorm:"foreignkey:cluster_id"`
	Details         datatypes.JSON `json:"payload" binding:"required"`
}

type HANAClusterDetails struct {
	SystemReplicationMode          string             `json:"system_replication_mode"`
	SystemReplicationOperationMode string             `json:"system_replication_operation_mode"`
	SecondarySyncState             string             `json:"secondary_sync_state"`
	SRHealthState                  string             `json:"sr_health_state"`
	CIBLastWritten                 time.Time          `json:"cib_last_written"`
	FencingType                    string             `json:"fencing_type"`
	StoppedResources               []*ClusterResource `json:"stopped_resources"`
	Nodes                          []*HANAClusterNode `json:"nodes"`
	SBDDevices                     []*SBDDevice       `json:"sbd_devices"`
}

type ClusterResource struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	FailCount int    `json:"fail_count"`
}

type HANAClusterNode struct {
	Name       string             `json:"name"`
	Site       string             `json:"site"`
	Attributes map[string]string  `json:"attributes"`
	Resources  []*ClusterResource `json:"resources"`
	VirtualIPs []string           `json:"virtual_ips"`
	HANAStatus string             `json:"hana_status"`
}

type SBDDevice struct {
	Device string `json:"device"`
	Status string `json:"status"`
}

func (c *Cluster) ToModel() *models.Cluster {
	// TODO: move to Tags entity when we will have it
	var tags []string
	for _, tag := range c.Tags {
		tags = append(tags, tag.Value)
	}

	health := models.HealthSummaryHealthUnknown
	if c.Health != nil {
		health = c.Health.Health
	}

	return &models.Cluster{
		ID:              c.ID,
		Name:            c.Name,
		ClusterType:     c.ClusterType,
		SID:             c.SID,
		ResourcesNumber: c.ResourcesNumber,
		HostsNumber:     c.HostsNumber,
		Health:          health,
		Tags:            tags,
	}
}

func (h *HANAClusterDetails) ToModel() *models.HANAClusterDetails {
	var stoppedResources []*models.ClusterResource
	for _, r := range h.StoppedResources {
		stoppedResources = append(stoppedResources, r.ToModel())
	}

	var nodes []*models.HANAClusterNode
	for _, n := range h.Nodes {
		nodes = append(nodes, n.ToModel())
	}

	var sbdDevices []*models.SBDDevice
	for _, s := range h.SBDDevices {
		sbdDevices = append(sbdDevices, s.ToModel())
	}

	return &models.HANAClusterDetails{
		SystemReplicationMode:          h.SystemReplicationMode,
		SystemReplicationOperationMode: h.SystemReplicationOperationMode,
		SecondarySyncState:             h.SecondarySyncState,
		SRHealthState:                  h.SRHealthState,
		CIBLastWritten:                 h.CIBLastWritten,
		FencingType:                    h.FencingType,
		StoppedResources:               stoppedResources,
		Nodes:                          nodes,
		SBDDevices:                     sbdDevices,
	}
}

func (r *ClusterResource) ToModel() *models.ClusterResource {
	return &models.ClusterResource{
		ID:        r.ID,
		Type:      r.Type,
		Role:      r.Role,
		Status:    r.Status,
		FailCount: r.FailCount,
	}
}

func (s *SBDDevice) ToModel() *models.SBDDevice {
	return &models.SBDDevice{
		Device: s.Device,
		Status: s.Status,
	}
}

func (n *HANAClusterNode) ToModel() *models.HANAClusterNode {
	var resources []*models.ClusterResource
	for _, r := range n.Resources {
		resources = append(resources, r.ToModel())
	}

	return &models.HANAClusterNode{
		Name:       n.Name,
		Site:       n.Site,
		Attributes: n.Attributes,
		Resources:  resources,
		VirtualIPs: n.VirtualIPs,
		HANAStatus: n.HANAStatus,
	}
}
