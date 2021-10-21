package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	consulApi "github.com/hashicorp/consul/api"

	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/sapsystem"

	"github.com/stretchr/testify/mock"
	consulMocks "github.com/trento-project/trento/internal/consul/mocks"

	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"

	"github.com/stretchr/testify/assert"
)

func TestApiListTag(t *testing.T) {
	tagsSAPSystems := []string{}

	tagsClusters := []string{
		"tag1",
		"tag2",
		"tag3",
	}

	tagsHosts := []string{
		"tag4",
		"tag5",
		"tag6",
	}

	tags := append(tagsSAPSystems, tagsClusters...)
	tags = append(tags, tagsHosts...)

	mockTagsService := new(services.MockTagsService)
	mockTagsService.On("GetAll").Return(tags, nil)
	mockTagsService.On("GetAll", "sapsystems").Return(tagsSAPSystems, nil)
	mockTagsService.On("GetAll", "hosts").Return(tagsHosts, nil)
	deps := setupTestDependencies()
	deps.tagsService = mockTagsService

	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/tags", nil)
	app.webEngine.ServeHTTP(resp, req)

	expectedBody, _ := json.Marshal(tags)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, expectedBody, resp.Body.Bytes())

	resp = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/tags?resource_type=hosts", nil)
	app.webEngine.ServeHTTP(resp, req)

	expectedBody, _ = json.Marshal([]string{
		"tag4",
		"tag5",
		"tag6",
	})
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, expectedBody, resp.Body.Bytes())

	resp = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/tags?resource_type=sapsystems", nil)
	app.webEngine.ServeHTTP(resp, req)

	expectedBody, _ = json.Marshal([]string{})
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, expectedBody, resp.Body.Bytes())
}

func setupTestApiHostTag(resourceId string) Dependencies {
	node := &consulApi.Node{
		Node: resourceId,
	}

	consulInst := new(consulMocks.Client)
	catalog := new(consulMocks.Catalog)
	catalogNode := &consulApi.CatalogNode{Node: node}
	catalog.On("Node", resourceId, (*consulApi.QueryOptions)(nil)).Return(catalogNode, nil, nil)
	catalog.On("Node", mock.Anything, (*consulApi.QueryOptions)(nil)).Return(nil, nil, nil)
	consulInst.On("Catalog").Return(catalog)

	deps := setupTestDependencies()
	deps.consul = consulInst

	return deps
}

func setupTestApiClusterTag(resourceID string) Dependencies {
	clustersListMap := map[string]interface{}{
		resourceID: map[string]interface{}{},
	}

	consulInst := new(consulMocks.Client)
	kv := new(consulMocks.KV)
	consulInst.On("KV").Return(kv)
	consulInst.On("WaitLock", mock.Anything).Return(nil)
	kv.On("ListMap", consul.KvClustersPath, consul.KvClustersPath).Return(clustersListMap, nil)

	deps := setupTestDependencies()
	deps.consul = consulInst

	return deps
}

func setupTestApiSAPSystemTag(resourceID string) Dependencies {
	systemList := sapsystem.SAPSystemsList{
		&sapsystem.SAPSystem{
			SID: resourceID,
		},
	}

	mockSAPSystemsService := new(services.MockSAPSystemsService)
	mockSAPSystemsService.On("GetSAPSystemsById", resourceID).Return(systemList, nil)
	mockSAPSystemsService.On("GetSAPSystemsById", mock.Anything).Return(sapsystem.SAPSystemsList{}, nil)

	deps := setupTestDependencies()
	deps.sapSystemsService = mockSAPSystemsService

	return deps
}

func setupTestApiDatabaseTag(resourceID string) Dependencies {
	return setupTestApiSAPSystemTag(resourceID)
}

