package web

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/mock"
	consulMock "github.com/trento-project/trento/internal/consul/mocks"
	hostsServiceMock "github.com/trento-project/trento/web/service/mocks"

	"github.com/trento-project/trento/web/models"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

func TestNewHostsHealthContainer(t *testing.T) {
	hosts := []models.Host{
		{Health: consulApi.HealthPassing},
		{Health: "passing"},
		{Health: consulApi.HealthWarning},
		{Health: "warning"},
		{Health: consulApi.HealthCritical},
		{Health: "critical"},
	}

	hCont := NewHostsHealthContainer(hosts)

	expectedHealth := &HealthContainer{
		PassingCount:  2,
		WarningCount:  2,
		CriticalCount: 2,
	}

	assert.Equal(t, expectedHealth, hCont)
}

func TestHostsListHandler(t *testing.T) {
	hosts := []models.Host{
		{
			Name:          "foo",
			Address:       "192.168.1.1",
			Health:        consulApi.HealthPassing,
			Environment:   "env1",
			Landscape:     "land1",
			SAPSystem:     "sys1",
			CloudProvider: "azure",
		},
		{
			Name:          "bar",
			Address:       "192.168.1.2",
			Health:        consulApi.HealthCritical,
			Environment:   "env2",
			Landscape:     "land2",
			SAPSystem:     "sys2",
			CloudProvider: "aws",
		},
		{
			Name:          "buzz",
			Address:       "192.168.1.3",
			Health:        consulApi.HealthWarning,
			Environment:   "env3",
			Landscape:     "land3",
			SAPSystem:     "sys3",
			CloudProvider: "gcp",
		},
	}

	hostServiceMock := new(hostsServiceMock.IHostsService)
	hostServiceMock.On("GetHosts", mock.Anything, mock.Anything).Return(hosts)
	hostServiceMock.On("GetHostsCount").Return(len(hosts))
	hostServiceMock.On("GetHostsSAPSystems").Return([]string{"sys1", "sys2", "sys3"})
	hostServiceMock.On("GetHostsLandscapes").Return([]string{"land1", "land2", "land3"})
	hostServiceMock.On("GetHostsEnvironments").Return([]string{"env1", "env2", "env3"})
	deps := Dependencies{
		hostsService: hostServiceMock,
		engine:       gin.Default(),
	}

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
	assert.Regexp(t, regexp.MustCompile("<div.*alert-success.*<i.*check_circle.*</i>.*Passing.*1"), minified)
	assert.Regexp(t, regexp.MustCompile("<div.*alert-warning.*<i.*warning.*</i>.*Warning.*1"), minified)
	assert.Regexp(t, regexp.MustCompile("<div.*alert-danger.*<i.*error.*</i>.*Critical.*1"), minified)
	assert.Regexp(t, regexp.MustCompile("<td.*<i.*success.*check_circle.*</i></td><td>.*foo.*</td><td>192.168.1.1</td><td>.*azure.*</td><td>.*sys1.*</td><td>.*land1.*</td><td>.*env1.*</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td.*<i.*success.*check_circle.*</i></td><td>.*foo.*</td><td>192.168.1.1</td><td>.*azure.*</td><td>.*sys1.*</td><td>.*land1.*</td><td>.*env1.*</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<select name=environment.*>.*env1.*env2.*</select>"), minified)
	assert.Regexp(t, regexp.MustCompile("<select name=landscape.*>.*land1.*land2.*land3.*</select>"), minified)
	assert.Regexp(t, regexp.MustCompile("<select name=sap_system.*>.*sys1.*sys2.*sys3.*</select>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td.*<i.*success.*check_circle.*</i></td><td>.*foo.*</td><td>192.168.1.1</td><td>.*azure.*</td><td>.*sys1.*</td><td>.*land1.*</td><td>.*env1.*</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td.*<i.*danger.*error.*</i></td><td>.*bar.*</td><td>192.168.1.2</td><td>.*aws.*</td><td>.*sys2.*</td><td>.*land2.*</td><td>.*env2.*</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td.*<i.*warning.*warning.*</i></td><td>.*buzz.*</td><td>192.168.1.3</td><td>.*gcp.*</td><td>.*sys3.*</td><td>.*land3.*</td><td>.*env3.*</td>"), minified)
}

