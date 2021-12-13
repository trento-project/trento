package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

func TestApiClusterCheckResultsHandler(t *testing.T) {
	results := &models.ClusterCheckResults{
		Hosts: map[string]*models.CheckHost{
			"host1": &models.CheckHost{
				Reachable: true,
				Msg:       "",
			},
			"host2": &models.CheckHost{
				Reachable: false,
				Msg:       "error connecting",
			},
		},
		Checks: []models.ClusterCheckResult{
			models.ClusterCheckResult{
				ID:          "ABCDEF",
				Group:       "group 1",
				Description: "description 1",
				Hosts: map[string]*models.Check{
					"host1": &models.Check{
						Result: models.CheckPassing,
						Msg:    "some random message",
					},
					"host2": &models.Check{
						Result: models.CheckPassing,
					},
				},
			},
			models.ClusterCheckResult{
				ID:          "123456",
				Group:       "group 1",
				Description: "description 2",
				Hosts: map[string]*models.Check{
					"host1": &models.Check{
						Result: models.CheckWarning,
					},
					"host2": &models.Check{
						Result: models.CheckCritical,
					},
				},
			},
		},
	}

	mockChecksService := new(services.MockChecksService)
	mockChecksService.On(
		"GetChecksResultAndMetadataByCluster", "47d1190ffb4f781974c8356d7f863b03").Return(results, nil)

	deps := setupTestDependencies()
	deps.checksService = mockChecksService

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/api/clusters/47d1190ffb4f781974c8356d7f863b03/results", nil)

	app.webEngine.ServeHTTP(resp, req)

	expectedBody, _ := json.Marshal(gin.H{
		"hosts": gin.H{
			"host1": gin.H{
				"reachable": true,
				"msg":       "",
			},
			"host2": gin.H{
				"reachable": false,
				"msg":       "error connecting",
			},
		},
		"checks": []gin.H{
			gin.H{
				"id":          "ABCDEF",
				"group":       "group 1",
				"description": "description 1",
				"hosts": gin.H{
					"host1": gin.H{
						"result": "passing",
						"msg":    "some random message",
					},
					"host2": gin.H{
						"result": "passing",
					},
				},
			},
			gin.H{
				"id":          "123456",
				"group":       "group 1",
				"description": "description 2",
				"hosts": gin.H{
					"host1": gin.H{
						"result": "warning",
					},
					"host2": gin.H{
						"result": "critical",
					},
				},
			},
		},
	})
	assert.JSONEq(t, string(expectedBody), resp.Body.String())
	assert.Equal(t, 200, resp.Code)
}

func TestApiClusterCheckResultsHandler500(t *testing.T) {
	mockChecksService := new(services.MockChecksService)
	mockChecksService.On(
		"GetChecksResultAndMetadataByCluster", "47d1190ffb4f781974c8356d7f863b03").Return(
		&models.ClusterCheckResults{}, fmt.Errorf("kaboom"))

	deps := setupTestDependencies()
	deps.checksService = mockChecksService

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/api/clusters/47d1190ffb4f781974c8356d7f863b03/results", nil)

	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)
}

func TestApiCreateChecksCatalogHandler(t *testing.T) {
	expectedCatalog := models.ChecksCatalog{
		&models.Check{
			ID:             "id1",
			Name:           "name1",
			Group:          "group1",
			Description:    "description1",
			Remediation:    "remediation1",
			Implementation: "implementation1",
			Labels:         "labels1",
			Premium:        true,
		},
		&models.Check{
			ID:             "id2",
			Name:           "name2",
			Group:          "group2",
			Description:    "description2",
			Remediation:    "remediation2",
			Implementation: "implementation2",
			Labels:         "labels2",
		},
	}
	mockChecksService := new(services.MockChecksService)
	mockChecksService.On("CreateChecksCatalog", expectedCatalog).Return(nil)
	mockChecksService.On("CreateChecksCatalog", models.ChecksCatalog(nil)).Return(fmt.Errorf("error"))

	deps := setupTestDependencies()
	deps.checksService = mockChecksService

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	// 200 scenario
	sendData := JSONChecksCatalog{
		&JSONCheck{
			ID:             "id1",
			Name:           "name1",
			Group:          "group1",
			Description:    "description1",
			Remediation:    "remediation1",
			Implementation: "implementation1",
			Labels:         "labels1",
			Premium:        true,
		},
		&JSONCheck{
			ID:             "id2",
			Name:           "name2",
			Group:          "group2",
			Description:    "description2",
			Remediation:    "remediation2",
			Implementation: "implementation2",
			Labels:         "labels2",
		},
	}

	resp := httptest.NewRecorder()
	body, _ := json.Marshal(&sendData)
	req := httptest.NewRequest("PUT", "/api/checks/catalog", bytes.NewBuffer(body))

	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	// 500 scenario
	resp = httptest.NewRecorder()

	sendData = JSONChecksCatalog{}
	body, _ = json.Marshal(&sendData)
	req = httptest.NewRequest("PUT", "/api/checks/catalog", bytes.NewBuffer(body))

	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)

	mockChecksService.AssertExpectations(t)
}

