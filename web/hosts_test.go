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

func TestHostsListHandler(t *testing.T) {
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

	fooHealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthPassing,
		},
	}

	barHealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthCritical,
		},
	}

	filters := map[string]interface{}{
		"env1": map[string]interface{}{
			"landscapes": map[string]interface{}{
				"land1": map[string]interface{}{
					"sapsystems": map[string]interface{}{
						"sys1": map[string]interface{}{},
					},
				},
				"land2": map[string]interface{}{
					"sapsystems": map[string]interface{}{
						"sys2": map[string]interface{}{},
					},
				},
			},
		},
		"env2": map[string]interface{}{
			"landscapes": map[string]interface{}{
				"land3": map[string]interface{}{
					"sapsystems": map[string]interface{}{
						"sys3": map[string]interface{}{},
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

	query := &consulApi.QueryOptions{Filter: ""}
	catalog.On("Nodes", (*consulApi.QueryOptions)(query)).Return(nodes, nil, nil)

	filterSys1 := &consulApi.QueryOptions{
		Filter: "(Meta[\"trento-sap-environment\"] == \"env1\") and (Meta[\"trento-sap-landscape\"] == \"land1\") and (Meta[\"trento-sap-system\"] == \"sys1\")"}
	catalog.On("Nodes", (filterSys1)).Return(nodes, nil, nil)

	filterSys2 := &consulApi.QueryOptions{
		Filter: "(Meta[\"trento-sap-environment\"] == \"env1\") and (Meta[\"trento-sap-landscape\"] == \"land2\") and (Meta[\"trento-sap-system\"] == \"sys2\")"}
	catalog.On("Nodes", (filterSys2)).Return(nodes, nil, nil)

	filterSys3 := &consulApi.QueryOptions{
		Filter: "(Meta[\"trento-sap-environment\"] == \"env2\") and (Meta[\"trento-sap-landscape\"] == \"land3\") and (Meta[\"trento-sap-system\"] == \"sys3\")"}
	catalog.On("Nodes", (filterSys3)).Return(nodes, nil, nil)

	health.On("Node", "foo", (*consulApi.QueryOptions)(nil)).Return(fooHealthChecks, nil, nil)
	health.On("Node", "bar", (*consulApi.QueryOptions)(nil)).Return(barHealthChecks, nil, nil)

	deps := DefaultDependencies()
	deps.consul = consulInst

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/hosts", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	consulInst.AssertExpectations(t)
	catalog.AssertExpectations(t)
	health.AssertExpectations(t)

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
	assert.Contains(t, minified, "Hosts")
	assert.Regexp(t, regexp.MustCompile("<select name=trento-sap-environment.*>.*env1.*env2.*</select>"), minified)
	assert.Regexp(t, regexp.MustCompile("<select name=trento-sap-landscape.*>.*land1.*land2.*land3.*</select>"), minified)
	assert.Regexp(t, regexp.MustCompile("<select name=trento-sap-system.*>.*sys1.*sys2.*sys3.*</select>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>foo</td><td>192.168.1.1</td><td>.*land1.*</td><td>.*passing.*</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>bar</td><td>192.168.1.2</td><td>.*land2.*</td><td>.*critical.*</td>"), minified)
}
