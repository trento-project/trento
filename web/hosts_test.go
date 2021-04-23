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

	filterEnv := &consulApi.KVPair{
		Value: []byte("[\"env1\", \"env2\"]"),
	}
	filterLand := &consulApi.KVPair{
		Value: []byte("[\"land1\", \"land2\"]"),
	}
	filterSys := &consulApi.KVPair{
		Value: []byte("[\"sys1\", \"sys2\"]"),
	}

	consul := new(mocks.Client)
	catalog := new(mocks.Catalog)
	health := new(mocks.Health)
	kv := new(mocks.KV)

	consul.On("Catalog").Return(catalog)
	consul.On("Health").Return(health)
	consul.On("KV").Return(kv)

	query := &consulApi.QueryOptions{Filter: ""}
	catalog.On("Nodes", (*consulApi.QueryOptions)(query)).Return(nodes, nil, nil)

	health.On("Node", "foo", (*consulApi.QueryOptions)(nil)).Return(fooHealthChecks, nil, nil)
	health.On("Node", "bar", (*consulApi.QueryOptions)(nil)).Return(barHealthChecks, nil, nil)

	kv.On("Get", "trento/filters/sap-environments", (*consulApi.QueryOptions)(nil)).Return(filterEnv, nil, nil)
	kv.On("Get", "trento/filters/sap-landscapes", (*consulApi.QueryOptions)(nil)).Return(filterLand, nil, nil)
	kv.On("Get", "trento/filters/sap-systems", (*consulApi.QueryOptions)(nil)).Return(filterSys, nil, nil)

	deps := DefaultDependencies()
	deps.consul = consul

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

	consul.AssertExpectations(t)
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
	assert.Regexp(t, regexp.MustCompile("<select name=trento-sap-landscape.*>.*land1.*land2.*</select>"), minified)
	assert.Regexp(t, regexp.MustCompile("<select name=trento-sap-system.*>.*sys1.*sys2.*</select>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>foo</td><td>192.168.1.1</td><td>.*land1.*</td><td>.*passing.*</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>bar</td><td>192.168.1.2</td><td>.*land2.*</td><td>.*critical.*</td>"), minified)
}
