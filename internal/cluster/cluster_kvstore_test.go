package cluster

import (
	"fmt"
	"os"
	"testing"

	"github.com/ClusterLabs/ha_cluster_exporter/collector/pacemaker/cib"
	"github.com/ClusterLabs/ha_cluster_exporter/collector/pacemaker/crmmon"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/consul/mocks"
)

type CrmConfig struct {
	ClusterProperties []cib.Attribute `xml:"cluster_property_set>nvpair"`
}

type Configuration struct {
	CrmConfig CrmConfig `xml:"crm_config"`
}

func TestStore(t *testing.T) {
	host, _ := os.Hostname()
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)

	consulInst.On("KV").Return(kv)

	kvPath := fmt.Sprintf("%s%s", consul.KvClustersPath, "cluster_name")

	expectedPutMap := map[string]interface{}{
		"cib": map[string]interface{}{
			"Configuration": map[string]interface{}{
				"Constraints": map[string]interface{}{
					"RscLocations": []struct {
						Id       string "xml:\"id,attr\""
						Node     string "xml:\"node,attr\""
						Resource string "xml:\"rsc,attr\""
						Role     string "xml:\"role,attr\""
						Score    string "xml:\"score,attr\""
					}(nil),
				},
				"CrmConfig": map[string]interface{}{
					"ClusterProperties": []cib.Attribute{
						cib.Attribute{
							Id:    "cib-bootstrap-options-cluster-name",
							Name:  "",
							Value: "cluster_name",
						},
					},
				},
				"Nodes": []struct {
					Id                 string          "xml:\"id,attr\""
					Uname              string          "xml:\"uname,attr\""
					InstanceAttributes []cib.Attribute "xml:\"instance_attributes>nvpair\""
				}(nil),
				"Resources": map[string]interface{}{
					"Clones":     []cib.Clone(nil),
					"Masters":    []cib.Clone(nil),
					"Primitives": []cib.Primitive(nil),
				},
			},
		},
		"crmmon": map[string]interface{}{
			"Clones": []crmmon.Clone(nil),
			"Groups": []crmmon.Group(nil),
			"NodeAttributes": map[string]interface{}{
				"Nodes": []struct {
					Name       string "xml:\"name,attr\""
					Attributes []struct {
						Name  string "xml:\"name,attr\""
						Value string "xml:\"value,attr\""
					} "xml:\"attribute\""
				}(nil),
			},
			"NodeHistory": map[string]interface{}{
				"Nodes": []struct {
					Name            string "xml:\"name,attr\""
					ResourceHistory []struct {
						Name               string "xml:\"id,attr\""
						MigrationThreshold int    "xml:\"migration-threshold,attr\""
						FailCount          int    "xml:\"fail-count,attr\""
					} "xml:\"resource_history\""
				}(nil),
			},
			"Nodes": []crmmon.Node{
				crmmon.Node{
					Name:             "othernode",
					Id:               "",
					Online:           false,
					Standby:          false,
					StandbyOnFail:    false,
					Maintenance:      false,
					Pending:          false,
					Unclean:          false,
					Shutdown:         false,
					ExpectedUp:       false,
					DC:               false,
					ResourcesRunning: 0,
					Type:             "",
				},
				crmmon.Node{
					Name:             host,
					Id:               "",
					Online:           false,
					Standby:          false,
					StandbyOnFail:    false,
					Maintenance:      false,
					Pending:          false,
					Unclean:          false,
					Shutdown:         false,
					ExpectedUp:       false,
					DC:               true,
					ResourcesRunning: 0,
					Type:             "",
				},
			},
			"Resources": []crmmon.Resource(nil),
			"Summary": map[string]interface{}{
				"ClusterOptions": map[string]interface{}{
					"StonithEnabled": false,
				},
				"LastChange": map[string]interface{}{
					"Time": "",
				},
				"Nodes": map[string]interface{}{
					"Number": 0,
				},
				"Resources": map[string]interface{}{
					"Blocked":  0,
					"Disabled": 0,
					"Number":   0,
				},
			},
			"Version": "1.2.3",
		},
	}

	kv.On("DeleteTree", kvPath, (*consulApi.WriteOptions)(nil)).Return(nil, nil)
	kv.On("PutMap", kvPath, expectedPutMap).Return(nil, nil)
	testLock := consulApi.Lock{}
	consulInst.On("LockTrento", consul.KvClustersPath).Return(&testLock, nil)

	root := new(cib.Root)

	crmConfig := struct {
		ClusterProperties []cib.Attribute `xml:"cluster_property_set>nvpair"`
	}{
		ClusterProperties: []cib.Attribute{
			cib.Attribute{
				Id:    "cib-bootstrap-options-cluster-name",
				Value: "cluster_name",
			},
		},
	}

	root.Configuration.CrmConfig = crmConfig

	c := Cluster{
		Cib: *root,
		Crmmon: crmmon.Root{
			Version: "1.2.3",
			Nodes: []crmmon.Node{
				crmmon.Node{
					Name: "othernode",
					DC:   false,
				},
				crmmon.Node{
					Name: host,
					DC:   true,
				},
			},
		},
	}

	result := c.Store(consulInst)
	assert.Equal(t, nil, result)
}

func TestLoad(t *testing.T) {
	host, _ := os.Hostname()
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)

	listMap := map[string]interface{}{
		"test_cluster": map[string]interface{}{
			"cib": map[string]interface{}{
				"Configuration": map[string]interface{}{
					"CrmConfig": map[string]interface{}{
						"ClusterProperties": []interface{}{
							map[string]interface{}{
								"Id":    "cib-bootstrap-options-cluster-name",
								"Value": "cluster_name",
							},
						},
					},
				},
			},
			"crmmon": map[string]interface{}{
				"Version": "1.2.3",
				"Nodes": []interface{}{
					map[string]interface{}{
						"Name": "othernode",
						"DC":   false,
					},
					map[string]interface{}{
						"Name": host,
						"DC":   true,
					},
				},
			},
		},
	}

	kv.On("ListMap", consul.KvClustersPath, consul.KvClustersPath).Return(listMap, nil)
	consulInst.On("WaitLock", consul.KvClustersPath).Return(nil)

	consulInst.On("KV").Return(kv)

	c, _ := Load(consulInst)

	root := new(cib.Root)

	crmConfig := struct {
		ClusterProperties []cib.Attribute `xml:"cluster_property_set>nvpair"`
	}{
		ClusterProperties: []cib.Attribute{
			cib.Attribute{
				Id:    "cib-bootstrap-options-cluster-name",
				Value: "cluster_name",
			},
		},
	}

	root.Configuration.CrmConfig = crmConfig

	expectedCluster := map[string]*Cluster{
		"test_cluster": &Cluster{
			Cib: *root,
			Crmmon: crmmon.Root{
				Version: "1.2.3",
				Nodes: []crmmon.Node{
					crmmon.Node{
						Name: "othernode",
						DC:   false,
					},
					crmmon.Node{
						Name: host,
						DC:   true,
					},
				},
			},
		},
	}

	assert.Equal(t, expectedCluster, c)
}
