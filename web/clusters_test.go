package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"testing"

	consulApi "github.com/hashicorp/consul/api"

	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/trento-project/trento/internal/cloud"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/consul/mocks"
	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
	serviceMocks "github.com/trento-project/trento/web/services/mocks"
)

func clustersListMap() map[string]interface{} {
	listMap := map[string]interface{}{
		"47d1190ffb4f781974c8356d7f863b03": map[string]interface{}{
			"discovered_data": map[string]interface{}{
				"cib": map[string]interface{}{
					"Configuration": map[string]interface{}{
						"Resources": map[string]interface{}{
							"Clones": []interface{}{
								map[string]interface{}{
									"Primitive": map[string]interface{}{
										"Type": "SAPHanaTopology",
										"InstanceAttributes": []interface{}{
											map[string]interface{}{
												"Name":  "SID",
												"Value": "PRD",
											},
										},
									},
								},
							},
							"Groups": []interface{}{
								map[string]interface{}{
									"Primitives": []interface{}{
										map[string]interface{}{
											"Id": "ip",
											"InstanceAttributes": []interface{}{
												map[string]interface{}{
													"Name":  "ip",
													"Value": "10.123.123.123",
												},
											},
										},
									},
								},
							},
						},
					},
					"CrmConfig": map[string]interface{}{
						"ClusterProperties": []interface{}{
							map[string]interface{}{
								"Id":    "cib-bootstrap-options-cluster-name",
								"Value": "hana_cluster",
							},
						},
					},
				},
				"crmmon": map[string]interface{}{
					"Clones": []interface{}{
						map[string]interface{}{
							"Resources": []interface{}{
								map[string]interface{}{
									"Agent": "ocf::suse:SAPHana",
									"Node": map[string]interface{}{
										"Name": "test_node_1",
									},
								},
								map[string]interface{}{
									"Agent": "ocf::suse:SAPHanaTopology",
									"Node": map[string]interface{}{
										"Name": "test_node_1",
									},
								},
							},
						},
					},
					"Nodes": []interface{}{
						map[string]interface{}{
							"Name": "test_node_1",
						},
						map[string]interface{}{
							"Name": "test_node_2",
						},
					},
					"NodeAttributes": map[string]interface{}{
						"Nodes": []interface{}{
							map[string]interface{}{
								"Name": "test_node_1",
								"Attributes": []interface{}{
									map[string]interface{}{
										"Name":  "hana_prd_srmode",
										"Value": "sync",
									},
									map[string]interface{}{
										"Name":  "hana_prd_op_mode",
										"Value": "logreplay",
									},
									map[string]interface{}{
										"Name":  "hana_prd_roles",
										"Value": "4:P:master1:master:worker:master",
									},
									map[string]interface{}{
										"Name":  "hana_prd_sync_state",
										"Value": "PRIM",
									},
									map[string]interface{}{
										"Name":  "hana_prd_site",
										"Value": "site1",
									},
								},
							},
							map[string]interface{}{
								"Name": "test_node_2",
								"Attributes": []interface{}{
									map[string]interface{}{
										"Name":  "hana_prd_srmode",
										"Value": "sync",
									},
									map[string]interface{}{
										"Name":  "hana_prd_op_mode",
										"Value": "logreplay",
									},
									map[string]interface{}{
										"Name":  "hana_prd_roles",
										"Value": "4:S:master1:master:worker:master",
									},
									map[string]interface{}{
										"Name":  "hana_prd_sync_state",
										"Value": "SFAIL",
									},
									map[string]interface{}{
										"Name":  "hana_prd_site",
										"Value": "site2",
									},
								},
							},
						},
					},
					"Summary": map[string]interface{}{
						"Nodes": map[string]interface{}{
							"Number": 3,
						},
						"LastChange": map[string]interface{}{
							"Time": "Wed Jun 30 18:11:37 2021",
						},
						"Resources": map[string]interface{}{
							"Number": 5,
						},
					},
					"Resources": []interface{}{
						map[string]interface{}{
							"Id":     "ip",
							"Agent":  "ocf::heartbeat:IPaddr2",
							"Role":   "Started",
							"Active": true,
							"Node": map[string]interface{}{
								"Name": "test_node_1",
							},
						},
						map[string]interface{}{
							"Id":     "sbd",
							"Agent":  "stonith:external/sbd",
							"Role":   "Started",
							"Active": true,
							"Node": map[string]interface{}{
								"Name": "test_node_1",
							},
						},
						map[string]interface{}{
							"Id":     "dummy_failed",
							"Agent":  "dummy",
							"Role":   "Started",
							"Failed": true,
							"Node": map[string]interface{}{
								"Name": "test_node_1",
							},
						},
					},
				},
				"name": "hana_cluster",
				"id":   "47d1190ffb4f781974c8356d7f863b03",
			},
		},
		"e2f2eb50aef748e586a7baa85e0162cf": map[string]interface{}{
			"discovered_data": map[string]interface{}{
				"cib": map[string]interface{}{
					"Configuration": map[string]interface{}{
						"CrmConfig": map[string]interface{}{
							"ClusterProperties": []interface{}{
								map[string]interface{}{
									"Id":    "cib-bootstrap-options-cluster-name",
									"Value": "netweaver_cluster",
								},
							},
						},
					},
				},
				"crmmon": map[string]interface{}{
					"Summary": map[string]interface{}{
						"Nodes": map[string]interface{}{
							"Number": 2,
						},
						"Resources": map[string]interface{}{
							"Number": 10,
						},
					},
				},
				"name": "netweaver_cluster",
				"id":   "e2f2eb50aef748e586a7baa85e0162cf",
			},
		},
		"e27d313a674375b2066777a89ee346b9": map[string]interface{}{
			"discovered_data": map[string]interface{}{
				"cib": map[string]interface{}{
					"Configuration": map[string]interface{}{
						"CrmConfig": map[string]interface{}{
							"ClusterProperties": []interface{}{
								map[string]interface{}{
									"Id":    "cib-bootstrap-options-cluster-name",
									"Value": "netweaver_cluster",
								},
							},
						},
					},
				},
				"crmmon": map[string]interface{}{
					"Summary": map[string]interface{}{
						"Nodes": map[string]interface{}{
							"Number": 2,
						},
						"Resources": map[string]interface{}{
							"Number": 10,
						},
					},
				},
				"name": "netweaver_cluster",
				"id":   "e27d313a674375b2066777a89ee346b9",
			},
		},
	}

	return listMap
}

