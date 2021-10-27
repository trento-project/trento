package datapipeline

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/models"
)

func TestClusterListHandler(t *testing.T) {
	db, err := helpers.SetupTestDatabase()
	if err != nil {
		t.Skip(err)
	}

	tx := db.Begin()
	defer tx.Rollback()

	tx.AutoMigrate(&models.Cluster{})
	tx.Create(&models.Cluster{
		Name:        "test_cluster",
		ID:          "test_id",
		ClusterType: models.ClusterTypeUnknown,
	})

	jsonFile, err := os.Open("../../test/fixtures/cluster_discovery_hana_scale_up.json")
	if err != nil {
		panic(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	dataCollectedEvent := &DataCollectedEvent{
		ID:            1,
		AgentID:       "agent_id",
		DiscoveryType: ClusterDiscovery,
		Payload:       byteValue,
	}

	ClusterListHandler(dataCollectedEvent, tx)

	var cluster models.Cluster
	tx.First(&cluster)

	assert.EqualValues(t,
		models.Cluster{
			Name:            "test_cluster",
			ID:              "test_id",
			ClusterType:     models.ClusterTypeScaleUp,
			SIDs:            []string{"DEV"},
			ResourcesNumber: 5,
			HostsNumber:     3,
		}, cluster)
}

func TestTransformClusterListData_HANAScaleUp(t *testing.T) {
	jsonFile, err := os.Open("../../test/fixtures/cluster_discovery_hana_scale_up.json")
	if err != nil {
		panic(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var clusterIn cluster.Cluster
	json.Unmarshal(byteValue, &clusterIn)
	clusterOut, _ := transformClusterData(&clusterIn)

	assert.EqualValues(t,
		&models.Cluster{
			Name:            "test_cluster",
			ID:              "test_id",
			ClusterType:     models.ClusterTypeScaleUp,
			SIDs:            []string{"DEV"},
			ResourcesNumber: 5,
			HostsNumber:     3,
		}, clusterOut)
}

func TestTransformClusterListData_Unknown(t *testing.T) {
	jsonFile, err := os.Open("../../test/fixtures/cluster_discovery_unknown.json")
	if err != nil {
		panic(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var clusterIn cluster.Cluster
	json.Unmarshal(byteValue, &clusterIn)
	clusterOut, _ := transformClusterData(&clusterIn)

	assert.EqualValues(t,
		&models.Cluster{
			Name:        "test_cluster",
			ID:          "test_id",
			ClusterType: models.ClusterTypeUnknown,
		}, clusterOut)
}
