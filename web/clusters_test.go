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
	consul_internal "github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/consul/mocks"
)

func TestClustersListHandler(t *testing.T) {
	clusters := consulApi.KVPairs{
		&consulApi.KVPair{
			Key: consul_internal.KvClustersPath + "/",
		},
		&consulApi.KVPair{
			Key: consul_internal.KvClustersPath + "/cluster1/",
		},
		&consulApi.KVPair{
			Key: consul_internal.KvClustersPath + "/cluster2/",
		},
	}

	consul := new(mocks.Client)
	kv := new(mocks.KV)

	consul.On("KV").Return(kv)

	kv.On("List", consul_internal.KvClustersPath+"/", (*consulApi.QueryOptions)(nil)).Return(clusters, nil, nil)

	deps := DefaultDependencies()
	deps.consul = consul

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/clusters", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	consul.AssertExpectations(t)
	kv.AssertExpectations(t)

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
	assert.Contains(t, minified, "Clusters")
	assert.Regexp(t, regexp.MustCompile("<td>cluster1</td><td>2</td><td>4</td><td>.*passing.*</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>cluster2</td><td>2</td><td>4</td><td>.*passing.*</td>"), minified)
}
