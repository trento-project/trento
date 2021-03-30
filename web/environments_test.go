package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"

	"github.com/SUSE/console-for-sap-applications/test/mock_consul"
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

	ctrl := gomock.NewController(t)
	consul := mock_consul.NewMockClient(ctrl)
	catalog := mock_consul.NewMockCatalog(ctrl)
	health := mock_consul.NewMockHealth(ctrl)

	consul.EXPECT().Catalog().Return(catalog).AnyTimes()
	consul.EXPECT().Health().Return(health).AnyTimes()

	catalog.EXPECT().Datacenters().Return(datacenters, nil)
	catalog.EXPECT().Nodes(nil).Return(nodes, nil, nil)

	health.EXPECT().Node("foo", nil).Return(fooHealthChecks, nil, nil).AnyTimes()
	health.EXPECT().Node("bar", nil).Return(barHealthChecks, nil, nil).AnyTimes()

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

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "<span class=\"environment-name\">test-environment</span>")
	assert.Contains(t, resp.Body.String(), "<span class=\"node-name\">foo</span>")
	assert.Contains(t, resp.Body.String(), "<span class=\"node-name\">bar</span>")
	assert.Contains(t, resp.Body.String(), "<span class=\"node-health health-passing\">passing</span>")
	assert.Contains(t, resp.Body.String(), "<span class=\"node-health health-critical\">critical</span>")
}
