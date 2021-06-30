package web

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/consul/mocks"
)

func setupEnvironmentsTest() (*mocks.Client, *mocks.Catalog) {
	nodes1 := []*consulApi.Node{
		{
			Node:       "node1",
			Datacenter: "dc",
			Address:    "192.168.1.1",
			Meta: map[string]string{
				"trento-sap-environments": "land1",
			},
		},
		{
			Node:       "node2",
			Datacenter: "dc",
			Address:    "192.168.1.2",
			Meta: map[string]string{
				"trento-sap-environments": "land1",
			},
		},
	}

	nodes2 := []*consulApi.Node{
		{
			Node:       "node3",
			Datacenter: "dc1",
			Address:    "192.168.1.2",
			Meta: map[string]string{
				"trento-sap-environments": "land2",
			},
		},
		{
			Node:       "node4",
			Datacenter: "dc",
			Address:    "192.168.1.3",
			Meta: map[string]string{
				"trento-sap-environments": "land2",
			},
		},
	}

	node1HealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthPassing,
		},
	}

	node2HealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthPassing,
		},
	}

	node3HealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthPassing,
		},
	}

	node4HealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthCritical,
		},
	}

	filters := map[string]interface{}{
		"env1": map[string]interface{}{
			"name": "env1",
			"landscapes": map[string]interface{}{
				"land1": map[string]interface{}{
					"name": "land1",
					"sapsystems": map[string]interface{}{
						"sys1": map[string]interface{}{
							"name": "sys1",
						},
					},
				},
				"land2": map[string]interface{}{
					"name": "land2",
					"sapsystems": map[string]interface{}{
						"sys2": map[string]interface{}{
							"name": "sys2",
						},
					},
				},
			},
		},
		"env2": map[string]interface{}{
			"name": "env2",
			"landscapes": map[string]interface{}{
				"land3": map[string]interface{}{
					"name": "land3",
					"sapsystems": map[string]interface{}{
						"sys3": map[string]interface{}{
							"name": "sys3",
						},
					},
				},
			},
		},
	}

	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	health := new(mocks.Health)
	kv := new(mocks.KV)

	consulInst.On("Catalog").Return(catalog)
	consulInst.On("Health").Return(health)
	consulInst.On("KV").Return(kv)

	kv.On("ListMap", consul.KvEnvironmentsPath, consul.KvEnvironmentsPath).Return(filters, nil)

	filterSys1 := &consulApi.QueryOptions{
		Filter: "(Meta[\"trento-sap-environment\"] == \"env1\") and (Meta[\"trento-sap-landscape\"] == \"land1\") and (Meta[\"trento-sap-system\"] == \"sys1\")"}
	catalog.On("Nodes", filterSys1).Return(nodes1, nil, nil)

	filterSys2 := &consulApi.QueryOptions{
		Filter: "(Meta[\"trento-sap-environment\"] == \"env1\") and (Meta[\"trento-sap-landscape\"] == \"land2\") and (Meta[\"trento-sap-system\"] == \"sys2\")"}
	catalog.On("Nodes", filterSys2).Return(nodes1, nil, nil)

	filterSys3 := &consulApi.QueryOptions{
		Filter: "(Meta[\"trento-sap-environment\"] == \"env2\") and (Meta[\"trento-sap-landscape\"] == \"land3\") and (Meta[\"trento-sap-system\"] == \"sys3\")"}
	catalog.On("Nodes", filterSys3).Return(nodes2, nil, nil)

	health.On("Node", "node1", (*consulApi.QueryOptions)(nil)).Return(node1HealthChecks, nil, nil)
	health.On("Node", "node2", (*consulApi.QueryOptions)(nil)).Return(node2HealthChecks, nil, nil)
	health.On("Node", "node3", (*consulApi.QueryOptions)(nil)).Return(node3HealthChecks, nil, nil)
	health.On("Node", "node4", (*consulApi.QueryOptions)(nil)).Return(node4HealthChecks, nil, nil)

	return consulInst, catalog
}