func checksCatalog() map[string]*models.Check {

	checksByGroup := map[string]*models.Check{
		"1.1.1": &models.Check{
			ID:             "1.1.1",
			Name:           "check 1",
			Group:          "group 1",
			Description:    "description 1",
			Remediation:    "remediation 1",
			Implementation: "implementation 1",
			Labels:         "labels 1",
		},
		"1.1.1.runtime": &models.Check{
			ID:             "1.1.1.runtime",
			Name:           "check 1 (runtime)",
			Group:          "group 1",
			Description:    "description 1",
			Remediation:    "remediation 1",
			Implementation: "implementation 1",
			Labels:         "labels 1",
		},
		"1.1.2": &models.Check{
			ID:             "1.1.2",
			Name:           "check 2",
			Group:          "group 1",
			Description:    "description 2",
			Remediation:    "remediation 2",
			Implementation: "implementation 2",
			Labels:         "labels 2",
		},
		"1.2.3": &models.Check{
			ID:             "1.2.3",
			Name:           "check 3",
			Group:          "group 2",
			Description:    "description 3",
			Remediation:    "remediation 3",
			Implementation: "implementation 3",
			Labels:         "labels 3",
		},
	}

	return checksByGroup
}

func checksCatalogByGroup() map[string]map[string]*models.Check {

	checksByGroup := map[string]map[string]*models.Check{
		"group 1": {
			"1.1.1": &models.Check{
				ID:             "1.1.1",
				Name:           "check 1",
				Group:          "group 1",
				Description:    "description 1",
				Remediation:    "remediation 1",
				Implementation: "implementation 1",
				Labels:         "labels 1",
			},
			"1.1.1.runtime": &models.Check{
				ID:             "1.1.1.runtime",
				Name:           "check 1 (runtime)",
				Group:          "group 1",
				Description:    "description 1",
				Remediation:    "remediation 1",
				Implementation: "implementation 1",
				Labels:         "labels 1",
			},
			"1.1.2": &models.Check{
				ID:             "1.1.2",
				Name:           "check 2",
				Group:          "group 1",
				Description:    "description 2",
				Remediation:    "remediation 2",
				Implementation: "implementation 2",
				Labels:         "labels 2",
			},
		},
		"group 2": {
			"1.2.3": &models.Check{
				ID:             "1.2.3",
				Name:           "check 3",
				Group:          "group 2",
				Description:    "description 3",
				Remediation:    "remediation 3",
				Implementation: "implementation 3",
				Labels:         "labels 3",
			},
		},
	}

	return checksByGroup
}

