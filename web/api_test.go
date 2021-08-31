package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/consul/api"

	"github.com/gin-gonic/gin"

	"github.com/trento-project/trento/internal/consul"

	"github.com/stretchr/testify/mock"

	"github.com/trento-project/trento/internal/consul/mocks"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

func setupTestApiHostTag() Dependencies {
	node := &consulApi.Node{
		Node: "suse",
	}

	consulInst := new(mocks.Client)

	kv := new(mocks.KV)
	kv.On("PutMap", "trento/v0/tags/hosts/suse/cool_rabbit/", map[string]interface{}(nil)).Return(nil)
	kv.On("DeleteTree", "trento/v0/tags/hosts/suse/cool_rabbit/", (*api.WriteOptions)(nil)).Return(nil, nil)
	consulInst.On("KV").Return(kv)

	catalog := new(mocks.Catalog)
	catalogNode := &consulApi.CatalogNode{Node: node}
	catalog.On("Node", "suse", (*consulApi.QueryOptions)(nil)).Return(catalogNode, nil, nil)
	catalog.On("Node", mock.Anything, (*consulApi.QueryOptions)(nil)).Return(nil, nil, nil)
	consulInst.On("Catalog").Return(catalog)

	deps := DefaultDependencies()
	deps.consul = consulInst

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
	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	catalog.On("Node", "suse", mock.Anything).Return(nil, nil, fmt.Errorf("kaboom"))
	consulInst.On("Catalog").Return(catalog)

	deps := DefaultDependencies()
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
	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	catalog.On("Node", "suse", mock.Anything).Return(nil, nil, fmt.Errorf("kaboom"))
	consulInst.On("Catalog").Return(catalog)

	deps := DefaultDependencies()
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

	consulInst := new(mocks.Client)
	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)
	consulInst.On("WaitLock", mock.Anything).Return(nil)
	kv.On("ListMap", consul.KvClustersPath, consul.KvClustersPath).Return(clustersListMap, nil)
	kv.On("PutMap", "trento/v0/tags/clusters/47d1190ffb4f781974c8356d7f863b03/cool_rabbit/", map[string]interface{}(nil)).Return(nil)
	kv.On("DeleteTree", "trento/v0/tags/clusters/47d1190ffb4f781974c8356d7f863b03/cool_rabbit/", (*api.WriteOptions)(nil)).Return(nil, nil)

	deps := DefaultDependencies()
	deps.consul = consulInst

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
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)
	consulInst.On("WaitLock", mock.Anything).Return(fmt.Errorf("kaboom"))

	deps := DefaultDependencies()
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
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)
	consulInst.On("WaitLock", mock.Anything).Return(fmt.Errorf("kaboom"))

	deps := DefaultDependencies()
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

func setupTestApiSAPSystemTag() Dependencies {
	sapSystemMap := map[string]interface{}{
		"HA1": map[string]interface{}{
			"sid": "HA1",
		},
	}

	nodes := []*consulApi.Node{
		{
			Node: "test_host",
		},
	}

	consulInst := new(mocks.Client)
	consulInst.On("WaitLock", mock.Anything).Return(nil)

	catalog := new(mocks.Catalog)
	catalog.On("Nodes", mock.Anything).Return(nodes, nil, nil)
	consulInst.On("Catalog").Return(catalog)

	kv := new(mocks.KV)
	path := fmt.Sprintf(consul.KvHostsSAPSystemPath, "test_host")
	kv.On("ListMap", path, path).Return(sapSystemMap, nil)
	kv.On("PutMap", "trento/v0/tags/sapsystems/HA1/cool_rabbit/", map[string]interface{}(nil)).Return(nil)
	kv.On("DeleteTree", "trento/v0/tags/sapsystems/HA1/cool_rabbit/", (*api.WriteOptions)(nil)).Return(nil, nil)
	consulInst.On("KV").Return(kv)

	deps := DefaultDependencies()
	deps.consul = consulInst

	return deps
}

func TestApiSAPSystemCreateTagHandler(t *testing.T) {
	deps := setupTestApiSAPSystemTag()
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
	deps := setupTestApiSAPSystemTag()
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
	deps := setupTestApiSAPSystemTag()
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
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)
	consulInst.On("WaitLock", mock.Anything).Return(fmt.Errorf("kaboom"))

	deps := DefaultDependencies()
	deps.consul = consulInst

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
	deps := setupTestApiSAPSystemTag()
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
	deps := setupTestApiSAPSystemTag()
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "/api/sapsystems/non-existing-id/tags/cool_rabbit", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	assert.Equal(t, 404, resp.Code)
}

func TestApiSAPSystemDeleteTagHandler500(t *testing.T) {
	consulInst := new(mocks.Client)
	kv := new(mocks.KV)
	consulInst.On("KV").Return(kv)
	consulInst.On("WaitLock", mock.Anything).Return(fmt.Errorf("kaboom"))

	deps := DefaultDependencies()
	deps.consul = consulInst

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
