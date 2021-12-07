package web

import (
	"net/http/httptest"
	"regexp"
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/internal/sapsystem"
	"github.com/trento-project/trento/internal/sapsystem/sapcontrol"

	consulMocks "github.com/trento-project/trento/internal/consul/mocks"
	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

var sapSystemsList = sapsystem.SAPSystemsList{
	&sapsystem.SAPSystem{
		Id:   "systemId1",
		SID:  "HA1",
		Type: sapsystem.Application,
		Instances: map[string]*sapsystem.SAPInstance{
			"ASCS00": &sapsystem.SAPInstance{
				Host: "netweaver01",
				SAPControl: &sapsystem.SAPControl{
					Properties: map[string]*sapcontrol.InstanceProperty{
						"SAPSYSTEMNAME": &sapcontrol.InstanceProperty{
							Property:     "SAPSYSTEMNAME",
							Propertytype: "string",
							Value:        "HA1",
						},
						"SAPSYSTEM": &sapcontrol.InstanceProperty{
							Property:     "SAPSYSTEM",
							Propertytype: "string",
							Value:        "00",
						},
					},
					Instances: map[string]*sapcontrol.SAPInstance{
						"netweaver01": &sapcontrol.SAPInstance{
							Hostname:      "netweaver01",
							InstanceNr:    0,
							Features:      "MESSAGESERVER|ENQUE",
							HttpPort:      50013,
							HttpsPort:     50014,
							StartPriority: "0.5",
							Dispstatus:    "SAPControl-GREEN",
						},
						"netweaver02": &sapcontrol.SAPInstance{
							Hostname:   "netweaver02",
							InstanceNr: 10,
							Features:   "ENQREP",
						},
					},
				},
			},
		},
		Profile: map[string]interface{}{
			"dbs": map[string]interface{}{
				"hdb": map[string]interface{}{
					"dbname": "PRD",
				},
			},
			"SAPDBHOST": "192.168.1.5",
		},
	},
	&sapsystem.SAPSystem{
		Id:   "systemId1",
		SID:  "HA1",
		Type: sapsystem.Application,
		Instances: map[string]*sapsystem.SAPInstance{
			"ERS10": &sapsystem.SAPInstance{
				Host: "netweaver02",
				SAPControl: &sapsystem.SAPControl{
					Properties: map[string]*sapcontrol.InstanceProperty{
						"SAPSYSTEMNAME": &sapcontrol.InstanceProperty{
							Property:     "SAPSYSTEMNAME",
							Propertytype: "string",
							Value:        "HA1",
						},
						"SAPSYSTEM": &sapcontrol.InstanceProperty{
							Property:     "SAPSYSTEM",
							Propertytype: "string",
							Value:        "10",
						},
					},
					Instances: map[string]*sapcontrol.SAPInstance{
						"netweaver01": &sapcontrol.SAPInstance{
							Hostname:   "netweaver01",
							InstanceNr: 0,
							Features:   "MESSAGESERVER|ENQUE",
						},
						"netweaver02": &sapcontrol.SAPInstance{
							Hostname:   "netweaver02",
							InstanceNr: 10,
							Features:   "ENQREP",
						},
					},
				},
			},
		},
	},
	// Test duplicated icon
	&sapsystem.SAPSystem{
		Id:        "systemId2",
		SID:       "DEV",
		Type:      sapsystem.Application,
		Instances: map[string]*sapsystem.SAPInstance{},
	},
	&sapsystem.SAPSystem{
		Id:        "systemId3",
		SID:       "DEV",
		Type:      sapsystem.Application,
		Instances: map[string]*sapsystem.SAPInstance{},
	},
}

var sapDatabasesList = sapsystem.SAPSystemsList{
	&sapsystem.SAPSystem{
		Id:   "systemId2",
		SID:  "PRD",
		Type: sapsystem.Database,
		Instances: map[string]*sapsystem.SAPInstance{
			"HDB00": &sapsystem.SAPInstance{
				Host: "hana01",
				Type: sapsystem.Database,
				SAPControl: &sapsystem.SAPControl{
					Properties: map[string]*sapcontrol.InstanceProperty{
						"SAPSYSTEMNAME": &sapcontrol.InstanceProperty{
							Property:     "SAPSYSTEMNAME",
							Propertytype: "string",
							Value:        "PRD",
						},
						"SAPSYSTEM": &sapcontrol.InstanceProperty{
							Property:     "SAPSYSTEM",
							Propertytype: "string",
							Value:        "00",
						},
					},
					Instances: map[string]*sapcontrol.SAPInstance{
						"hana01": &sapcontrol.SAPInstance{
							Hostname:   "hana01",
							InstanceNr: 0,
							Features:   "HDB_WORKER",
						},
					},
				},
				SystemReplication: sapsystem.SystemReplication{
					"local_site_id": "1",
					"site": map[string]interface{}{
						"1": map[string]interface{}{
							"REPLICATION_MODE": "PRIMARY",
						},
					},
					"overall_replication_status": "ACTIVE",
				},
			},
		},
	},
}

func TestSAPSystemsListHandler(t *testing.T) {
	sapSystemsService := new(services.MockSAPSystemsService)

	deps := setupTestDependencies()
	deps.sapSystemsService = sapSystemsService
	sapSystemsService.On("GetAllApplications", mock.Anything, mock.Anything).Return(models.SAPSystemList{
		{
			ID:     "application_id",
			SID:    "HA1",
			Type:   models.SAPSystemTypeApplication,
			Tags:   []string{"tag1"},
			DBName: "PRD",
			DBHost: "192.168.1.5",
			Instances: []*models.SAPSystemInstance{
				{
					InstanceNumber: "00",
					SID:            "HA1",
					Features:       "MESSAGESERVER|ENQUE",
					HostID:         "host_id_1",
					Hostname:       "netweaver01",
					ClusterName:    "netweaver_cluster",
					ClusterID:      "cluster_id",
					Type:           "application",
				},
				{
					InstanceNumber: "10",
					SID:            "HA1",
					Features:       "ENQREP",
					HostID:         "host_id_2",
					Hostname:       "netweaver02",
					ClusterName:    "netweaver_cluster",
					ClusterID:      "cluster_id",
					Type:           "application",
				},
			},
			AttachedDatabase: &models.SAPSystem{
				ID:  "database_id",
				SID: "PRD",
				Instances: []*models.SAPSystemInstance{
					{
						InstanceNumber:          "00",
						SID:                     "PRD",
						Features:                "HDB_WORKER",
						HostID:                  "host_id_3",
						Hostname:                "hana01",
						ClusterName:             "hana_cluster",
						ClusterID:               "cluster_id_2",
						SystemReplication:       "Primary",
						SystemReplicationStatus: "SOK",
						Type:                    "database",
					},
				},
			},
		},
		{
			ID:               "duplicated_sid_1",
			SID:              "DEV",
			Type:             models.SAPSystemTypeApplication,
			HasDuplicatedSID: true,
			AttachedDatabase: &models.SAPSystem{},
		},
		{
			ID:               "duplicated_sid_2",
			SID:              "DEV",
			Type:             models.SAPSystemTypeApplication,
			HasDuplicatedSID: true,
			AttachedDatabase: &models.SAPSystem{},
		},
	}, nil)
	sapSystemsService.On("GetApplicationsCount").Return(1, nil)
	sapSystemsService.On("GetAllApplicationsSIDs").Return([]string{"HA1"}, nil)
	sapSystemsService.On("GetAllApplicationsTags").Return([]string{"tag1"}, nil)

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/sapsystems", nil)

	app.webEngine.ServeHTTP(resp, req)
	sapSystemsService.AssertExpectations(t)

	responseBody := minifyHtml(resp.Body.String())

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, responseBody, "SAP Systems")
	assert.Regexp(t, regexp.MustCompile("<a href=/sapsystems/application_id>HA1</a></td><td></td><td><a href=/databases/database_id>PRD</a></td><td>PRD</td><td>192.168.1.5</td><td>.*<input.*value=tag1.*>.*</td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<td>HA1</td><td>MESSAGESERVER\\|ENQUE</td><td>00</td><td></td><td><a href=/clusters/cluster_id>netweaver_cluster</a></td><td><a href=/hosts/host_id_1>netweaver01</a></td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<td>HA1</td><td>ENQREP</td><td>10</td><td></td><td><a href=/clusters/cluster_id>netweaver_cluster</a></td><td><a href=/hosts/host_id_2>netweaver02</a></td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("(?s)<td>PRD</td><td>HDB_WORKER</td><td>00</td><td>HANA Primary.*SOK.*</td><td><a href=/clusters/cluster_id_2>hana_cluster</a></td><td><a href=/hosts/host_id_3>hana01</a></td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("(?s)<td><i .*This SAP system SID exists multiple times.*warning.*<a href=/sapsystems/duplicated_sid_1>DEV</a></td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("(?s)<td><i .*This SAP system SID exists multiple times.*warning.*<a href=/sapsystems/duplicated_sid_2>DEV</a></td>"), responseBody)
}

func TestSAPDatabaseListHandler(t *testing.T) {
	sapSystemsService := new(services.MockSAPSystemsService)

	deps := setupTestDependencies()
	deps.sapSystemsService = sapSystemsService
	sapSystemsService.On("GetAllDatabases", mock.Anything, mock.Anything).Return(models.SAPSystemList{
		{
			ID:   "database_id",
			SID:  "PRD",
			Type: models.SAPSystemTypeDatabase,
			Tags: []string{"tag1"},
			Instances: []*models.SAPSystemInstance{
				{
					InstanceNumber:          "00",
					Features:                "HDB_WORKER",
					HostID:                  "host_id",
					Hostname:                "hana01",
					ClusterName:             "hana_cluster",
					ClusterID:               "cluster_id",
					SystemReplication:       "Primary",
					SystemReplicationStatus: "SOK",
					SID:                     "PRD",
					Type:                    models.SAPSystemTypeDatabase,
				},
			},
		},
	}, nil)
	sapSystemsService.On("GetDatabasesCount").Return(1, nil)
	sapSystemsService.On("GetAllDatabasesSIDs").Return([]string{"PRD"}, nil)
	sapSystemsService.On("GetAllDatabasesTags").Return([]string{"tag1"}, nil)

	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/databases", nil)

	app.webEngine.ServeHTTP(resp, req)
	sapSystemsService.AssertExpectations(t)

	responseBody := minifyHtml(resp.Body.String())

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, responseBody, "HANA Databases")

	assert.Regexp(t, regexp.MustCompile("<td><a href=/databases/database_id>PRD</a></td><td></td><td><input class=tags-input value=tag1"), responseBody)
	assert.Regexp(t, regexp.MustCompile("(?s)<td>PRD</td><td>HDB_WORKER</td><td>00</td>.*HANA Primary.*SOK.*<td><a href=/clusters/cluster_id>hana_cluster</a></td><td><a href=/hosts/host_id>hana01</a></td>"), responseBody)
}

func TestSAPResourceHandler(t *testing.T) {
	consulInst := new(consulMocks.Client)
	health := new(consulMocks.Health)
	consulInst.On("Health").Return(health)
	sapSystemsService := new(services.MockSAPSystemsConsulService)
	hostsService := new(services.MockHostsConsulService)

	deps := setupTestDependencies()
	deps.consul = consulInst
	deps.sapSystemsConsulService = sapSystemsService
	deps.hostsConsulService = hostsService

	host := hosts.NewHost(consulApi.Node{
		Node:    "netweaver01",
		Address: "192.168.10.10",
		Meta: map[string]string{
			"trento-sap-systems":      "PRD",
			"trento-sap-systems-type": "Application",
			"trento-sap-systems-id":   "systemId",
			"trento-cloud-provider":   "azure",
			"trento-agent-version":    "0",
			"trento-ha-cluster-id":    "e2f2eb50aef748e586a7baa85e0162cf",
			"trento-ha-cluster":       "banana",
		},
	},
		consulInst)
	hostList := hosts.HostList{
		&host,
	}

	passHealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthPassing,
		},
	}

	health.On("Node", "netweaver01", (*consulApi.QueryOptions)(nil)).Return(passHealthChecks, nil, nil)
	sapSystemsService.On("GetSAPSystemsById", "systemId").Return(sapSystemsList, nil)
	hostsService.On("GetHostsBySystemId", "systemId").Return(hostList, nil)

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/sapsystems/systemId", nil)

	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	responseBody := minifyHtml(resp.Body.String())

	sapSystemsService.AssertExpectations(t)
	hostsService.AssertExpectations(t)
	consulInst.AssertExpectations(t)

	assert.Contains(t, responseBody, "SAP System details")
	assert.Contains(t, responseBody, "PRD")
	// Layout
	assert.Regexp(t, regexp.MustCompile("<tr><td>netweaver01</td><td>00</td><td>MESSAGESERVER\\|ENQUE</td><td>50013</td><td>50014</td><td>0.5</td><td><span.*primary.*>SAPControl-GREEN</span></td></tr>"), responseBody)
	// Host
	assert.Regexp(t, regexp.MustCompile("<tr><td>.*check_circle.*</td><td><a href=/hosts/netweaver01>netweaver01</a></td><td>192.168.10.10</td><td>azure</td><td><a href=/clusters/e2f2eb50aef748e586a7baa85e0162cf>banana</a></td><td><a href=/sapsystems/systemId>PRD</a></td><td>v0</td></tr>"), responseBody)
}

func TestSAPResourceHandler404Error(t *testing.T) {
	sapSystemsService := new(services.MockSAPSystemsConsulService)

	deps := setupTestDependencies()
	deps.sapSystemsConsulService = sapSystemsService

	sapSystemsService.On("GetSAPSystemsById", "foobar").Return(sapsystem.SAPSystemsList{}, nil)

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/sapsystems/foobar", nil)
	req.Header.Set("Accept", "text/html")

	app.webEngine.ServeHTTP(resp, req)

	sapSystemsService.AssertExpectations(t)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.Code)
	assert.Contains(t, resp.Body.String(), "Not Found")
}

func minifyHtml(input string) string {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepEndTags:         true,
	})
	minified, err := m.String("text/html", input)
	if err != nil {
		panic(err)
	}
	return minified
}
