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

	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
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

		config := setupTestConfig()
		app, err := NewAppWithDeps(config, deps)
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
