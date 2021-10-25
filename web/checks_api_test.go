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
		Hosts: map[string]*models.Host{
			"host1": &models.Host{
				Reachable: true,
				Msg:       "",
			},
			"host2": &models.Host{
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

	req, err := http.NewRequest("GET", "/api/clusters/47d1190ffb4f781974c8356d7f863b03/results", nil)
	if err != nil {
		t.Fatal(err)
	}

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

	req, err := http.NewRequest("GET", "/api/clusters/47d1190ffb4f781974c8356d7f863b03/results", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)
}

func TestApiCheckGetSelectedHandler(t *testing.T) {
	expectedValue := models.SelectedChecks{
		ID:             "group1",
		SelectedChecks: []string{"ABCDEF", "123456"},
	}

	mockChecksService := new(services.MockChecksService)
	mockChecksService.On(
		"GetSelectedChecksById", "group1").Return(expectedValue, nil)
	mockChecksService.On(
		"GetSelectedChecksById", "otherId").Return(models.SelectedChecks{}, fmt.Errorf("not found"))

	deps := setupTestDependencies()
	deps.checksService = mockChecksService

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	// 200 scenario
	resp := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/api/checks/group1/selected", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.webEngine.ServeHTTP(resp, req)

	var selectedChecks JSONSelectedChecks
	json.Unmarshal(resp.Body.Bytes(), &selectedChecks)

	assert.Equal(t, 200, resp.Code)
	assert.ElementsMatch(t, expectedValue.SelectedChecks, selectedChecks.SelectedChecks)

	// 404 scenario
	resp = httptest.NewRecorder()

	req, err = http.NewRequest("GET", "/api/checks/otherId/selected", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.webEngine.ServeHTTP(resp, req)

	json.Unmarshal(resp.Body.Bytes(), &selectedChecks)

	assert.Equal(t, 404, resp.Code)

	mockChecksService.AssertExpectations(t)
}

func TestApiCheckCreateSelectedHandler(t *testing.T) {
	mockChecksService := new(services.MockChecksService)
	mockChecksService.On(
		"CreateSelectedChecks", "group1", []string{"ABCDEF", "123456"}).Return(nil)
	mockChecksService.On(
		"CreateSelectedChecks", "otherId", []string{"ABCDEF", "123456"}).Return(fmt.Errorf("not storing"))

	deps := setupTestDependencies()
	deps.checksService = mockChecksService

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	// 200 scenario
	sentValue := JSONSelectedChecks{
		SelectedChecks: []string{"ABCDEF", "123456"},
	}
	resp := httptest.NewRecorder()
	body, _ := json.Marshal(&sentValue)
	req, err := http.NewRequest("POST", "/api/checks/group1/selected", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	app.webEngine.ServeHTTP(resp, req)

	var selectedChecks JSONSelectedChecks
	json.Unmarshal(resp.Body.Bytes(), &selectedChecks)

	assert.Equal(t, 201, resp.Code)
	assert.ElementsMatch(t, sentValue.SelectedChecks, selectedChecks.SelectedChecks)

	// 500 scenario
	resp = httptest.NewRecorder()

	req, err = http.NewRequest("POST", "/api/checks/otherId/selected", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)

	mockChecksService.AssertExpectations(t)
}