func checksResult() *models.Results {

	checksResult := &models.Results{
		Checks: map[string]*models.ChecksByHost{
			"1.1.1": &models.ChecksByHost{
				Hosts: map[string]*models.Check{
					"test_node_1": &models.Check{
						Result: true,
					},
					"test_node_2": &models.Check{
						Result: true,
					},
				},
			},
			"1.1.2": &models.ChecksByHost{
				Hosts: map[string]*models.Check{
					"test_node_1": &models.Check{
						Result: false,
					},
					"test_node_2": &models.Check{
						Result: false,
					},
				},
			},
		},
	}

	return checksResult
}

func checksResultUnreachable() *models.Results {

	checksResult := &models.Results{
		Checks: map[string]*models.ChecksByHost{
			"1.1.1": &models.ChecksByHost{
				Hosts: map[string]*models.Check{
					"test_node_1": &models.Check{
						Result: true,
					},
				},
			},
			"1.1.2": &models.ChecksByHost{
				Hosts: map[string]*models.Check{
					"test_node_1": &models.Check{
						Result: false,
					},
				},
			},
		},
	}

	return checksResult
}

func aggregatedByCluster() *services.AggregatedCheckData {
	return &services.AggregatedCheckData{
		PassingCount:  2,
		WarningCount:  0,
		CriticalCount: 0,
	}
}

func aggregatedByClusterCritical() *services.AggregatedCheckData {
	return &services.AggregatedCheckData{
		PassingCount:  2,
		WarningCount:  0,
		CriticalCount: 1,
	}
}

func aggregatedByClusterEmpty() *services.AggregatedCheckData {
	return &services.AggregatedCheckData{
		PassingCount:  0,
		WarningCount:  0,
		CriticalCount: 0,
	}
}

func checksResultByHost() map[string]*services.AggregatedCheckData {
	return map[string]*services.AggregatedCheckData{
		"test_node_1": &services.AggregatedCheckData{
			PassingCount:  2,
			WarningCount:  0,
			CriticalCount: 0,
		},
		"test_node_2": &services.AggregatedCheckData{
			PassingCount:  2,
			WarningCount:  0,
			CriticalCount: 1,
		},
	}
}

func azureMeta(userIndex int) map[string]interface{} {
	return map[string]interface{}{
		"provider": cloud.Azure,
		"metadata": &cloud.AzureMetadata{
			Compute: cloud.Compute{
				OsProfile: cloud.OsProfile{
					AdminUserName: fmt.Sprintf("defuser%d", userIndex),
				},
			},
		},
	}
}

func TestClustersListHandler(t *testing.T) {
	consulInst := new(mocks.Client)
	checksMocks := new(serviceMocks.ChecksService)

	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)
	kv.On("ListMap", consul.KvClustersPath, consul.KvClustersPath).Return(clustersListMap(), nil)
	consulInst.On("WaitLock", consul.KvClustersPath).Return(nil)

	checksMocks.On("GetAggregatedChecksResultByCluster", "47d1190ffb4f781974c8356d7f863b03").Return(
		aggregatedByCluster(), nil)
	checksMocks.On("GetAggregatedChecksResultByCluster", "e2f2eb50aef748e586a7baa85e0162cf").Return(
		aggregatedByClusterCritical(), nil)
	checksMocks.On("GetAggregatedChecksResultByCluster", "e27d313a674375b2066777a89ee346b9").Return(
		aggregatedByClusterEmpty(), nil)

	tags := map[string]interface{}{
		"tag1": struct{}{},
	}

	tagsPath := "trento/v0/tags/clusters/47d1190ffb4f781974c8356d7f863b03/"
	kv.On("ListMap", tagsPath, tagsPath).Return(tags, nil)
	tagsPath = "trento/v0/tags/clusters/e2f2eb50aef748e586a7baa85e0162cf/"
	kv.On("ListMap", tagsPath, tagsPath).Return(tags, nil)
	tagsPath = "trento/v0/tags/clusters/e27d313a674375b2066777a89ee346b9/"
	kv.On("ListMap", tagsPath, tagsPath).Return(nil, nil)

	deps := DefaultDependencies()
	deps.consul = consulInst
	deps.checksService = checksMocks

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/clusters", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	consulInst.AssertExpectations(t)
	checksMocks.AssertExpectations(t)
	kv.AssertExpectations(t)

	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepEndTags:         true,
	})
	minified, err := m.String("text/html", resp.Body.String())
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, minified, "Clusters")
	assert.Regexp(t, regexp.MustCompile("<td .*>.*check_circle.*</td><td>.*hana_cluster.*</td><td>.*47d1190ffb4f781974c8356d7f863b03.*</td><td>HANA scale-up</td><td>PRD</td><td>3</td><td>5</td><td><input.*value=tag1.*></td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td .*>.*error.*</td><td>.*duplicated.*netweaver_cluster.*</td><td>.*e2f2eb50aef748e586a7baa85e0162cf.*</td><td>Unknown</td><td></td><td>2</td><td>10</td><td><input.*value=tag1.*></td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td .*>.*fiber_manual_record.*</td><td>.*duplicated.*netweaver_cluster.*</td><td>.*e27d313a674375b2066777a89ee346b9.*</td><td>Unknown</td><td></td>"), minified)
}

