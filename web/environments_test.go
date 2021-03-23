package web

import (
/*
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/test/mock_consul"
*/
)

/*
func TestEnvironmentsListHandler(t *testing.T) {
	datacenters := []string{"test-environment"}
	nodes := []*api.Node{
		{
			Node: "foo",
			Datacenter: "test-environment",
		},
		{
			Node: "bar",
			Datacenter: "test-environment",
		},
	}

	ctrl := gomock.NewController(t)
	consul := mock_consul.NewMockClient(ctrl)
	catalog := mock_consul.NewMockCatalog(ctrl)
	consul.EXPECT().Catalog().Return(catalog).AnyTimes()
	catalog.EXPECT().Datacenters().Return(datacenters, nil)
	catalog.EXPECT().Nodes(nil).Return(nodes, nil, nil)
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
}
*/
