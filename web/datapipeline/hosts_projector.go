package datapipeline

import (
	"encoding/json"
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/trento-project/trento/internal/cloud"
	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/web/entities"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewHostsProjector(db *gorm.DB) *projector {
	hostsProjector := NewProjector("hosts", db)

	hostsProjector.AddHandler(HostDiscovery, hostsProjector_HostDiscoveryHandler)
	hostsProjector.AddHandler(CloudDiscovery, hostsProjector_CloudDiscoveryHandler)
	hostsProjector.AddHandler(ClusterDiscovery, hostsProjector_ClusterDiscoveryHandler)

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
		SSHAddress:   discoveredHost.SSHAddress,
		Name:         discoveredHost.HostName,
		IPAddresses:  filterIPAddresses(discoveredHost.HostIpAddresses),
		AgentVersion: discoveredHost.AgentVersion,
	}

	return storeHost(db, host,
		"name",
		"ip_addresses",
		"agent_version",
		"ssh_address",
	)
}

func hostsProjector_CloudDiscoveryHandler(dataCollectedEvent *DataCollectedEvent, db *gorm.DB) error {
	decoder := getPayloadDecoder(dataCollectedEvent.Payload)

	var discoveredCloud cloud.CloudInstance
	if err := decoder.Decode(&discoveredCloud); err != nil {
		log.Errorf("can't decode data: %s", err)
		return err
	}

	parsedCloudData := parseCloudData(discoveredCloud.Provider, discoveredCloud.Metadata)
	jsonCloudData, err := json.Marshal(parsedCloudData)
	if err != nil {
		log.Errorf("can't decode cloud data: %s", err)
		return err
	}

	host := entities.Host{
		AgentID:       dataCollectedEvent.AgentID,
		CloudProvider: discoveredCloud.Provider,
		CloudData:     (datatypes.JSON)(jsonCloudData),
	}

	return storeHost(db, host, "cloud_provider", "cloud_data")
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
		ClusterType: detectClusterType(&discoveredCluster),
	}

	return storeHost(db, host, "cluster_id", "cluster_name", "cluster_type")
}

func storeHost(db *gorm.DB, host entities.Host, updateColumns ...string) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "agent_id"},
		},
		DoUpdates: clause.AssignmentColumns(append(updateColumns, "updated_at")),
	}).Create(&host).Error
}

// filterIPAddresses filters out non-IPv4, loopback or invalid IP addresses
func filterIPAddresses(ipAddresses []string) []string {
	var filtered []string
	for _, ipAddress := range ipAddresses {
		ip := net.ParseIP(ipAddress)
		if ip == nil || ip.IsLoopback() || ip.To4() == nil {
			continue
		}

		filtered = append(filtered, ipAddress)
	}
	return filtered
}

func parseCloudData(provider string, metadata interface{}) *entities.AzureCloudData {
	switch provider {
	case "azure":
		cloudData := parseAzureCloudData(metadata)
		return &cloudData
	default:
		return nil
	}
}

func parseAzureCloudData(metadata interface{}) entities.AzureCloudData {
	var cloudData entities.AzureCloudData
	metadataMap, ok := metadata.(map[string]interface{})
	if !ok {
		log.Errorf("can't decode Azure metadata")
		return cloudData
	}

	computeInfo, ok := metadataMap["compute"].(map[string]interface{})
	if !ok {
		log.Errorf("can't decode compute metadata")
		return cloudData
	}

	if _, ok := computeInfo["name"]; ok {
		cloudData.VMName = computeInfo["name"].(string)
	}

	if _, ok := computeInfo["name"]; ok {
		cloudData.Offer = computeInfo["offer"].(string)
	}

	if _, ok := computeInfo["resourceGroupName"]; ok {
		cloudData.ResourceGroup = computeInfo["resourceGroupName"].(string)
	}

	if _, ok := computeInfo["sku"]; ok {
		cloudData.SKU = computeInfo["sku"].(string)
	}

	if _, ok := computeInfo["vmSize"]; ok {
		cloudData.VMSize = computeInfo["vmSize"].(string)
	}

	if _, ok := computeInfo["location"]; ok {
		cloudData.Location = computeInfo["location"].(string)
	}

	osProfile, ok := computeInfo["osProfile"].(map[string]interface{})
	if !ok {
		log.Errorf("can't decode os profile")
		return cloudData
	}

	if _, ok := osProfile["adminUsername"]; ok {
		cloudData.AdminUsername = osProfile["adminUsername"].(string)
	}

	storageProfile, ok := computeInfo["storageProfile"].(map[string]interface{})
	if !ok {
		log.Errorf("can't decode storage profile")
		return cloudData
	}

	if _, ok := storageProfile["dataDisks"]; ok {
		dataDisks, ok := storageProfile["dataDisks"].([]interface{})
		if !ok {
			log.Errorf("can't decode Azure data disks")
			return cloudData
		}
		cloudData.DataDisksNumber = len(dataDisks)
	}

	return cloudData
}
