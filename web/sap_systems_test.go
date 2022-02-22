package web

import (
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"

	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

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

	responseBody := resp.Body.String()

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, responseBody, "SAP Systems")
	assert.Contains(t, responseBody, "HA1")
	assert.Contains(t, responseBody, "PRD")
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
	sapSystemsService.On("GetAllDatabasesSIDs").Return([]string{"PRD"}, nil)
	sapSystemsService.On("GetAllDatabasesTags").Return([]string{"tag1"}, nil)

	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/databases", nil)

	app.webEngine.ServeHTTP(resp, req)
	sapSystemsService.AssertExpectations(t)

	responseBody := resp.Body.String()

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, responseBody, "HANA Databases")
	assert.Contains(t, responseBody, "PRD")
	assert.Contains(t, responseBody, "HDB_WORKER")
}

func TestSAPResourceHandler(t *testing.T) {
	sapSystemsService := new(services.MockSAPSystemsService)
	hostsService := new(services.MockHostsService)

	sapSystemsService.On("GetByID", "sap_system_id").Return(&models.SAPSystem{
		ID:   "sap_system_id",
		SID:  "PRD",
		Type: models.SAPSystemTypeApplication,
		Instances: []*models.SAPSystemInstance{
			{
				InstanceNumber: "00",
				SAPHostname:    "netweaver01",
				Features:       "MESSAGESERVER|ENQUE",
				HttpPort:       50013,
				HttpsPort:      50014,
				Status:         "SAPControl-GREEN",
				StartPriority:  "0.5",
			},
		},
	}, nil)
	hostsService.On("GetAllBySAPSystemID", "sap_system_id").Return(models.HostList{
		{
			ID:            "netweaver01",
			Name:          "netweaver01",
			AgentVersion:  "v0",
			IPAddresses:   []string{"192.168.10.10"},
			Health:        "passing",
			CloudProvider: "azure",
			ClusterID:     "cluster_id",
			ClusterName:   "netweaver",
			ClusterType:   models.ClusterTypeHANAScaleOut,
			SAPSystems: []*models.SAPSystem{
				{
					ID:  "sap_system_id",
					SID: "PRD",
				},
			},
		},
	}, nil)

	deps := setupTestDependencies()
	deps.sapSystemsService = sapSystemsService
	deps.hostsService = hostsService

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/sapsystems/sap_system_id", nil)

	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	responseBody := minifyHtml(resp.Body.String())

	assert.Contains(t, responseBody, "SAP System details")
	assert.Contains(t, responseBody, "PRD")
	// Layout
	assert.Regexp(t, regexp.MustCompile("<tr><td>netweaver01</td><td>00</td><td>MESSAGESERVER\\|ENQUE</td><td>50013</td><td>50014</td><td>0.5</td><td><span.*primary.*>SAPControl-GREEN</span></td></tr>"), responseBody)
	// Host
	assert.Regexp(t, regexp.MustCompile("<tr><td>.*check_circle.*</td><td .*><a href=/hosts/netweaver01>netweaver01</a></td><td>192.168.10.10</td><td>azure</td><td><a href=/clusters/cluster_id>netweaver</a></td><td>v0</td></tr>"), responseBody)
}

func TestSAPResourceHandler404Error(t *testing.T) {
	sapSystemsService := new(services.MockSAPSystemsService)
	sapSystemsService.On("GetByID", mock.Anything).Return(nil, nil)

	deps := setupTestDependencies()
	deps.sapSystemsService = sapSystemsService

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
