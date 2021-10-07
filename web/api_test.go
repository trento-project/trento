package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/sapsystem"

	"github.com/stretchr/testify/mock"
	consulMocks "github.com/trento-project/trento/internal/consul/mocks"
	servicesMocks "github.com/trento-project/trento/web/services/mocks"

	"github.com/trento-project/trento/web/models"

	consulApi "github.com/hashicorp/consul/api"
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

	mockTagsService := new(servicesMocks.TagsService)
	mockTagsService.On("GetAll").Return(tags, nil)
	mockTagsService.On("GetAll", "sapsystems").Return(tagsSAPSystems, nil)
	mockTagsService.On("GetAll", "hosts").Return(tagsHosts, nil)
	deps := testDependencies()
	deps.tagsService = mockTagsService

	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/tags", nil)
	app.ServeHTTP(resp, req)

	expectedBody, _ := json.Marshal(tags)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, expectedBody, resp.Body.Bytes())

	resp = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/tags?resource_type=hosts", nil)
	app.ServeHTTP(resp, req)

	expectedBody, _ = json.Marshal([]string{
		"tag4",
		"tag5",
		"tag6",
	})
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, expectedBody, resp.Body.Bytes())

	resp = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/tags?resource_type=sapsystems", nil)
	app.ServeHTTP(resp, req)

	expectedBody, _ = json.Marshal([]string{})
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, expectedBody, resp.Body.Bytes())
}

func setupTestApiHostTag() Dependencies {
	node := &consulApi.Node{
		Node: "suse",
	}

	consulInst := new(consulMocks.Client)
	catalog := new(consulMocks.Catalog)
	catalogNode := &consulApi.CatalogNode{Node: node}
	catalog.On("Node", "suse", (*consulApi.QueryOptions)(nil)).Return(catalogNode, nil, nil)
	catalog.On("Node", mock.Anything, (*consulApi.QueryOptions)(nil)).Return(nil, nil, nil)
	consulInst.On("Catalog").Return(catalog)

	mockTagsService := new(servicesMocks.TagsService)
	mockTagsService.On("Create", "cool_rabbit", models.TagHostResourceType, "suse").Return(nil)
	mockTagsService.On("Delete", "cool_rabbit", models.TagHostResourceType, "suse").Return(nil)

	deps := testDependencies()
	deps.consul = consulInst
	deps.tagsService = mockTagsService

	return deps
}

func TestApiHostCreateTagHandler(t *testing.T) {
	deps := setupTestApiHostTag()
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	body, _ := json.Marshal(&JSONTag{"cool_rabbit"})
	req, err := http.NewRequest("POST", "/api/hosts/suse/tags", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	expectedBody, _ := json.Marshal(gin.H{
		"tag": "cool_rabbit",
	})
	assert.Equal(t, expectedBody, resp.Body.Bytes())
	assert.Equal(t, 201, resp.Code)
}

func TestApiHostCreateTagHandler404(t *testing.T) {
	deps := setupTestApiHostTag()
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	body, _ := json.Marshal(&JSONTag{"cool_rabbit"})
	req, err := http.NewRequest("POST", "/api/hosts/non-existing-host/tags", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 404, resp.Code)
}

func TestApiHostCreateTagHandler400(t *testing.T) {
	deps := setupTestApiHostTag()
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	invalidJSON := []byte("ABC€")
	req, err := http.NewRequest("POST", "/api/hosts/suse/tags", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code)
}

func TestApiHostCreateTagHandler500(t *testing.T) {
	consulInst := new(consulMocks.Client)
	catalog := new(consulMocks.Catalog)
	catalog.On("Node", "suse", mock.Anything).Return(nil, nil, fmt.Errorf("kaboom"))
	consulInst.On("Catalog").Return(catalog)

	deps := testDependencies()
	deps.consul = consulInst

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	body, _ := json.Marshal(&JSONTag{"cool_rabbit"})
	req, err := http.NewRequest("POST", "/api/hosts/suse/tags", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)
}

func TestApiHostDeleteTagHandler(t *testing.T) {
	deps := setupTestApiHostTag()
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "/api/hosts/suse/tags/cool_rabbit", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 204, resp.Code)
}

func TestApiHostDeleteTagHandler404(t *testing.T) {
	deps := setupTestApiHostTag()
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "/api/hosts/non-existing-host/tags/cool_rabbit", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 404, resp.Code)
}

func TestApiHostDeleteTagHandler500(t *testing.T) {
	consulInst := new(consulMocks.Client)
	catalog := new(consulMocks.Catalog)
	catalog.On("Node", "suse", mock.Anything).Return(nil, nil, fmt.Errorf("kaboom"))
	consulInst.On("Catalog").Return(catalog)

	deps := testDependencies()
	deps.consul = consulInst

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "/api/hosts/suse/tags/cool_rabbit", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)
}

