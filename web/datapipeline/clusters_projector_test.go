package datapipeline

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/datatypes"
)

func TestClustersProjector_ClusterDiscoveryHandler(t *testing.T) {
	db := helpers.SetupTestDatabase(t)

	tx := db.Begin()
	defer tx.Rollback()

	tx.AutoMigrate(&entities.Cluster{}, &entities.HealthState{})
	tx.Create(&entities.Cluster{
		Name:        "test_cluster",
		ID:          "test_id",
		ClusterType: models.ClusterTypeUnknown,
	})

	jsonFile, err := os.Open("./test/fixtures/discovery/cluster/cluster_discovery_hana_scale_up.json")
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

	var health entities.HealthState
	tx.First(&health)

	var partialHealth map[string]string
	json.Unmarshal(health.PartialHealths, &partialHealth)

	assert.Equal(t, "5dfbd28f35cbfb38969f9b99243ae8d4", cluster.ID)
	assert.Equal(t, models.ClusterTypeHANAScaleUp, cluster.ClusterType)
	assert.Equal(t, "PRD", cluster.SID)
	assert.Equal(t, 8, cluster.ResourcesNumber)
	assert.Equal(t, 2, cluster.HostsNumber)
	assert.NotNil(t, cluster.Details)
	assert.Equal(t, health.ID, cluster.ID)
	assert.Equal(t, "critical", health.Health)
	assert.Equal(t, map[string]string{"hana_sr_health": "critical"}, partialHealth)
}

func TestTransformClusterData_HANAScaleUp(t *testing.T) {
	jsonFile, err := os.Open("./test/fixtures/discovery/cluster/cluster_discovery_hana_scale_up.json")
	if err != nil {
		panic(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var clusterIn cluster.Cluster
	json.Unmarshal(byteValue, &clusterIn)
	clusterOut, _ := transformClusterData(&clusterIn)

	expectedHANAClusterDetails, _ := json.Marshal(
		&entities.HANAClusterDetails{
			SystemReplicationMode:          "sync",
			SystemReplicationOperationMode: "logreplay",
			SecondarySyncState:             "SFAIL",
			SRHealthState:                  "4",
			CIBLastWritten:                 time.Date(2021, time.November, 6, 19, 8, 41, 0, time.UTC),
			FencingType:                    "external/sbd",
			StoppedResources: []*entities.ClusterResource{
				{
					ID:        "stopped_dummy_resource",
					Type:      "",
					Role:      "",
					Status:    "",
					FailCount: 0,
				},
			},
			Nodes: []*entities.HANAClusterNode{
				{
					Name:       "vmhana01",
					Site:       "Site1",
					VirtualIPs: []string{"10.74.1.12"},
					HANAStatus: models.HANAStatusPrimary,
					Attributes: map[string]string{
						"hana_prd_clone_state":         "PROMOTED",
						"hana_prd_op_mode":             "logreplay",
						"hana_prd_remoteHost":          "vmhana02",
						"hana_prd_roles":               "4:P:master1:master:worker:master",
						"hana_prd_site":                "Site1",
						"hana_prd_srmode":              "sync",
						"hana_prd_sync_state":          "PRIM",
						"hana_prd_version":             "2.00.030.00.1522210459",
						"hana_prd_vhost":               "vmhana01",
						"lpa_prd_lpt":                  "1636225720",
						"master-rsc_SAPHana_PRD_HDB00": "150",
					},
					Resources: []*entities.ClusterResource{
						{
							ID:        "stonith-sbd",
							Type:      "stonith:external/sbd",
							Role:      "Started",
							Status:    "active",
							FailCount: 0,
						},
						{
							ID:        "rsc_exporter_PRD_HDB00",
							Type:      "systemd:prometheus-hanadb_exporter@PRD_HDB00",
							Role:      "Started",
							Status:    "active",
							FailCount: 0,
						},
						{
							ID:        "rsc_ip_PRD_HDB00",
							Type:      "ocf::heartbeat:IPaddr2",
							Role:      "Started",
							Status:    "active",
							FailCount: 0,
						},
						{
							ID:        "rsc_socat_PRD_HDB00",
							Type:      "ocf::heartbeat:azure-lb",
							Role:      "Started",
							Status:    "active",
							FailCount: 0,
						},
						{
							ID:        "rsc_SAPHana_PRD_HDB00",
							Type:      "ocf::suse:SAPHana",
							Role:      "Master",
							Status:    "active",
							FailCount: 0,
						},
						{
							ID:        "rsc_SAPHanaTopology_PRD_HDB00",
							Type:      "ocf::suse:SAPHanaTopology",
							Role:      "Started",
							Status:    "active",
							FailCount: 0,
						},
					},
				},
				{
					Name: "vmhana02",
					Site: "Site2",
					Attributes: map[string]string{
						"hana_prd_clone_state":         "DEMOTED",
						"hana_prd_op_mode":             "logreplay",
						"hana_prd_remoteHost":          "vmhana01",
						"hana_prd_roles":               "4:S:master1:master:worker:master",
						"hana_prd_site":                "Site2",
						"hana_prd_srmode":              "sync",
						"hana_prd_sync_state":          "SFAIL",
						"hana_prd_version":             "2.00.030.00.1522210459",
						"hana_prd_vhost":               "vmhana02",
						"lpa_prd_lpt":                  "10",
						"master-rsc_SAPHana_PRD_HDB00": "-INFINITY",
					},
					Resources: []*entities.ClusterResource{
						{
							ID:        "rsc_SAPHana_PRD_HDB00",
							Type:      "ocf::suse:SAPHana",
							Role:      "Slave",
							Status:    "active",
							FailCount: 1,
						},
						{
							ID:        "rsc_SAPHanaTopology_PRD_HDB00",
							Type:      "ocf::suse:SAPHanaTopology",
							Role:      "Started",
							Status:    "active",
							FailCount: 0,
						},
					},
					VirtualIPs: nil,
					HANAStatus: models.HANAStatusFailed,
				},
			},
			SBDDevices: []*entities.SBDDevice{
				{
					Device: "/dev/disk/by-id/scsi-SLIO-ORG_IBLOCK_649b292b-ae9d-49a4-8002-2e602a0ab56e",
					Status: "healthy",
				},
				{
					Device: "/dev/disk/by-id/scsi-SLIO-ORG_IBLOCK_649b292b-ae9d-49a4-8002-2e602a012345",
					Status: "unhealthy",
				},
			},
		},
	)

	assert.EqualValues(t,
		&entities.Cluster{
			Name:            "hana_cluster",
			ID:              "5dfbd28f35cbfb38969f9b99243ae8d4",
			ClusterType:     models.ClusterTypeHANAScaleUp,
			SID:             "PRD",
			ResourcesNumber: 8,
			HostsNumber:     2,
			Details:         expectedHANAClusterDetails,
		}, clusterOut)
}

func TestTransformClusterData_Unknown(t *testing.T) {
	jsonFile, err := os.Open("./test/fixtures/discovery/cluster/cluster_discovery_unknown.json")
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
			Details:     datatypes.JSON{},
		}, clusterOut)
}