func TestHostHandler(t *testing.T) {
	consulInst := new(consulMock.Client)
	catalog := new(consulMock.Catalog)
	health := new(consulMock.Health)
	kv := new(consulMock.KV)

	consulInst.On("Catalog").Return(catalog)
	consulInst.On("Health").Return(health)
	consulInst.On("KV").Return(kv)

	node := &consulApi.Node{
		Node:       "test_host",
		Datacenter: "dc1",
		Address:    "192.168.1.1",
		Meta: map[string]string{
			"trento-sap-environment": "env1",
			"trento-sap-system":      "sys1",
			"trento-sap-landscape":   "land1",
		},
	}

	catalogNode := &consulApi.CatalogNode{Node: node}
	catalog.On("Node", "test_host", (*consulApi.QueryOptions)(nil)).Return(catalogNode, nil, nil)

	healthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthPassing,
		},
	}
	health.On("Node", "test_host", (*consulApi.QueryOptions)(nil)).Return(healthChecks, nil, nil)

	sapsystemPath := "trento/v0/hosts/test_host/sapsystems/"
	consulInst.On("WaitLock", sapsystemPath).Return(nil)
	kv.On("ListMap", sapsystemPath, sapsystemPath).Return(nil, nil)

	cloudListMap := map[string]interface{}{
		"provider": "other",
	}
	cloudPath := "trento/v0/hosts/test_host/"
	cloudListMapPath := cloudPath + "cloud/"
	consulInst.On("WaitLock", cloudPath).Return(nil)
	kv.On("ListMap", cloudListMapPath, cloudListMapPath).Return(cloudListMap, nil)

	deps := DefaultDependencies()
	deps.consul = consulInst

	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/hosts/test_host", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

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
	assert.Contains(t, minified, "Host details")
	assert.Regexp(t, regexp.MustCompile("<dd.*>test_host</dd>"), minified)
	assert.Regexp(t, regexp.MustCompile("<a.*environments.*>env1</a>"), minified)
	assert.Regexp(t, regexp.MustCompile("<a.*landscapes.*>land1</a>"), minified)
	assert.Regexp(t, regexp.MustCompile("<a.*sapsystems.*>sys1</a>"), minified)
	assert.Regexp(t, regexp.MustCompile("<span.*>passing</span>"), minified)
}

func TestHostHandlerAzure(t *testing.T) {
	consulInst := new(consulMock.Client)
	catalog := new(consulMock.Catalog)
	health := new(consulMock.Health)
	kv := new(consulMock.KV)

	consulInst.On("Catalog").Return(catalog)
	consulInst.On("Health").Return(health)
	consulInst.On("KV").Return(kv)

	node := &consulApi.Node{
		Node:       "test_host",
		Datacenter: "dc1",
		Address:    "192.168.1.1",
		Meta: map[string]string{
			"trento-sap-environment": "env1",
			"trento-sap-system":      "sys1",
			"trento-sap-landscape":   "land1",
		},
	}

	catalogNode := &consulApi.CatalogNode{Node: node}
	catalog.On("Node", "test_host", (*consulApi.QueryOptions)(nil)).Return(catalogNode, nil, nil)

	healthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthPassing,
		},
	}
	health.On("Node", "test_host", (*consulApi.QueryOptions)(nil)).Return(healthChecks, nil, nil)

	sapsystemPath := "trento/v0/hosts/test_host/sapsystems/"
	consulInst.On("WaitLock", sapsystemPath).Return(nil)
	kv.On("ListMap", sapsystemPath, sapsystemPath).Return(nil, nil)

	cloudListMap := map[string]interface{}{
		"provider": "azure",
		"metadata": map[string]interface{}{
			"compute": map[string]interface{}{
				"name":     "vmtest_host",
				"location": "north",
				"vmsize":   "10gb",
				"storageprofile": map[string]interface{}{
					"datadisks": []interface{}{
						map[string]interface{}{
							"name": "value1",
						},
						map[string]interface{}{
							"name": "value2",
						},
					},
				},
				"offer":             "superoffer",
				"sku":               "gen2",
				"subscription":      "1234",
				"resourceid":        "resource1",
				"resourcegroupname": "group1",
			},
		},
	}
	cloudPath := "trento/v0/hosts/test_host/"
	cloudListMapPath := cloudPath + "cloud/"
	consulInst.On("WaitLock", cloudPath).Return(nil)
	kv.On("ListMap", cloudListMapPath, cloudListMapPath).Return(cloudListMap, nil)

	deps := DefaultDependencies()
	deps.consul = consulInst

	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/hosts/test_host", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

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
	assert.Contains(t, minified, "Cloud details")
	assert.Regexp(t, regexp.MustCompile("<dd.*>.*vmtest_host.*</dd>"), minified)
	assert.Regexp(t, regexp.MustCompile("<dd.*>north</dd>"), minified)
	assert.Regexp(t, regexp.MustCompile("<dd.*>10gb</dd>"), minified)
	assert.Regexp(t, regexp.MustCompile("<dd.*>2</dd>"), minified)
	assert.Regexp(t, regexp.MustCompile("<dd.*>superoffer</dd>"), minified)
	assert.Regexp(t, regexp.MustCompile("<dd.*>gen2</dd>"), minified)
}

func TestHostHandler404Error(t *testing.T) {
	consulInst := new(consulMock.Client)
	catalog := new(consulMock.Catalog)
	catalog.On("Node", "foobar", (*consulApi.QueryOptions)(nil)).Return((*consulApi.CatalogNode)(nil), nil, nil)
	consulInst.On("Catalog").Return(catalog)

	deps := DefaultDependencies()
	deps.consul = consulInst

	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/hosts/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.Code)
	assert.Contains(t, resp.Body.String(), "Not Found")
}

func TestHAChecksHandler404(t *testing.T) {
	var err error

	consulInst := new(consulMock.Client)
	catalog := new(consulMock.Catalog)
	catalog.On("Node", "foobar", (*consulApi.QueryOptions)(nil)).Return((*consulApi.CatalogNode)(nil), nil, nil)
	consulInst.On("Catalog").Return(catalog)

	deps := DefaultDependencies()
	deps.consul = consulInst

	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/hosts/foobar/ha-checks", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.Code)
	assert.Contains(t, resp.Body.String(), "Not Found")
}
