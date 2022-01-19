package web

import (
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
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
		"a615a35f65627be5a757319a0741127f": map[string]interface{}{
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
				"name": "other_cluster",
				"id":   "a615a35f65627be5a757319a0741127f",
			},
		},
	}

	return listMap
}

func TestClustersListHandler(t *testing.T) {
	clustersList := models.ClusterList{
		{
			ID:                "47d1190ffb4f781974c8356d7f863b03",
			Name:              "hana_cluster",
			ClusterType:       models.ClusterTypeHANAScaleUp,
			SID:               "PRD",
			ResourcesNumber:   5,
			HostsNumber:       3,
			Tags:              []string{"tag1"},
			Health:            models.CheckPassing,
			HasDuplicatedName: false,
		},
		{
			ID:                "a615a35f65627be5a757319a0741127f",
			Name:              "other_cluster",
			ClusterType:       models.ClusterTypeUnknown,
			SID:               "",
			Tags:              []string{"tag1"},
			Health:            models.CheckCritical,
			HasDuplicatedName: false,
		},
		{
			ID:                "e2f2eb50aef748e586a7baa85e0162cf",
			Name:              "netweaver_cluster",
			ClusterType:       models.ClusterTypeUnknown,
			SID:               "",
			ResourcesNumber:   10,
			HostsNumber:       2,
			Tags:              []string{"tag1"},
			Health:            models.CheckCritical,
			HasDuplicatedName: true,
		},
		{
			ID:                "e27d313a674375b2066777a89ee346b9",
			Name:              "netweaver_cluster",
			ClusterType:       models.ClusterTypeUnknown,
			SID:               "",
			Tags:              []string{"tag1"},
			Health:            models.CheckUndefined,
			HasDuplicatedName: true,
		},
	}

	mockClusterService := new(services.MockClustersService)
	mockClusterService.On("GetAll", mock.Anything, mock.Anything).Return(clustersList, nil)
	mockClusterService.On("GetCount").Return(4, nil)
	mockClusterService.On("GetAllClusterNames", mock.Anything).Return(
		[]string{"hana_cluster", "other_cluster", "netweaver_cluster"},
		nil,
	)
	mockClusterService.On("GetAllClusterTypes", mock.Anything).Return(
		[]string{models.ClusterTypeHANAScaleUp, models.ClusterTypeUnknown},
		nil,
	)
	mockClusterService.On("GetAllSIDs", mock.Anything).Return([]string{"PRD"}, nil)
	mockClusterService.On("GetAllTags", mock.Anything).Return([]string{"tag1"}, nil)
	deps := setupTestDependencies()
	deps.clustersService = mockClusterService

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/clusters", nil)

	app.webEngine.ServeHTTP(resp, req)

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
	assert.Regexp(t, regexp.MustCompile("<td .*>.*error.*</td><td>.*other_cluster.*</td><td>.*a615a35f65627be5a757319a0741127f.*</td><td>Unknown</td><td></td>"), minified)
	assert.Regexp(t, regexp.MustCompile("(?s)<td .*>.*error.*</td><td>.*duplicated.*netweaver_cluster.*</td><td>.*e2f2eb50aef748e586a7baa85e0162cf.*</td><td>Unknown</td><td></td><td>2</td><td>10</td><td><input.*value=tag1.*></td>"), minified)
	assert.Regexp(t, regexp.MustCompile("(?s)<td .*>.*fiber_manual_record.*</td><td>.*duplicated.*info.*netweaver_cluster.*</td><td>.*e27d313a674375b2066777a89ee346b9.*</td><td>Unknown</td><td></td>"), minified)
}