func TestClusterHandlerHANA(t *testing.T) {
	clusterId := "47d1190ffb4f781974c8356d7f863b03"

	nodes := []*consulApi.Node{
		{
			Node:    "test_node_1",
			Address: "192.168.1.1",
		},
		{
			Node:    "test_node_2",
			Address: "192.168.1.2",
		},
	}

	consulInst := new(mocks.Client)
	checksMocks := new(serviceMocks.ChecksService)

	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)
	kv.On("ListMap", consul.KvClustersPath, consul.KvClustersPath).Return(clustersListMap(), nil)
	consulInst.On("WaitLock", consul.KvClustersPath).Return(nil).Times(3)

	catalog := new(mocks.Catalog)
	filter := &consulApi.QueryOptions{Filter: "Meta[\"trento-ha-cluster-id\"] == \"" + clusterId + "\""}
	catalog.On("Nodes", filter).Return(nodes, nil, nil)
	consulInst.On("Catalog").Return(catalog)

	selectedChecksPath := fmt.Sprintf(consul.KvClustersChecksPath, clusterId)
	selectedChecksValue := &consulApi.KVPair{Value: []byte("1.1.1,1.1.2")}
	kv.On("Get", selectedChecksPath, (*consulApi.QueryOptions)(nil)).Return(selectedChecksValue, nil, nil)

	connData := map[string]interface{}{
		"test_node_1": "myuser1",
		"test_node_2": "myuser2",
	}
	connSettingsPath := fmt.Sprintf(consul.KvClustersConnectionPath, clusterId)
	kv.On("ListMap", connSettingsPath, connSettingsPath).Return(connData, nil)

	cloudPath1 := fmt.Sprintf(consul.KvHostsClouddataPath, "test_node_1")
	consulInst.On("WaitLock", path.Join(consul.KvHostsPath, "test_node_1")+"/").Return(nil)
	kv.On("ListMap", cloudPath1, cloudPath1).Return(azureMeta(1), nil)

	cloudPath2 := fmt.Sprintf(consul.KvHostsClouddataPath, "test_node_2")
	consulInst.On("WaitLock", path.Join(consul.KvHostsPath, "test_node_2")+"/").Return(nil)
	kv.On("ListMap", cloudPath2, cloudPath2).Return(azureMeta(2), nil)

	checksMocks.On("GetChecksCatalog").Return(checksCatalog(), nil)
	checksMocks.On("GetChecksCatalogByGroup").Return(checksCatalogByGroup(), nil)
	checksMocks.On("GetChecksResultByCluster", clusterId).Return(
		checksResult(), nil)
	checksMocks.On("GetAggregatedChecksResultByCluster", clusterId).Return(
		aggregatedByClusterCritical(), nil)
	checksMocks.On("GetAggregatedChecksResultByHost", clusterId).Return(
		checksResultByHost(), nil)

	deps := DefaultDependencies()
	deps.consul = consulInst
	deps.checksService = checksMocks

	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/clusters/"+clusterId, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

	consulInst.AssertExpectations(t)
	kv.AssertExpectations(t)
	checksMocks.AssertExpectations(t)

	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepEndTags:         true,
	})
	minified, err := m.String("text/html", resp.Body.String())
	if err != nil {
		panic(err)
	}

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "Cluster details")
	// Summary
	assert.Regexp(t, regexp.MustCompile("<strong>Cluster name:</strong><br><span.*>hana_cluster</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<strong>Cluster type:</strong><br><span.*>HANA scale-up</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<strong>HANA system replication mode:</strong><br><span.*>sync</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<strong>Stonith type:</strong><br><span.*>external/sbd</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<strong>HANA system replication operation mode:</strong><br><span.*>logreplay</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<strong>CIB last written:</strong><br><span.*>Wed Jun 30 18:11:37 2021</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<strong>HANA secondary sync state:</strong><br><span.*>SFAIL</span>"), minified)
	// Health
	assert.Regexp(t, regexp.MustCompile(".*check_circle.*alert-body.*Passing.*2"), minified)
	assert.Regexp(t, regexp.MustCompile(".*error.*alert-body.*Critical.*1"), minified)

	// Nodes
	assert.Regexp(t, regexp.MustCompile("<td.*check_circle.*<td><a.*href=/hosts/test_node_1.*>test_node_1</a></td><td>192.168.1.1</td><td>10.123.123.123</td><td><span .*>HANA Primary</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td.*error.*<td><a.*href=/hosts/test_node_2.*>test_node_2</a></td><td>192.168.1.2</td>"), minified)
	// Resources
	assert.Regexp(t, regexp.MustCompile("<td>sbd</td><td>stonith:external/sbd</td><td>Started</td><td>active</td><td>0</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>dummy_failed</td><td>dummy</td><td>Started</td><td>failed</td><td>0</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<h4>Stopped resources</h4><div.*><div.*><span .*>dummy_failed</span>"), minified)
	// Connection settings
	assert.Regexp(t, regexp.MustCompile("<td>test_node_1</td>.*<td><input.*id=username-test_node_1.*value=myuser1.*<td>defuser1"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>test_node_2</td>.*<td><input.*id=username-test_node_2.*value=myuser2.*<td>defuser2"), minified)
	// Selected checks
	assert.Regexp(t, regexp.MustCompile("id=0-1-1-1 checked>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>1.1.1</td><td>description 1</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("id=0-1-1-1-runtime>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>1.1.1.runtime</td><td>description 1</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("id=0-1-1-2 checked>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>1.1.2</td><td>description 2</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("id=1-1-2-3>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>1.2.3</td><td>description 3</td>"), minified)
	// Checks result modal
	assert.Regexp(t, regexp.MustCompile("<th.*>test_node_1.*<th.*>test_node_2"), minified)
	assert.Regexp(t, regexp.MustCompile("<td.*>1.1.1</td><td.*>description 1</td><td.*>.*check_circle.*check_circle"), minified)
	assert.Regexp(t, regexp.MustCompile("<td.*>1.1.2</td><td.*>description 2</td><td.*>.*error.*error"), minified)
}

func TestClusterHandlerUnreachableNodes(t *testing.T) {
	clusterId := "47d1190ffb4f781974c8356d7f863b03"

	nodes := []*consulApi.Node{
		{
			Node:    "test_node_1",
			Address: "192.168.1.1",
		},
		{
			Node:    "test_node_2",
			Address: "192.168.1.2",
		},
	}

	consulInst := new(mocks.Client)
	checksMocks := new(serviceMocks.ChecksService)

	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)
	kv.On("ListMap", consul.KvClustersPath, consul.KvClustersPath).Return(clustersListMap(), nil)
	consulInst.On("WaitLock", consul.KvClustersPath).Return(nil).Times(3)

	catalog := new(mocks.Catalog)
	filter := &consulApi.QueryOptions{Filter: "Meta[\"trento-ha-cluster-id\"] == \"" + clusterId + "\""}
	catalog.On("Nodes", filter).Return(nodes, nil, nil)
	consulInst.On("Catalog").Return(catalog)

	selectedChecksPath := fmt.Sprintf(consul.KvClustersChecksPath, clusterId)
	selectedChecksValue := &consulApi.KVPair{Value: []byte("1.1.1,1.1.2")}
	kv.On("Get", selectedChecksPath, (*consulApi.QueryOptions)(nil)).Return(selectedChecksValue, nil, nil)

	connData := map[string]interface{}{
		"test_node_1": "myuser1",
		"test_node_2": "myuser2",
	}
	connSettingsPath := fmt.Sprintf(consul.KvClustersConnectionPath, clusterId)
	kv.On("ListMap", connSettingsPath, connSettingsPath).Return(connData, nil)

	cloudPath1 := fmt.Sprintf(consul.KvHostsClouddataPath, "test_node_1")
	consulInst.On("WaitLock", path.Join(consul.KvHostsPath, "test_node_1")+"/").Return(nil)
	kv.On("ListMap", cloudPath1, cloudPath1).Return(azureMeta(1), nil)

	cloudPath2 := fmt.Sprintf(consul.KvHostsClouddataPath, "test_node_2")
	consulInst.On("WaitLock", path.Join(consul.KvHostsPath, "test_node_2")+"/").Return(nil)
	kv.On("ListMap", cloudPath2, cloudPath2).Return(azureMeta(2), nil)

	checksMocks.On("GetChecksCatalog").Return(checksCatalog(), nil)
	checksMocks.On("GetChecksCatalogByGroup").Return(checksCatalogByGroup(), nil)
	checksMocks.On("GetChecksResultByCluster", clusterId).Return(checksResultUnreachable(), nil)
	checksMocks.On("GetAggregatedChecksResultByCluster", clusterId).Return(
		aggregatedByClusterCritical(), nil)
	checksMocks.On("GetAggregatedChecksResultByHost", clusterId).Return(
		checksResultByHost(), nil)

	deps := DefaultDependencies()
	deps.consul = consulInst
	deps.checksService = checksMocks

	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/clusters/"+clusterId, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

	consulInst.AssertExpectations(t)
	kv.AssertExpectations(t)
	checksMocks.AssertExpectations(t)

	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepEndTags:         true,
	})
	minified, err := m.String("text/html", resp.Body.String())
	if err != nil {
		panic(err)
	}

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	// Checks result modal
	assert.Regexp(t, regexp.MustCompile("<th.*>test_node_1.*<th.*>.*warning.*test_node_2"), minified)
	assert.Regexp(t, regexp.MustCompile("<td.*>1.1.1</td><td.*>description 1</td><td.*>.*check_circle.*sync_problem"), minified)
	assert.Regexp(t, regexp.MustCompile("<td.*>1.1.2</td><td.*>description 2</td><td.*>.*error.*sync_problem"), minified)
}

func TestClusterHandlerAlert(t *testing.T) {
	clusterId := "47d1190ffb4f781974c8356d7f863b03"
	nodes := []*consulApi.Node{
		{
			Node:    "test_node_1",
			Address: "192.168.1.1",
		},
		{
			Node:    "test_node_2",
			Address: "192.168.1.2",
		},
	}

	consulInst := new(mocks.Client)
	checksMocks := new(serviceMocks.ChecksService)

	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)
	kv.On("ListMap", consul.KvClustersPath, consul.KvClustersPath).Return(clustersListMap(), nil)
	consulInst.On("WaitLock", consul.KvClustersPath).Return(nil).Times(3)

	catalog := new(mocks.Catalog)
	filter := &consulApi.QueryOptions{Filter: "Meta[\"trento-ha-cluster-id\"] == \"" + clusterId + "\""}
	catalog.On("Nodes", filter).Return(nodes, nil, nil)
	consulInst.On("Catalog").Return(catalog)

	selectedChecksPath := fmt.Sprintf(consul.KvClustersChecksPath, clusterId)
	selectedChecksValue := &consulApi.KVPair{Value: []byte("1.1.1,1.2.3")}
	kv.On("Get", selectedChecksPath, (*consulApi.QueryOptions)(nil)).Return(selectedChecksValue, nil, nil)

	connData := map[string]interface{}{
		"test_node_1": "myuser1",
		"test_node_2": "myuser2",
	}
	connSettingsPath := fmt.Sprintf(consul.KvClustersConnectionPath, clusterId)
	kv.On("ListMap", connSettingsPath, connSettingsPath).Return(connData, nil)

	cloudPath1 := fmt.Sprintf(consul.KvHostsClouddataPath, "test_node_1")
	consulInst.On("WaitLock", path.Join(consul.KvHostsPath, "test_node_1")+"/").Return(nil)
	kv.On("ListMap", cloudPath1, cloudPath1).Return(azureMeta(1), nil)

	cloudPath2 := fmt.Sprintf(consul.KvHostsClouddataPath, "test_node_2")
	consulInst.On("WaitLock", path.Join(consul.KvHostsPath, "test_node_2")+"/").Return(nil)
	kv.On("ListMap", cloudPath2, cloudPath2).Return(azureMeta(2), nil)

	checksMocks.On("GetChecksCatalog").Return(nil, fmt.Errorf("catalog error"))
	checksMocks.On("GetChecksCatalogByGroup").Return(nil, fmt.Errorf("catalog error"))
	checksMocks.On("GetChecksResultByCluster", clusterId).Return(nil, fmt.Errorf("catalog error"))
	checksMocks.On("GetAggregatedChecksResultByCluster", clusterId).Return(
		aggregatedByClusterCritical(), nil)
	checksMocks.On("GetAggregatedChecksResultByHost", clusterId).Return(
		checksResultByHost(), nil)

	deps := DefaultDependencies()
	deps.consul = consulInst
	deps.checksService = checksMocks

	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/clusters/"+clusterId, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

	consulInst.AssertExpectations(t)
	kv.AssertExpectations(t)
	checksMocks.AssertExpectations(t)

	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepEndTags:         true,
	})
	minified, err := m.String("text/html", resp.Body.String())
	if err != nil {
		panic(err)
	}

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "Cluster details")

	assert.Regexp(t, regexp.MustCompile("Error loading the checks catalog"), minified)
	assert.Equal(t, 1, strings.Count(minified, "Error loading the checks catalog"))
}

