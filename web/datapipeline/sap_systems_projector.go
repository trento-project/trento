package datapipeline

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/internal/sapsystem"
	"github.com/trento-project/trento/web/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewSAPSystemsProjector(db *gorm.DB) *projector {
	SAPSystemsProjector := NewProjector("sapsystems", db)

	SAPSystemsProjector.AddHandler(SAPsystemDiscovery, SAPSystemsProjector_SAPSystemsDiscoveryHandler)

	return SAPSystemsProjector
}

func SAPSystemsProjector_SAPSystemsDiscoveryHandler(dataCollectedEvent *DataCollectedEvent, db *gorm.DB) error {
	decoder := getPayloadDecoder(dataCollectedEvent.Payload)
	var discoveredSAPSystems sapsystem.SAPSystemsList
	if err := decoder.Decode(&discoveredSAPSystems); err != nil {
		log.Errorf("can't decode data: %s", err)
		return err
	}

	for _, s := range discoveredSAPSystems {
		instance := entities.SAPSystemInstance{
			AgentID: dataCollectedEvent.AgentID,
			ID:      s.Id,
			SID:     s.SID,
		}

		switch s.Type {
		case 1:
			instance.Type = sapsystem.SAPSystemsDatabase
		case 2:
			instance.Type = sapsystem.SAPSystemsApplication
			instance.DBHost = fmt.Sprint(s.Profile["SAPDBHOST"])
		}

		for _, i := range s.Instances {
			var features string
			var instanceNumber string

			if p, ok := i.SAPControl.Properties["SAPSYSTEM"]; ok {
				instanceNumber = p.Value
			} else {
				log.Warnf("Instance Number not found for %s", s.SID)
				continue
			}

			for _, SAPControlInstance := range i.SAPControl.Instances {
				if instanceNumber == fmt.Sprintf("%02d", SAPControlInstance.InstanceNr) {
					features = SAPControlInstance.Features
					break
				}
			}

			instance.Features = features
			instance.InstanceNumber = instanceNumber

			err := storeSAPInstance(db, instance, "id", "sid", "type", "features", "instance_number")
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func storeSAPInstance(db *gorm.DB, sapInstance entities.SAPSystemInstance, updateColumns ...string) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "agent_id"},
			{Name: "id"},
			{Name: "instance_number"},
		},
		DoUpdates: clause.AssignmentColumns(append(updateColumns, "updated_at")),
	}).Create(&sapInstance).Error
}
