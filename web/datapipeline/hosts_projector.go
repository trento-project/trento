package datapipeline

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/trento-project/trento/internal/cloud"
	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/internal/sapsystem"
	"github.com/trento-project/trento/web/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewHostsProjector(db *gorm.DB) *projector {
	hostsProjector := NewProjector("hosts", db)

	hostsProjector.AddHandler(HostDiscovery, hostsProjector_HostDiscoveryHandler)
	hostsProjector.AddHandler(CloudDiscovery, hostsProjector_CloudDiscoveryHandler)
	hostsProjector.AddHandler(ClusterDiscovery, hostsProjector_ClusterDiscoveryHandler)
	hostsProjector.AddHandler(SAPsystemDiscovery, hostsProjector_SAPSystemDiscoveryHandler)

	return hostsProjector
}

func hostsProjector_HostDiscoveryHandler(dataCollectedEvent *DataCollectedEvent, db *gorm.DB) error {
	decoder := getPayloadDecoder(dataCollectedEvent.Payload)

	var discoveredHost hosts.DiscoveredHost
	if err := decoder.Decode(&discoveredHost); err != nil {
		log.Errorf("can't decode data: %s", err)
		return err
	}

	host := entities.Host{
		AgentID:      dataCollectedEvent.AgentID,
		Name:         discoveredHost.HostName,
		IPAddresses:  discoveredHost.HostIpAddresses,
		AgentVersion: discoveredHost.AgentVersion,
	}

	return storeHost(db, host,
		"name",
		"ip_addresses",
		"agent_version",
	)
}

func hostsProjector_CloudDiscoveryHandler(dataCollectedEvent *DataCollectedEvent, db *gorm.DB) error {
	decoder := getPayloadDecoder(dataCollectedEvent.Payload)

	var discoveredCloud cloud.CloudInstance
	if err := decoder.Decode(&discoveredCloud); err != nil {
		log.Errorf("can't decode data: %s", err)
		return err
	}

	host := entities.Host{
		AgentID:       dataCollectedEvent.AgentID,
		CloudProvider: discoveredCloud.Provider,
	}

	return storeHost(db, host, "cloud_provider")
}

func hostsProjector_ClusterDiscoveryHandler(dataCollectedEvent *DataCollectedEvent, db *gorm.DB) error {
	decoder := getPayloadDecoder(dataCollectedEvent.Payload)

	var discoveredCluster cluster.Cluster
	if err := decoder.Decode(&discoveredCluster); err != nil {
		log.Errorf("can't decode data: %s", err)
		return err
	}

	host := entities.Host{
		AgentID:     dataCollectedEvent.AgentID,
		ClusterID:   discoveredCluster.Id,
		ClusterName: discoveredCluster.Name,
	}

	return storeHost(db, host, "cluster_id", "cluster_name")
}

func hostsProjector_SAPSystemDiscoveryHandler(dataCollectedEvent *DataCollectedEvent, db *gorm.DB) error {
	decoder := getPayloadDecoder(dataCollectedEvent.Payload)

	var discoveredSAPSystems sapsystem.SAPSystemsList
	if err := decoder.Decode(&discoveredSAPSystems); err != nil {
		log.Errorf("can't decode data: %s", err)
		return err
	}

	var sids []string
	for _, s := range discoveredSAPSystems {
		fmt.Println(s.SID)
		sids = append(sids, s.SID)
	}

	host := entities.Host{
		AgentID: dataCollectedEvent.AgentID,
		SIDs:    sids,
	}

	return storeHost(db, host, "sids")
}

func storeHost(db *gorm.DB, host entities.Host, updateColumns ...string) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "agent_id"},
		},
		DoUpdates: clause.AssignmentColumns(append(updateColumns, "updated_at")),
	}).Create(&host).Error
}