func TestParseHANAStatus_Primary(t *testing.T) {
	node := &entities.HANAClusterNode{
		Attributes: map[string]string{
			"hana_prd_roles":      "4:P:master1:master:worker:master",
			"hana_prd_sync_state": "PRIM",
		},
	}

	status := parseHANAStatus(node, "PRD")
	assert.Equal(t, models.HANAStatusPrimary, status)
}

func TestParseHANAStatus_Secondary(t *testing.T) {
	node := &entities.HANAClusterNode{
		Attributes: map[string]string{
			"hana_prd_roles":      "4:S:master1:master:worker:master",
			"hana_prd_sync_state": "SOK",
		},
	}

	status := parseHANAStatus(node, "PRD")
	assert.Equal(t, models.HANAStatusSecondary, status)
}

func TestParseHANAStatus_Unknown(t *testing.T) {
	node := &entities.HANAClusterNode{
		Attributes: map[string]string{
			"hana_prd_roles": "4:S:master1:master:worker:master",
		},
	}

	status := parseHANAStatus(node, "PRD")
	assert.Equal(t, models.HANAStatusUnknown, status)

	node = &entities.HANAClusterNode{
		Attributes: map[string]string{
			"hana_prd_sync_state": "SOK",
		},
	}

	status = parseHANAStatus(node, "PRD")
	assert.Equal(t, models.HANAStatusUnknown, status)
}

func TestParseHANAStatus_Failed(t *testing.T) {
	node := &entities.HANAClusterNode{
		Attributes: map[string]string{
			"hana_prd_roles":      "4:P:master1:master:worker:master",
			"hana_prd_sync_state": "SFAIL",
		},
	}

	status := parseHANAStatus(node, "PRD")
	assert.Equal(t, models.HANAStatusFailed, status)

	node = &entities.HANAClusterNode{
		Attributes: map[string]string{
			"hana_prd_roles":      "4:S:master1:master:worker:master",
			"hana_prd_sync_state": "PRIM",
		},
	}

	status = parseHANAStatus(node, "PRD")
	assert.Equal(t, models.HANAStatusFailed, status)
}

func TestParseSecondarySyncState(t *testing.T) {
	nodes := []*entities.HANAClusterNode{
		&entities.HANAClusterNode{
			Attributes: map[string]string{
				"hana_prd_roles":      "4:P:master1:master:worker:master",
				"hana_prd_sync_state": "PRIM",
			},
		},
		&entities.HANAClusterNode{
			Attributes: map[string]string{
				"hana_prd_roles":      "4:S:master1:master:worker:master",
				"hana_prd_sync_state": "SOK",
			},
		},
	}

	state := parseHANASecondarySyncState(nodes, "PRD")
	assert.Equal(t, "SOK", state)
}

func TestParseSecondaryHealthState(t *testing.T) {
	nodes := []*entities.HANAClusterNode{
		&entities.HANAClusterNode{
			Attributes: map[string]string{
				"hana_prd_roles":      "4:P:master1:master:worker:master",
				"hana_prd_sync_state": "PRIM",
			},
		},
		&entities.HANAClusterNode{
			Attributes: map[string]string{
				"hana_prd_roles":      "4:S:master1:master:worker:master",
				"hana_prd_sync_state": "SOK",
			},
		},
	}

	state := parseHANAHealthState(nodes, "PRD")
	assert.Equal(t, "4", state)
}