func TestClusterHandlerHANA(t *testing.T) {
	clusterID := "47d1190ffb4f781974c8356d7f863b03"

	clustersService := new(services.MockClustersService)
	clustersService.On("GetByID", clusterID).Return(&models.Cluster{
		ID:            clusterID,
		Name:          "hana_cluster",
		ClusterType:   models.ClusterTypeHANAScaleUp,
		SID:           "PRD",
		Tags:          []string{"tag1"},
		Health:        models.CheckCritical,
		PassingCount:  2,
		CriticalCount: 1,
		Details: &models.HANAClusterDetails{
			SystemReplicationMode:          "sync",
			SystemReplicationOperationMode: "logreplay",
			SecondarySyncState:             "SFAIL",
			SRHealthState:                  "1",
			FencingType:                    "external/sbd",
			CIBLastWritten:                 time.Date(2021, time.June, 30, 18, 11, 37, 0, time.UTC),
			StoppedResources: []*models.ClusterResource{
				{
					ID:        "dummy_failed",
					Type:      "dummy",
					Role:      "Started",
					Status:    "failed",
					FailCount: 0,
				},
			},
			Nodes: []*models.HANAClusterNode{
				{
					HostID:      "host1",
					Name:        "test_node_1",
					IPAddresses: []string{"192.168.1.1"},
					VirtualIPs:  []string{"10.123.123.123"},
					HANAStatus:  "Primary",
					Health:      models.HostHealthPassing,
					Resources: []*models.ClusterResource{
						{
							ID:        "dummy_failed",
							Type:      "dummy",
							Role:      "Started",
							Status:    "failed",
							FailCount: 0,
						},
						{
							ID:        "sbd",
							Type:      "stonith:external/sbd",
							Role:      "Started",
							Status:    "active",
							FailCount: 0,
						},
					},
				},
				{
					HostID:      "host2",
					Name:        "test_node_2",
					IPAddresses: []string{"192.168.1.2"},
					HANAStatus:  "Failed",
					Health:      models.HostHealthCritical,
				},
			},
		},
	}, nil)

	deps := setupTestDependencies()
	deps.clustersService = clustersService

	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/clusters/"+clusterID, nil)
	req.Header.Set("Accept", "text/html")

	app.webEngine.ServeHTTP(resp, req)

	clustersService.AssertExpectations(t)

	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepEndTags:         true,
	})
	minified, err := m.String("text/html", resp.Body.String())
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "Cluster details")
	// Summary
	assert.Regexp(t, regexp.MustCompile("<strong>SID:</strong><br><span.*>PRD</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<strong>Cluster name:</strong><br><span.*>hana_cluster</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<strong>Cluster type:</strong><br><span.*>HANA scale-up</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<strong>HANA system replication mode:</strong><br><span.*>sync</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<strong>Fencing type:</strong><br><span.*>external/sbd</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<strong>HANA system replication operation mode:</strong><br><span.*>logreplay</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<strong>CIB last written:</strong><br><span.*>Jun 30, 2021 18:11:37 UTC</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<strong>SAPHanaSR health state:</strong>.*text-danger.*"), minified)
	assert.Regexp(t, regexp.MustCompile("<strong>HANA secondary sync state:</strong><br><span.*>SFAIL</span>"), minified)
	// Health
	assert.Regexp(t, regexp.MustCompile(".*check_circle.*alert-body.*Passing.*2"), minified)
	assert.Regexp(t, regexp.MustCompile(".*error.*alert-body.*Critical.*1"), minified)

	// Nodes
	assert.Regexp(t, regexp.MustCompile("<td.*check_circle.*<td.*><a.*href=/hosts/host1.*>test_node_1</a></td><td.*>192\\.168\\.1\\.1</td><td.*>10\\.123\\.123\\.123</td><td.*><span .*>HANA Primary</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td.*error.*<td.*><a.*href=/hosts/host2.*>test_node_2</a></td><td.*>192\\.168\\.1\\.2</td>.*<span .*danger.*>HANA Failed</span>"), minified)
	// Resources
	assert.Regexp(t, regexp.MustCompile("<td>sbd</td><td>stonith:external/sbd</td><td>Started</td><td>active</td><td>0</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>dummy_failed</td><td>dummy</td><td>Started</td><td>failed</td><td>0</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<h4>Stopped resources</h4><div.*><div.*><span .*>dummy_failed</span>"), minified)
}
