package web

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	consulApi "github.com/hashicorp/consul/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/consul/mocks"
)

var clustersListMap = map[string]interface{}{
	"47d1190ffb4f781974c8356d7f863b03": map[string]interface{}{
		"id": "47d1190ffb4f781974c8356d7f863b03",
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
				"CrmConfig": map[string]interface{}{
					"ClusterProperties": []interface{}{
						map[string]interface{}{
							"Id":    "cib-bootstrap-options-cluster-name",
							"Value": "hana_cluster",
						},
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
	},
	"e2f2eb50aef748e586a7baa85e0162cf": map[string]interface{}{
		"id": "e2f2eb50aef748e586a7baa85e0162cf",
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
	},
	"e27d313a674375b2066777a89ee346b9": map[string]interface{}{
		"id": "e27d313a674375b2066777a89ee346b9",
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
	},
}

func TestClustersListHandler(t *testing.T) {
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

	node1HealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthPassing,
		},
	}

	node2HealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthCritical,
		},
	}

	consulInst := new(mocks.Client)

	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)
	kv.On("ListMap", consul.KvClustersPath, consul.KvClustersPath).Return(clustersListMap, nil)
	consulInst.On("WaitLock", consul.KvClustersPath).Return(nil)

	catalog := new(mocks.Catalog)
	catalog.On("Nodes", mock.Anything).Return(nodes, nil, nil)
	consulInst.On("Catalog").Return(catalog)

	health := new(mocks.Health)
	health.On("Node", "test_node_1", (*consulApi.QueryOptions)(nil)).Return(node1HealthChecks, nil, nil)
	health.On("Node", "test_node_2", (*consulApi.QueryOptions)(nil)).Return(node2HealthChecks, nil, nil)
	consulInst.On("Health").Return(health)

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
	assert.Regexp(t, regexp.MustCompile("<td .*>.*error.*</td><td>.*hana_cluster.*</td><td>.*47d1190ffb4f781974c8356d7f863b03.*</td><td>HANA scale-up</td><td>PRD</td><td>3</td><td>5</td><td>.*<option.*>tag1</option>.*</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td .*>.*fiber_manual_record.*</td><td>.*duplicated.*netweaver_cluster.*</td><td>.*e2f2eb50aef748e586a7baa85e0162cf.*</td><td>Unknown</td><td></td><td>.*<option.*>tag1</option>.*</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td .*>.*fiber_manual_record.*</td><td>.*duplicated.*netweaver_cluster.*</td><td>.*e27d313a674375b2066777a89ee346b9.*</td><td>Unknown</td><td></td>"), minified)
}

func TestClusterHandlerHANA(t *testing.T) {
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

	node1HealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthPassing,
		},
	}

	node2HealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthCritical,
		},
	}

	consulInst := new(mocks.Client)

	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)
	kv.On("ListMap", consul.KvClustersPath, consul.KvClustersPath).Return(clustersListMap, nil)
	consulInst.On("WaitLock", consul.KvClustersPath).Return(nil)

	catalog := new(mocks.Catalog)
	filter := &consulApi.QueryOptions{Filter: "Meta[\"trento-ha-cluster-id\"] == \"47d1190ffb4f781974c8356d7f863b03\""}
	catalog.On("Nodes", filter).Return(nodes, nil, nil)
	consulInst.On("Catalog").Return(catalog)

	health := new(mocks.Health)
	health.On("Node", "test_node_1", (*consulApi.QueryOptions)(nil)).Return(node1HealthChecks, nil, nil)
	health.On("Node", "test_node_2", (*consulApi.QueryOptions)(nil)).Return(node2HealthChecks, nil, nil)
	consulInst.On("Health").Return(health)

	deps := DefaultDependencies()
	deps.consul = consulInst

	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/clusters/47d1190ffb4f781974c8356d7f863b03", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

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
	// Nodes
	assert.Regexp(t, regexp.MustCompile("<td><a.*href=/hosts/test_node_1.*>test_node_1</a></td><td>192.168.1.1</td><td>10.123.123.123</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td><a.*href=/hosts/test_node_2.*>test_node_2</a></td><td>192.168.1.2</td>"), minified)
	// Resources
	assert.Regexp(t, regexp.MustCompile("<td>sbd</td><td>stonith:external/sbd</td><td>Started</td><td>active</td><td>0</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>dummy_failed</td><td>dummy</td><td>Started</td><td>failed</td><td>0</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<h4>Stopped resources</h4><div.*><div.*><span .*>dummy_failed</span>"), minified)
}

func TestClusterHandlerGeneric(t *testing.T) {
	consulInst := new(mocks.Client)

	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)

	kv.On("ListMap", consul.KvClustersPath, consul.KvClustersPath).Return(clustersListMap, nil)
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
	kv.On("ListMap", consul.KvClustersPath, consul.KvClustersPath).Return(clustersListMap, nil)

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
