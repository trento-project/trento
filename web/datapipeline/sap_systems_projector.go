package datapipeline

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/internal/sapsystem"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
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

	// deletes all obsolete instances if no sap system was discovered
	if len(discoveredSAPSystems) == 0 {
		return db.
			Where("agent_id = ?", dataCollectedEvent.AgentID).
			Delete(&entities.SAPSystemInstance{}).
			Error
	}

	for _, s := range discoveredSAPSystems {
		var sapSystemType, dbHost, dbName, dbAddress string
		var tenants []string

		switch s.Type {
		case 1:
			for _, tenant := range s.Databases {
				tenants = append(tenants, tenant.Database)
			}
			sapSystemType = models.SAPSystemTypeDatabase
		case 2:
			sapSystemType = models.SAPSystemTypeApplication
			dbHost = fmt.Sprint(s.Profile["SAPDBHOST"])

			if hdb, ok := s.Profile["dbs/hdb/dbname"]; ok {
				dbName = hdb.(string)
			}
			dbAddress = s.DBAddress
		case 3:
			log.Infof("SAP diagnostics agent with %s identifier found. Skipping projection", s.SID)
			continue
		}

		var instances []entities.SAPSystemInstance
		for _, i := range s.Instances {
			instance := entities.SAPSystemInstance{
				AgentID:   dataCollectedEvent.AgentID,
				ID:        s.Id,
				SID:       s.SID,
				Type:      sapSystemType,
				Tenants:   tenants,
				DBHost:    dbHost,
				DBName:    dbName,
				DBAddress: dbAddress,
			}

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
			instance.SystemReplication = parseReplicationMode(i.SystemReplication)
			instance.SystemReplicationStatus = parseReplicationStatus(i.SystemReplication)
			addSAPControlData(&instance, i.SAPControl)

			instances = append(instances, instance)
		}

		err := storeSAPInstances(db,
			instances,
			"id", "sid", "type", "features", "instance_number",
			"system_replication", "system_replication_status",
			"sap_hostname", "start_priority", "http_port", "https_port", "status",
			"tenants", "db_host", "db_name", "db_address")
		if err != nil {
			return err
		}

		// delete obsolete instances
		if len(instances) > 0 {
			for _, instance := range instances {
				db = db.Where("NOT (id = ? AND instance_number = ?)", instance.ID, instance.InstanceNumber).
					Where("agent_id = ?", dataCollectedEvent.AgentID)
			}
		} else {
			db = db.Where("id = ? AND agent_id = ?", s.Id, dataCollectedEvent.AgentID)
		}

		err = db.Delete(&entities.SAPSystemInstance{}).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func storeSAPInstances(db *gorm.DB, sapInstances []entities.SAPSystemInstance, updateColumns ...string) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "agent_id"},
			{Name: "id"},
			{Name: "instance_number"},
		},
		DoUpdates: clause.AssignmentColumns(append(updateColumns, "updated_at")),
	}).Create(&sapInstances).Error
}

func parseReplicationMode(r sapsystem.SystemReplication) string {
	localSite, ok := r["local_site_id"]
	if !ok {
		return ""
	}

	replicationKey := fmt.Sprintf("site/%s/REPLICATION_MODE", localSite)
	mode, ok := r[replicationKey]
	if !ok {
		return ""
	}

	switch mode {
	case "PRIMARY":
		return "Primary"
	case "":
		return ""
	default: // SYNC, SYNCMEM, ASYNC, UNKNOWN
		return "Secondary"
	}
}

// Find status information at: https://help.sap.com/viewer/4e9b18c116aa42fc84c7dbfd02111aba/2.0.04/en-US/aefc55a27003440792e34ece2125dc89.html
func parseReplicationStatus(s sapsystem.SystemReplication) string {
	status, ok := s["overall_replication_status"]
	if !ok {
		return ""
	}

	status = fmt.Sprintf("%v", status)

	switch status {
	case "ACTIVE":
		return "SOK"
	case "ERROR":
		return "SFAIL"
	default: // UNKNOWN, INITIALIZING, SYNCING
		return ""
	}
}

func addSAPControlData(instance *entities.SAPSystemInstance, sapControl *sapsystem.SAPControl) {
	for _, i := range sapControl.Instances {
		if instance.InstanceNumber == fmt.Sprintf("%02d", i.InstanceNr) {
			instance.StartPriority = i.StartPriority
			instance.Status = (string)(i.Dispstatus)
			instance.SAPHostname = i.Hostname
			instance.HttpPort = (int)(i.HttpPort)
			instance.HttpsPort = (int)(i.HttpsPort)
		}
	}
}
