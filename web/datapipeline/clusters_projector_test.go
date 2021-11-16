package datapipeline

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
)

func TestClustersProjector_ClusterDiscoveryHandler(t *testing.T) {
	db := helpers.SetupTestDatabase(t)

	tx := db.Begin()
	defer tx.Rollback()

	tx.AutoMigrate(&entities.Cluster{})
	tx.Create(&entities.Cluster{
		Name:        "test_cluster",
		ID:          "test_id",
		ClusterType: models.ClusterTypeUnknown,
	})

	jsonFile, err := os.Open("./test/fixtures/cluster_discovery_hana_scale_up.json")
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

	clustersProjector_ClusterDiscoveryHandler(dataCollectedEvent, tx)

	var cluster entities.Cluster
	tx.First(&cluster)

	assert.Equal(t, "test_id", cluster.ID)
	assert.Equal(t, models.ClusterTypeScaleUp, cluster.ClusterType)
	assert.Equal(t, pq.StringArray{"DEV"}, cluster.SIDs)
	assert.Equal(t, 5, cluster.ResourcesNumber)
	assert.Equal(t, 3, cluster.HostsNumber)
}

func TestTransformClusterListData_HANAScaleUp(t *testing.T) {
	jsonFile, err := os.Open("./test/fixtures/cluster_discovery_hana_scale_up.json")
	if err != nil {
		panic(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var clusterIn cluster.Cluster
	json.Unmarshal(byteValue, &clusterIn)
	clusterOut, _ := transformClusterData(&clusterIn)

	assert.EqualValues(t,
		&entities.Cluster{
			Name:            "test_cluster",
			ID:              "test_id",
			ClusterType:     models.ClusterTypeScaleUp,
			SIDs:            []string{"DEV"},
			ResourcesNumber: 5,
			HostsNumber:     3,
		}, clusterOut)
}

func TestTransformClusterListData_Unknown(t *testing.T) {
	jsonFile, err := os.Open("./test/fixtures/cluster_discovery_unknown.json")
	if err != nil {
		panic(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var clusterIn cluster.Cluster
	json.Unmarshal(byteValue, &clusterIn)
	clusterOut, _ := transformClusterData(&clusterIn)

	assert.EqualValues(t,
		&entities.Cluster{
			Name:        "test_cluster",
			ID:          "test_id",
			ClusterType: models.ClusterTypeUnknown,
		}, clusterOut)
}