func TestClusterHandlerGeneric(t *testing.T) {
	consulInst := new(mocks.Client)

	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)

	kv.On("ListMap", consul.KvClustersPath, consul.KvClustersPath).Return(clustersListMap(), nil)
	consulInst.On("WaitLock", consul.KvClustersPath).Return(nil)

	catalog := new(mocks.Catalog)
	filter := &consulApi.QueryOptions{Filter: "Meta[\"trento-ha-cluster-id\"] == \"e2f2eb50aef748e586a7baa85e0162cf\""}
	catalog.On("Nodes", filter).Return(nil, nil, nil)
	consulInst.On("Catalog").Return(catalog)

	deps := DefaultDependencies()
	deps.consul = consulInst

	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/clusters/e2f2eb50aef748e586a7baa85e0162cf", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "Cluster details")
	assert.Contains(t, resp.Body.String(), "netweaver_cluster")
	assert.NotContains(t, resp.Body.String(), "HANA scale-out")
	assert.NotContains(t, resp.Body.String(), "HANA scale-up")
}

func TestClusterHandler404Error(t *testing.T) {
	var err error

	kv := new(mocks.KV)
	kv.On("ListMap", consul.KvClustersPath, consul.KvClustersPath).Return(clustersListMap(), nil)

	consulInst := new(mocks.Client)
	consulInst.On("KV").Return(kv)
	consulInst.On("WaitLock", consul.KvClustersPath).Return(nil)

	deps := DefaultDependencies()
	deps.consul = consulInst

	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/clusters/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.Code)
	assert.Contains(t, resp.Body.String(), "Not Found")
}

func TestSaveChecksHandler(t *testing.T) {
	var err error

	expectedConnections := map[string]interface{}{
		"host1": "myuser1",
		"host2": "myuser2",
	}

	kv := new(mocks.KV)
	kv.On("PutTyped", fmt.Sprintf(consul.KvClustersChecksPath, "foobar"), "1.2.3").Return(nil)
	kv.On("PutMap", fmt.Sprintf(consul.KvClustersConnectionPath, "foobar"), expectedConnections).Return(nil)

	consulInst := new(mocks.Client)
	consulInst.On("KV").Return(kv)
	testLock := consulApi.Lock{}
	consulInst.On("AcquireLockKey", consul.KvClustersPath).Return(&testLock, nil).Times(2)

	deps := DefaultDependencies()
	deps.consul = consulInst

	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	data := url.Values{}
	data.Set("check_ids[]", "1.2.3")
	data.Set("username-host1", "myuser1")
	data.Set("username-host2", "myuser2")

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/clusters/foobar/settings", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "text/html")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	app.ServeHTTP(resp, req)

	assert.NoError(t, err)
	assert.Equal(t, 302, resp.Code)

	consulInst.AssertExpectations(t)
	kv.AssertExpectations(t)
}