func setupTestApiClusterTag() Dependencies {
	clustersListMap := map[string]interface{}{
		"47d1190ffb4f781974c8356d7f863b03": map[string]interface{}{},
	}

	consulInst := new(consulMocks.Client)
	kv := new(consulMocks.KV)
	consulInst.On("KV").Return(kv)
	consulInst.On("WaitLock", mock.Anything).Return(nil)
	kv.On("ListMap", consul.KvClustersPath, consul.KvClustersPath).Return(clustersListMap, nil)

	mockTagsService := new(servicesMocks.TagsService)
	mockTagsService.On("Create", "cool_rabbit", models.TagClusterResourceType, "47d1190ffb4f781974c8356d7f863b03").Return(nil)
	mockTagsService.On("Delete", "cool_rabbit", models.TagClusterResourceType, "47d1190ffb4f781974c8356d7f863b03").Return(nil)

	deps := testDependencies()
	deps.consul = consulInst
	deps.tagsService = mockTagsService

	return deps
}

func TestApiClusterCreateTagHandler(t *testing.T) {
	deps := setupTestApiClusterTag()
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	body, _ := json.Marshal(&JSONTag{"cool_rabbit"})
	req, err := http.NewRequest("POST", "/api/clusters/47d1190ffb4f781974c8356d7f863b03/tags", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	expectedBody, _ := json.Marshal(gin.H{
		"tag": "cool_rabbit",
	})
	assert.Equal(t, expectedBody, resp.Body.Bytes())
	assert.Equal(t, 201, resp.Code)
}

func TestApiClusterCreateTagHandler404(t *testing.T) {
	deps := setupTestApiClusterTag()
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	body, _ := json.Marshal(&JSONTag{"cool_rabbit"})
	req, err := http.NewRequest("POST", "/api/clusters/non-existing-id/tags", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 404, resp.Code)
}

func TestApiClusterCreateTagHandler400(t *testing.T) {
	deps := setupTestApiClusterTag()
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	invalidJSON := []byte("ABC€")
	req, err := http.NewRequest("POST", "/api/clusters/47d1190ffb4f781974c8356d7f863b03/tags", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code)
}

func TestApiClusterCreateTagHandler500(t *testing.T) {
	consulInst := new(consulMocks.Client)
	kv := new(consulMocks.KV)
	consulInst.On("KV").Return(kv)
	consulInst.On("WaitLock", mock.Anything).Return(fmt.Errorf("kaboom"))

	deps := testDependencies()
	deps.consul = consulInst

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	body, _ := json.Marshal(&JSONTag{"cool_rabbit"})
	req, err := http.NewRequest("POST", "/api/clusters/47d1190ffb4f781974c8356d7f863b03/tags", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)
}

func TestApiClusterDeleteTagHandler(t *testing.T) {
	deps := setupTestApiClusterTag()
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "/api/clusters/47d1190ffb4f781974c8356d7f863b03/tags/cool_rabbit", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 204, resp.Code)
}

func TestApiClusterDeleteTagHandler404(t *testing.T) {
	deps := setupTestApiClusterTag()
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "/api/clusters/non-existing-id/tags/cool_rabbit", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 404, resp.Code)
}

func TestApiClusterDeleteTagHandler500(t *testing.T) {
	consulInst := new(consulMocks.Client)
	kv := new(consulMocks.KV)
	consulInst.On("KV").Return(kv)
	consulInst.On("WaitLock", mock.Anything).Return(fmt.Errorf("kaboom"))

	deps := testDependencies()
	deps.consul = consulInst

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "/api/clusters/47d1190ffb4f781974c8356d7f863b03/tags/cool_rabbit", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)
}

func setupTestApiSAPSystemTag(systemType string) Dependencies {
	systemList := sapsystem.SAPSystemsList{
		&sapsystem.SAPSystem{
			SID: "HA1",
		},
	}

	mockTagsService := new(servicesMocks.TagsService)
	mockTagsService.On("Create", "cool_rabbit", systemType, "HA1").Return(nil)
	mockTagsService.On("Delete", "cool_rabbit", systemType, "HA1").Return(nil)

	mockSAPSystemsService := new(servicesMocks.SAPSystemsService)
	mockSAPSystemsService.On("GetSAPSystemsBySid", "HA1").Return(systemList, nil)
	mockSAPSystemsService.On(
		"GetSAPSystemsBySid", "non-existing-sid").Return(sapsystem.SAPSystemsList{}, nil)

	deps := testDependencies()
	deps.sapSystemsService = mockSAPSystemsService
	deps.tagsService = mockTagsService

	return deps
}

