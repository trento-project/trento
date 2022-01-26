package entities

import (
	"sort"
	"time"

	"github.com/lib/pq"
	"github.com/trento-project/trento/web/models"
)

type SAPSystemInstance struct {
	ID                      string `gorm:"primaryKey"`
	AgentID                 string `gorm:"primaryKey"`
	Type                    string
	SID                     string `gorm:"column:sid"`
	InstanceNumber          string `gorm:"primaryKey"`
	Features                string
	Description             string
	StartPriority           string
	Status                  string
	SAPHostname             string
	HttpPort                int
	HttpsPort               int
	SystemReplication       string
	SystemReplicationStatus string
	DBHost                  string
	DBName                  string
	Tenants                 pq.StringArray `gorm:"type:text[]"`
	Host                    *Host          `gorm:"foreignKey:AgentID"`
	UpdatedAt               time.Time
	Tags                    []*models.Tag `gorm:"foreignKey:ResourceID"`
}

type SAPSystemInstances []*SAPSystemInstance

func (s SAPSystemInstances) ToModel() []*models.SAPSystem {
	set := make(map[string]*models.SAPSystem)

	for _, i := range s {

		sapSystem, ok := set[i.ID]
		if !ok {
			// TODO: move to Tags entity when we will have it
			var tags []string
			for _, tag := range i.Tags {
				tags = append(tags, tag.Value)
			}

			sapSystem = &models.SAPSystem{
				ID:     i.ID,
				Type:   i.Type,
				SID:    i.SID,
				DBName: i.DBName,
				DBHost: i.DBHost,
				Tags:   tags,
			}
			set[i.ID] = sapSystem
		}

		sapSystemInstance := &models.SAPSystemInstance{
			InstanceNumber:          i.InstanceNumber,
			Features:                i.Features,
			SystemReplication:       i.SystemReplication,
			SystemReplicationStatus: i.SystemReplicationStatus,
			SAPHostname:             i.SAPHostname,
			Status:                  i.Status,
			StartPriority:           i.StartPriority,
			HttpPort:                i.HttpPort,
			HttpsPort:               i.HttpsPort,
			Type:                    i.Type,
			SID:                     i.SID,
		}

		if i.Host != nil {
			sapSystemInstance.ClusterName = i.Host.ClusterName
			sapSystemInstance.ClusterID = i.Host.ClusterID
			sapSystemInstance.ClusterType = i.Host.ClusterType
			sapSystemInstance.HostID = i.Host.AgentID
			sapSystemInstance.Hostname = i.Host.Name
		}

		sapSystem.Instances = append(sapSystem.Instances, sapSystemInstance)
	}

	var sapSystems []*models.SAPSystem
	for _, sapSystem := range set {
		sapSystems = append(sapSystems, sapSystem)
	}
	sortBySID(sapSystems)

	return sapSystems
}

func sortBySID(sapSystems []*models.SAPSystem) {
	sort.Slice(sapSystems, func(i, j int) bool {
		return sapSystems[i].SID < sapSystems[j].SID
	})
}
