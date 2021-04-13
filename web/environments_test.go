package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/internal/consul/mocks"
)

func TestEnvironmentsListHandler(t *testing.T) {
	datacenters := []string{"test-environment"}
	nodes := []*consulApi.Node{
		{
			Node:       "foo",
			Datacenter: "test-environment",
		},
		{
			Node:       "bar",
			Datacenter: "test-environment",
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

	consul := new(mocks.Client)
	catalog := new(mocks.Catalog)
	health := new(mocks.Health)

	consul.On("Catalog").Return(catalog)
	consul.On("Health").Return(health)

	catalog.On("Datacenters").Return(datacenters, nil)
	catalog.On("Nodes", (*consulApi.QueryOptions)(nil)).Return(nodes, nil, nil)

	health.On("Node", "foo", (*consulApi.QueryOptions)(nil)).Return(fooHealthChecks, nil, nil)
	health.On("Node", "bar", (*consulApi.QueryOptions)(nil)).Return(barHealthChecks, nil, nil)

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
	health.AssertExpectations(t)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "<span class=\"environment-name\">test-environment</span>")
	assert.Contains(t, resp.Body.String(), "<span class=\"node-name\">foo</span>")
	assert.Contains(t, resp.Body.String(), "<span class=\"node-name\">bar</span>")
	assert.Contains(t, resp.Body.String(), "<span class=\"node-health health-passing\">passing</span>")
	assert.Contains(t, resp.Body.String(), "<span class=\"node-health health-critical\">critical</span>")
}