func TestApiSAPSystemCreateTagHandler(t *testing.T) {
	deps := setupTestApiSAPSystemTag(models.TagSAPSystemResourceType)
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	body, _ := json.Marshal(&JSONTag{"cool_rabbit"})
	req, err := http.NewRequest("POST", "/api/sapsystems/HA1/tags", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	expectedBody, _ := json.Marshal(gin.H{
		"tag": "cool_rabbit",
	})
	assert.Equal(t, expectedBody, resp.Body.Bytes())
	assert.Equal(t, 201, resp.Code)
}

func TestApiSAPSystemCreateTagHandler404(t *testing.T) {
	deps := setupTestApiSAPSystemTag(models.TagSAPSystemResourceType)
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	body, _ := json.Marshal(&JSONTag{"cool_rabbit"})
	req, err := http.NewRequest("POST", "/api/sapsystems/non-existing-sid/tags", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 404, resp.Code)
}

func TestApiSAPSystemCreateTagHandler400(t *testing.T) {
	deps := setupTestApiSAPSystemTag(models.TagSAPSystemResourceType)
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	invalidJSON := []byte("ABC€")
	req, err := http.NewRequest("POST", "/api/sapsystems/HA1/tags", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code)
}

func TestApiSAPSystemCreateTagHandler500(t *testing.T) {
	mockSAPSystemsService := new(servicesMocks.SAPSystemsService)
	mockSAPSystemsService.On(
		"GetSAPSystemsBySid", "HA1").Return(sapsystem.SAPSystemsList{}, fmt.Errorf("kaboom"))

	deps := testDependencies()
	deps.sapSystemsService = mockSAPSystemsService

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	body, _ := json.Marshal(&JSONTag{"cool_rabbit"})
	req, err := http.NewRequest("POST", "/api/sapsystems/HA1/tags", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)
}

func TestApiSAPSystemDeleteTagHandler(t *testing.T) {
	deps := setupTestApiSAPSystemTag(models.TagSAPSystemResourceType)
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "/api/sapsystems/HA1/tags/cool_rabbit", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 204, resp.Code)
}

func TestApiSAPSystemDeleteTagHandler404(t *testing.T) {
	deps := setupTestApiSAPSystemTag(models.TagSAPSystemResourceType)
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "/api/sapsystems/non-existing-sid/tags/cool_rabbit", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 404, resp.Code)
}

func TestApiSAPSystemDeleteTagHandler500(t *testing.T) {
	mockSAPSystemsService := new(servicesMocks.SAPSystemsService)
	mockSAPSystemsService.On(
		"GetSAPSystemsBySid", "HA1").Return(sapsystem.SAPSystemsList{}, fmt.Errorf("kaboom"))

	deps := testDependencies()
	deps.sapSystemsService = mockSAPSystemsService

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "/api/sapsystems/HA1/tags/cool_rabbit", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)
}

func TestApiDatabaseCreateTagHandler(t *testing.T) {
	deps := setupTestApiSAPSystemTag(models.TagDatabaseResourceType)
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	body, _ := json.Marshal(&JSONTag{"cool_rabbit"})
	req, err := http.NewRequest("POST", "/api/databases/HA1/tags", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	expectedBody, _ := json.Marshal(gin.H{
		"tag": "cool_rabbit",
	})
	assert.Equal(t, expectedBody, resp.Body.Bytes())
	assert.Equal(t, 201, resp.Code)
}

func TestApiDatabaseCreateTagHandler404(t *testing.T) {
	deps := setupTestApiSAPSystemTag(models.TagDatabaseResourceType)
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	body, _ := json.Marshal(&JSONTag{"cool_rabbit"})
	req, err := http.NewRequest("POST", "/api/databases/non-existing-sid/tags", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 404, resp.Code)
}

func TestApiDatabaseCreateTagHandler400(t *testing.T) {
	deps := setupTestApiSAPSystemTag(models.TagDatabaseResourceType)
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	invalidJSON := []byte("ABC€")
	req, err := http.NewRequest("POST", "/api/databases/HA1/tags", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code)
}

func TestApiDatabaseCreateTagHandler500(t *testing.T) {
	mockSAPSystemsService := new(servicesMocks.SAPSystemsService)
	mockSAPSystemsService.On(
		"GetSAPSystemsBySid", "HA1").Return(sapsystem.SAPSystemsList{}, fmt.Errorf("kaboom"))

	deps := testDependencies()
	deps.sapSystemsService = mockSAPSystemsService

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	body, _ := json.Marshal(&JSONTag{"cool_rabbit"})
	req, err := http.NewRequest("POST", "/api/databases/HA1/tags", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)
}

func TestApiDatabaseDeleteTagHandler(t *testing.T) {
	deps := setupTestApiSAPSystemTag(models.TagDatabaseResourceType)
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "/api/databases/HA1/tags/cool_rabbit", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 204, resp.Code)
}

func TestApiDatabaseDeleteTagHandler404(t *testing.T) {
	deps := setupTestApiSAPSystemTag(models.TagDatabaseResourceType)
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "/api/databases/non-existing-sid/tags/cool_rabbit", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 404, resp.Code)
}

func TestApiDatabaseDeleteTagHandler500(t *testing.T) {
	mockSAPSystemsService := new(servicesMocks.SAPSystemsService)
	mockSAPSystemsService.On(
		"GetSAPSystemsBySid", "HA1").Return(sapsystem.SAPSystemsList{}, fmt.Errorf("kaboom"))

	deps := testDependencies()
	deps.sapSystemsService = mockSAPSystemsService

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "/api/databases/HA1/tags/cool_rabbit", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)
}