func TestEnvironmentsListHandler(t *testing.T) {
	consulInst, catalog := setupEnvironmentsTest()

	deps := DefaultDependencies()
	deps.consul = consulInst

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/environments", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	consulInst.AssertExpectations(t)
	catalog.AssertExpectations(t)

	responseBody := minifyHtml(resp.Body.String())

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, responseBody, "Environments")
	assert.Regexp(t, regexp.MustCompile("<tr.*onclick=\"window.location='/environments/env1'\".*<td>env1</td><td>2</td><td>2</td><td>.*passing.*</td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<tr.*onclick=\"window.location='/environments/env1'\".*<td>env2</td><td>1</td><td>1</td><td>.*critical.*</td>"), responseBody)
}

func TestLandscapesListHandler(t *testing.T) {
	consulInst, catalog := setupEnvironmentsTest()

	deps := DefaultDependencies()
	deps.consul = consulInst

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/landscapes", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	consulInst.AssertExpectations(t)
	catalog.AssertExpectations(t)

	responseBody := minifyHtml(resp.Body.String())

	assert.Equal(t, 200, resp.Code)
	assert.Regexp(t, regexp.MustCompile("<tr.*onclick=\"window.location='/landscapes/land1\\?environment=env1'\"><td>land1</td><td>.*env1.*</td><td>1</td><td>2</td><td>.*passing.*</td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<tr.*onclick=\"window.location='/landscapes/land2\\?environment=env1'\".*<td>land2</td><td>.*env1.*<td>1</td><td>2</td><td>.*passing.*</td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<tr.*onclick=\"window.location='/landscapes/land3\\?environment=env2'\".*<td>land3</td><td>.*env2.*<td>1</td><td>2</td><td>.*critical.*</td>"), responseBody)
}

func TestSAPSystemsListHandler(t *testing.T) {
	consulInst, catalog := setupEnvironmentsTest()

	deps := DefaultDependencies()
	deps.consul = consulInst

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/sapsystems", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	consulInst.AssertExpectations(t)
	catalog.AssertExpectations(t)

	responseBody := minifyHtml(resp.Body.String())

	assert.Equal(t, 200, resp.Code)
	assert.Regexp(t, regexp.MustCompile("<tr.*onclick=\"window.location='/sapsystems/sys1\\?environment=env1&landscape=land1'\".*<td>sys1</td><td>.*passing.*</td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<tr.*onclick=\"window.location='/sapsystems/sys2\\?environment=env1&landscape=land2'\".*<td>sys2</td><td>.*passing.*</td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<tr.*onclick=\"window.location='/sapsystems/sys3\\?environment=env2&landscape=land3'\".*<td>sys3</td><td>.*critical.*</td>"), responseBody)
}

func TestLandscapeHandler(t *testing.T) {
	consulInst, _ := setupEnvironmentsTest()

	deps := DefaultDependencies()
	deps.consul = consulInst

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/landscapes/land1?environment=env1", nil)
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "Landscape details")
	assert.Contains(t, resp.Body.String(), "land1")
}

func TestLandscapeHandler404Error(t *testing.T) {
	consulInst, _ := setupEnvironmentsTest()

	deps := DefaultDependencies()
	deps.consul = consulInst

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/landscapes/foobar", nil)
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.Code)
	assert.Contains(t, resp.Body.String(), "Not Found")
}

func TestEnvironmentHandler(t *testing.T) {
	consulInst, _ := setupEnvironmentsTest()

	deps := DefaultDependencies()
	deps.consul = consulInst

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/environments/env1", nil)
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "Environment details")
	assert.Contains(t, resp.Body.String(), "env1")
}

func TestEnvironmentHandler404Error(t *testing.T) {
	consulInst, _ := setupEnvironmentsTest()

	deps := DefaultDependencies()
	deps.consul = consulInst

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/environments/foobar", nil)
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.Code)
	assert.Contains(t, resp.Body.String(), "Not Found")
}

func TestSAPSystemHandler(t *testing.T) {
	consulInst, _ := setupEnvironmentsTest()

	deps := DefaultDependencies()
	deps.consul = consulInst

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/sapsystems/sys1?environment=env1&landscape=land1", nil)
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "SAP System details")
	assert.Contains(t, resp.Body.String(), "sys1")
}

func TestSAPSystemHandler404Error(t *testing.T) {
	consulInst, _ := setupEnvironmentsTest()

	deps := DefaultDependencies()
	deps.consul = consulInst

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/sapsystems/foobar", nil)
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

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
