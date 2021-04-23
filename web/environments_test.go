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
	"github.com/trento-project/trento/internal/consul/mocks"
)

func setupTest() (*mocks.Client, *mocks.Catalog) {
	nodes := []*consulApi.Node{
		{
			Node:       "foo",
			Datacenter: "dc1",
			Address:    "192.168.1.1",
			Meta: map[string]string{
				"trento-sap-environments": "land1",
			},
		},
		{
			Node:       "bar",
			Datacenter: "dc",
			Address:    "192.168.1.2",
			Meta: map[string]string{
				"trento-sap-environments": "land2",
			},
		},
	}

	filters := []string{
		"trento/environments/",
		"trento/environments/env1/",
		"trento/environments/env1/landscapes/",
		"trento/environments/env1/landscapes/land1/",
		"trento/environments/env1/landscapes/land1/sapsystems/",
		"trento/environments/env1/landscapes/land1/sapsystems/sys1/",
		"trento/environments/env1/landscapes/land2/",
		"trento/environments/env1/landscapes/land2/sapsystems/",
		"trento/environments/env1/landscapes/land2/sapsystems/sys2/",
		"trento/environments/env2/",
		"trento/environments/env2/landscapes/",
		"trento/environments/env2/landscapes/land3/",
		"trento/environments/env2/landscapes/land3/sapsystems/",
		"trento/environments/env2/landscapes/land3/sapsystems/sys3/",
	}

	consul := new(mocks.Client)
	catalog := new(mocks.Catalog)
	kv := new(mocks.KV)

	consul.On("Catalog").Return(catalog)
	consul.On("KV").Return(kv)

	kv.On("Keys", "trento/environments", "", (*consulApi.QueryOptions)(nil)).Return(filters, nil, nil)

	filterSys1 := &consulApi.QueryOptions{
		Filter: "(Meta[\"trento-sap-environment\"] == \"env1\") and (Meta[\"trento-sap-landscape\"] == \"land1\") and (Meta[\"trento-sap-system\"] == \"sys1\")"}
	catalog.On("Nodes", (filterSys1)).Return(nodes, nil, nil)

	filterSys2 := &consulApi.QueryOptions{
		Filter: "(Meta[\"trento-sap-environment\"] == \"env1\") and (Meta[\"trento-sap-landscape\"] == \"land2\") and (Meta[\"trento-sap-system\"] == \"sys2\")"}
	catalog.On("Nodes", (filterSys2)).Return(nodes, nil, nil)

	filterSys3 := &consulApi.QueryOptions{
		Filter: "(Meta[\"trento-sap-environment\"] == \"env2\") and (Meta[\"trento-sap-landscape\"] == \"land3\") and (Meta[\"trento-sap-system\"] == \"sys3\")"}
	catalog.On("Nodes", (filterSys3)).Return(nodes, nil, nil)

	return consul, catalog
}

func TestEnvironmentsListHandler(t *testing.T) {
	consul, catalog := setupTest()

	deps := DefaultDependencies()
	deps.consul = consul

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

	consul.AssertExpectations(t)
	catalog.AssertExpectations(t)

	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepEndTags:         true,
	})
	minified, err := m.String("text/html", resp.Body.String())
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, minified, "Environments")
	assert.Regexp(t, regexp.MustCompile("<tr.*onclick=\"window.location='/environments/env1'\".*<td>env1</td><td>2</td><td>2</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<tr.*onclick=\"window.location='/environments/env2'\".*<td>env2</td><td>1</td><td>1</td>"), minified)
}

func TestLandscapesListHandler(t *testing.T) {
	consul, catalog := setupTest()

	deps := DefaultDependencies()
	deps.consul = consul

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/environments/env1", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	consul.AssertExpectations(t)
	catalog.AssertExpectations(t)

	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepEndTags:         true,
	})
	minified, err := m.String("text/html", resp.Body.String())
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 200, resp.Code)
	assert.Regexp(t, regexp.MustCompile("<tr.*onclick=\"window.location='/environments/env1/landscapes/land1'\".*<td>land1</td><td>1</td><td>2</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<tr.*onclick=\"window.location='/environments/env1/landscapes/land2'\".*<td>land2</td><td>1</td><td>2</td>"), minified)
}

func TestSAPSystemsListHandler(t *testing.T) {
	consul, catalog := setupTest()

	deps := DefaultDependencies()
	deps.consul = consul

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/environments/env1/landscapes/land1", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	consul.AssertExpectations(t)
	catalog.AssertExpectations(t)

	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepEndTags:         true,
	})
	minified, err := m.String("text/html", resp.Body.String())
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 200, resp.Code)
	assert.Regexp(t, regexp.MustCompile("<tr.*onclick=\"window.location='/environments/env1/landscapes/land1/sapsystems/sys1'\".*<td>sys1</td>"), minified)
}
