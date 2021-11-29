package entities

import (
	"time"

	"github.com/trento-project/trento/web/models"
)

type SAPSystemInstance struct {
	ID                string `gorm:"primaryKey"`
	AgentID           string `gorm:"primaryKey"`
	Type              string
	SID               string `gorm:"column:sid"`
	InstanceNumber    string `gorm:"primaryKey"`
	Features          string
	Description       string
	SystemReplication string
	DBHost            string
	Host              *Host `gorm:"foreignKey:AgentID"`
	UpdatedAt         time.Time
}

type SAPSystemInstances []*SAPSystemInstance

func (s SAPSystemInstances) ToModel() []*models.SAPSystem {
	set := make(map[string]*models.SAPSystem)

	for _, i := range s {
		sapSystem, ok := set[i.ID]
		if !ok {
			sapSystem = &models.SAPSystem{
				ID:   i.ID,
				Type: i.Type,
				SID:  i.SID,
			}
			set[i.ID] = sapSystem
		}
		sapSystem.Instances = append(sapSystem.Instances,
			&models.SAPSystemInstance{
				InstanceNumber: i.InstanceNumber,
				Features:       i.Features,
			})
	}

	var sapSystems []*models.SAPSystem
	for _, sapSystem := range set {
		sapSystems = append(sapSystems, sapSystem)
	}
	return sapSystems
}
