package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/mock"

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

	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/tags", nil)
	app.webEngine.ServeHTTP(resp, req)

	expectedBody, _ := json.Marshal(tags)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, expectedBody, resp.Body.Bytes())

	resp = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/tags?resource_type=hosts", nil)
	app.webEngine.ServeHTTP(resp, req)

	expectedBody, _ = json.Marshal([]string{
		"tag4",
		"tag5",
		"tag6",
	})
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, expectedBody, resp.Body.Bytes())

	resp = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/tags?resource_type=sapsystems", nil)
	app.webEngine.ServeHTTP(resp, req)

	expectedBody, _ = json.Marshal([]string{})
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, expectedBody, resp.Body.Bytes())
}

func setupTestApiHostTag(resourceID string) Dependencies {
	host := &models.Host{
		ID: resourceID,
	}

	hostsService := new(services.MockHostsService)
	hostsService.On("GetByID", resourceID).Return(host, nil)
	hostsService.On("GetByID", mock.Anything).Return(nil, nil)

	deps := setupTestDependencies()
	deps.hostsService = hostsService

	return deps
}

func setupTestApiClusterTag(resourceID string) Dependencies {

	cluster := &models.Cluster{
		ID: resourceID,
	}

	clustersService := new(services.MockClustersService)
	clustersService.On("GetByID", resourceID).Return(cluster, nil)
	clustersService.On("GetByID", mock.Anything).Return(nil, nil)

	deps := setupTestDependencies()
	deps.clustersService = clustersService

	return deps
}

func setupTestApiSAPSystemTag(resourceID string) Dependencies {
	sapSystem := &models.SAPSystem{
		ID: resourceID,
	}

	mockSAPSystemsService := new(services.MockSAPSystemsService)
	mockSAPSystemsService.On("GetByID", resourceID).Return(sapSystem, nil)
	mockSAPSystemsService.On("GetByID", mock.Anything).Return(nil, nil)

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

		config := setupTestConfig()
		app, err := NewAppWithDeps(config, deps)
		if err != nil {
			t.Fatal(err)
		}

		t.Run(fmt.Sprintf("Create %s tag", tc.resourceType), func(t *testing.T) {
			resp := httptest.NewRecorder()

			body, _ := json.Marshal(&JSONTag{tag})
			url := fmt.Sprintf("/api/%s/%s/tags", tc.resourceType, resourceID)
			req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))

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
			req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))

			app.webEngine.ServeHTTP(resp, req)

			assert.Equal(t, 404, resp.Code)
		})

		t.Run(fmt.Sprintf("Create %s tag 400", tc.resourceType), func(t *testing.T) {
			resp := httptest.NewRecorder()

			invalidJSON := []byte("ABC€")
			url := fmt.Sprintf("/api/%s/%s/tags", tc.resourceType, resourceID)
			req := httptest.NewRequest("POST", url, bytes.NewBuffer(invalidJSON))

			app.webEngine.ServeHTTP(resp, req)

			assert.Equal(t, 400, resp.Code)
		})

		t.Run(fmt.Sprintf("Create %s tag 500", tc.resourceType), func(t *testing.T) {
			resp := httptest.NewRecorder()

			body, _ := json.Marshal(&JSONTag{errorTag})
			url := fmt.Sprintf("/api/%s/%s/tags", tc.resourceType, resourceID)
			req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))

			app.webEngine.ServeHTTP(resp, req)

			assert.Equal(t, 500, resp.Code)
		})

		t.Run(fmt.Sprintf("Delete %s tag", tc.resourceType), func(t *testing.T) {
			resp := httptest.NewRecorder()

			url := fmt.Sprintf("/api/%s/%s/tags/%s", tc.resourceType, resourceID, tag)
			req := httptest.NewRequest("DELETE", url, nil)

			app.webEngine.ServeHTTP(resp, req)

			assert.Equal(t, 204, resp.Code)
		})

		t.Run(fmt.Sprintf("Delete %s tag 404", tc.resourceType), func(t *testing.T) {
			resp := httptest.NewRecorder()

			url := fmt.Sprintf("/api/%s/%s/tags/%s", tc.resourceType, notFoundResourceID, tag)
			req := httptest.NewRequest("DELETE", url, nil)

			app.webEngine.ServeHTTP(resp, req)

			assert.Equal(t, 404, resp.Code)
		})

		t.Run(fmt.Sprintf("Delete %s tag 500", tc.resourceType), func(t *testing.T) {
			resp := httptest.NewRecorder()

			url := fmt.Sprintf("/api/%s/%s/tags/%s", tc.resourceType, resourceID, errorTag)
			req := httptest.NewRequest("DELETE", url, nil)

			app.webEngine.ServeHTTP(resp, req)

			assert.Equal(t, 500, resp.Code)
		})
	}
}
