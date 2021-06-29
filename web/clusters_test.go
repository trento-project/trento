package web

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/consul/mocks"
)

func TestClustersListHandler(t *testing.T) {
	listMap := map[string]interface{}{
		"test_cluster": map[string]interface{}{
			"cib": map[string]interface{}{
				"Configuration": map[string]interface{}{
					"CrmConfig": map[string]interface{}{
						"ClusterProperties": []interface{}{
							map[string]interface{}{
								"Id":    "cib-bootstrap-options-cluster-name",
								"Value": "test_cluster",
							},
						},
					},
				},
			},
			"crmmon": map[string]interface{}{
				"Summary": map[string]interface{}{
					"Nodes": map[string]interface{}{
						"Number": 3,
					},
					"Resources": map[string]interface{}{
						"Number": 5,
					},
				},
			},
			"name": "test_cluster",
		},
		"2nd_cluster": map[string]interface{}{
			"cib": map[string]interface{}{
				"Configuration": map[string]interface{}{
					"CrmConfig": map[string]interface{}{
						"ClusterProperties": []interface{}{
							map[string]interface{}{
								"Id":    "cib-bootstrap-options-cluster-name",
								"Value": "2nd_cluster",
							},
						},
					},
				},
			},
			"crmmon": map[string]interface{}{
				"Summary": map[string]interface{}{
					"Nodes": map[string]interface{}{
						"Number": 2,
					},
					"Resources": map[string]interface{}{
						"Number": 10,
					},
				},
			},
			"name": "2nd_cluster",
		},
	}

	consulInst := new(mocks.Client)
	kv := new(mocks.KV)

	consulInst.On("KV").Return(kv)

	kv.On("ListMap", consul.KvClustersPath, consul.KvClustersPath).Return(listMap, nil)
	consulInst.On("WaitLock", consul.KvClustersPath).Return(nil)

	deps := DefaultDependencies()
	deps.consul = consulInst

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

	consulInst.AssertExpectations(t)
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
	assert.Regexp(t, regexp.MustCompile("<td>test_cluster</td><td></td><td>3</td><td>5</td><td>.*passing.*</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>2nd_cluster</td><td></td><td>2</td><td>10</td><td>.*passing.*</td>"), minified)
}

func TestClusterHandler404Error(t *testing.T) {
	var err error

	kv := new(mocks.KV)
	kv.On("ListMap", consul.KvClustersPath, consul.KvClustersPath).Return(nil, nil)

	consulInst := new(mocks.Client)
	consulInst.On("KV").Return(kv)
	consulInst.On("WaitLock", consul.KvClustersPath).Return(nil)

	deps := DefaultDependencies()
	deps.consul = consulInst

	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/clusters/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.Code)
	assert.Contains(t, resp.Body.String(), "Not Found")
}