func TestApiResourceTag(t *testing.T) {
	cases := []struct {
		setupTest    func(resourceId string) Dependencies
		resourceType string
	}{
		{setupTestApiHostTag, models.TagHostResourceType},
		{setupTestApiClusterTag, models.TagClusterResourceType},
		{setupTestApiSAPSystemTag, models.TagSAPSystemResourceType},
		{setupTestApiDatabaseTag, models.TagDatabaseResourceType},
	}

	const tag = "tag"
	const errorTag = "guru_meditation"
	const resourceID = "resource_id"
	const notFoundResourceID = "not_found"

	for _, tc := range cases {
		deps := tc.setupTest(resourceID)

		mockTagsService := new(services.MockTagsService)
		mockTagsService.On("Create", tag, tc.resourceType, resourceID).Return(nil)
		mockTagsService.On("Delete", tag, tc.resourceType, resourceID).Return(nil)
		mockTagsService.On("Create", errorTag, tc.resourceType, resourceID).Return(fmt.Errorf("guru meditation"))
		mockTagsService.On("Delete", errorTag, tc.resourceType, resourceID).Return(fmt.Errorf("guru meditation"))
		deps.tagsService = mockTagsService

		app, err := NewAppWithDeps("", 80, deps)
		if err != nil {
			t.Fatal(err)
		}

		t.Run(fmt.Sprintf("Create %s tag", tc.resourceType), func(t *testing.T) {
			resp := httptest.NewRecorder()

			body, _ := json.Marshal(&JSONTag{tag})
			url := fmt.Sprintf("/api/%s/%s/tags", tc.resourceType, resourceID)
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			app.webEngine.ServeHTTP(resp, req)

			expectedBody, _ := json.Marshal(gin.H{
				"tag": tag,
			})
			assert.Equal(t, expectedBody, resp.Body.Bytes())
			assert.Equal(t, 201, resp.Code)
		})

		t.Run(fmt.Sprintf("Create %s tag 404", tc.resourceType), func(t *testing.T) {
			resp := httptest.NewRecorder()

			body, _ := json.Marshal(&JSONTag{tag})
			url := fmt.Sprintf("/api/%s/%s/tags", tc.resourceType, notFoundResourceID)
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			app.webEngine.ServeHTTP(resp, req)

			assert.Equal(t, 404, resp.Code)
		})

		t.Run(fmt.Sprintf("Create %s tag 400", tc.resourceType), func(t *testing.T) {
			resp := httptest.NewRecorder()

			invalidJSON := []byte("ABC€")
			url := fmt.Sprintf("/api/%s/%s/tags", tc.resourceType, resourceID)
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(invalidJSON))
			if err != nil {
				t.Fatal(err)
			}

			app.webEngine.ServeHTTP(resp, req)

			assert.Equal(t, 400, resp.Code)
		})

		t.Run(fmt.Sprintf("Create %s tag 500", tc.resourceType), func(t *testing.T) {
			resp := httptest.NewRecorder()

			body, _ := json.Marshal(&JSONTag{errorTag})
			url := fmt.Sprintf("/api/%s/%s/tags", tc.resourceType, resourceID)
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			app.webEngine.ServeHTTP(resp, req)

			assert.Equal(t, 500, resp.Code)
		})

		t.Run(fmt.Sprintf("Delete %s tag", tc.resourceType), func(t *testing.T) {
			resp := httptest.NewRecorder()

			url := fmt.Sprintf("/api/%s/%s/tags/%s", tc.resourceType, resourceID, tag)
			req, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				t.Fatal(err)
			}

			app.webEngine.ServeHTTP(resp, req)

			assert.Equal(t, 204, resp.Code)
		})

		t.Run(fmt.Sprintf("Delete %s tag 404", tc.resourceType), func(t *testing.T) {
			resp := httptest.NewRecorder()

			url := fmt.Sprintf("/api/%s/%s/tags/%s", tc.resourceType, notFoundResourceID, tag)
			req, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				t.Fatal(err)
			}

			app.webEngine.ServeHTTP(resp, req)

			assert.Equal(t, 404, resp.Code)
		})

		t.Run(fmt.Sprintf("Delete %s tag 500", tc.resourceType), func(t *testing.T) {
			resp := httptest.NewRecorder()

			url := fmt.Sprintf("/api/%s/%s/tags/%s", tc.resourceType, resourceID, errorTag)
			req, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				t.Fatal(err)
			}

			app.webEngine.ServeHTTP(resp, req)

			assert.Equal(t, 500, resp.Code)
		})
	}
}

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
	app, err := NewAppWithDeps("", 80, deps)
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
	app, err := NewAppWithDeps("", 80, deps)
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
	app, err := NewAppWithDeps("", 80, deps)
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
	app, err := NewAppWithDeps("", 80, deps)
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