func TestApiCheckGetSettingsByIdHandler(t *testing.T) {
	mockClustersService := new(services.MockClustersService)
	mockClustersService.On("GetClusterSettingsByID", "a615a35f65627be5a757319a0741127f").Return(&models.ClusterSettings{
		SelectedChecks: []string{"ABCDEF", "123456"},
		Hosts: []*models.HostConnection{
			{
				Name: "host1",
				User: "user1",
			},
			{
				Name: "host2",
				User: "user2",
			},
		},
	}, nil)

	deps := setupTestDependencies()
	deps.clustersService = mockClustersService

	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/checks/a615a35f65627be5a757319a0741127f/settings", nil)
	assert.NoError(t, err)

	app.webEngine.ServeHTTP(resp, req)

	var settings *JSONChecksSettings
	json.Unmarshal(resp.Body.Bytes(), &settings)

	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, []string{"ABCDEF", "123456"}, settings.SelectedChecks)
	assert.Equal(t, map[string]string{
		"host1": "user1",
		"host2": "user2",
	}, settings.ConnectionSettings)
	assert.Equal(t, []string{"host1", "host2"}, settings.Hostnames)
}

func TestApiCheckGetSettingsByIdHandler404(t *testing.T) {
	mockClustersService := new(services.MockClustersService)
	mockClustersService.On("GetClusterSettingsByID", "not_found").Return(nil, nil)

	deps := setupTestDependencies()
	deps.clustersService = mockClustersService

	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/checks/not_found/settings", nil)
	assert.NoError(t, err)

	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 404, resp.Code)
}

func TestApiCheckCreateConnectionByIdHandler(t *testing.T) {
	mockChecksService := new(services.MockChecksService)

	mockChecksService.On(
		"CreateSelectedChecks", "group1", []string{"ABCDEF", "123456"}).Return(nil)
	mockChecksService.On(
		"CreateSelectedChecks", "otherId", []string{"ABCDEF", "123456"}).Return(fmt.Errorf("not storing"))

	mockChecksService.On(
		"CreateConnectionSettings", "group1", "node1", "user1").Return(nil)
	mockChecksService.On(
		"CreateConnectionSettings", "group1", "node2", "user2").Return(nil)

	deps := setupTestDependencies()
	deps.checksService = mockChecksService

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	// 200 scenario
	sendData := JSONChecksSettings{
		SelectedChecks: []string{"ABCDEF", "123456"},
		ConnectionSettings: map[string]string{
			"node1": "user1",
			"node2": "user2",
		},
	}
	resp := httptest.NewRecorder()
	body, _ := json.Marshal(&sendData)
	req := httptest.NewRequest("POST", "/api/checks/group1/settings", bytes.NewBuffer(body))

	app.webEngine.ServeHTTP(resp, req)

	var connData JSONChecksSettings
	json.Unmarshal(resp.Body.Bytes(), &connData)

	assert.Equal(t, 201, resp.Code)
	assert.Equal(t, sendData, connData)

	// 500 scenario
	resp = httptest.NewRecorder()

	req = httptest.NewRequest("POST", "/api/checks/otherId/settings", bytes.NewBuffer(body))

	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)

	mockChecksService.AssertExpectations(t)
}
