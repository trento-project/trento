package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/consul/mocks"
	"github.com/trento-project/trento/internal/sapsystem"
)

func sapSystemsMap(instanceName string, host string, instanceNr int, features string) map[string]interface{} {
	return map[string]interface{}{
		"HA1": map[string]interface{}{
			"sid":  "HA1",
			"type": sapsystem.Application,
			"instances": map[string]interface{}{
				instanceName: map[string]interface{}{
					"name": instanceName,
					"host": host,
					"sapcontrol": map[string]interface{}{
						"properties": map[string]interface{}{
							"SAPSYSTEM": map[string]interface{}{
								"Value": fmt.Sprintf("%02d", instanceNr),
							},
						},
						"instances": map[string]interface{}{
							"sapha1as": map[string]interface{}{
								"hostname":      host,
								"instancenr":    instanceNr,
								"features":      features,
								"httpport":      50013,
								"httpsport":     50014,
								"startpriority": "0.5",
								"dispstatus":    "SAPControl-GREEN",
							},
						},
					},
				},
			},
		},
	}
}

func TestSAPSystemsListHandler(t *testing.T) {
	nodes := []*consulApi.Node{
		{
			Node: "netweaver01",
			Meta: map[string]string{
				"trento-ha-cluster":    "banana",
				"trento-ha-cluster-id": "e2f2eb50aef748e586a7baa85e0162cf",
			},
		},
		{
			Node: "netweaver02",
			Meta: map[string]string{
				"trento-ha-cluster":    "banana",
				"trento-ha-cluster-id": "e2f2eb50aef748e586a7baa85e0162cf",
			},
		},
		{
			Node: "netweaver03",
			Meta: map[string]string{
				"trento-ha-cluster":    "banana",
				"trento-ha-cluster-id": "e2f2eb50aef748e586a7baa85e0162cf",
			},
		},
		{
			Node: "netweaver04",
			Meta: map[string]string{
				"trento-ha-cluster":    "banana",
				"trento-ha-cluster-id": "e2f2eb50aef748e586a7baa85e0162cf",
			},
		},
	}

	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	kv := new(mocks.KV)

	catalog.On("Nodes", mock.Anything).Return(nodes, nil, nil)
	consulInst.On("Catalog").Return(catalog)

	consulInst.On("WaitLock", mock.Anything).Return(nil)
	p := fmt.Sprintf(consul.KvHostsSAPSystemPath, "netweaver01")
	m := sapSystemsMap("ERS10", "netweaver01", 10, "ENQREP")
	kv.On("ListMap", p, p).Return(m, nil)

	p = fmt.Sprintf(consul.KvHostsSAPSystemPath, "netweaver02")
	m = sapSystemsMap("ASCS00", "netweaver02", 0, "MESSAGESERVER|ENQUE")
	kv.On("ListMap", p, p).Return(m, nil)

	p = fmt.Sprintf(consul.KvHostsSAPSystemPath, "netweaver03")
	m = sapSystemsMap("D01", "netweaver03", 1, "ABAP|GATEWAY|ICMAN|IGS")
	kv.On("ListMap", p, p).Return(m, nil)

	p = fmt.Sprintf(consul.KvHostsSAPSystemPath, "netweaver04")
	m = sapSystemsMap("D02", "netweaver04", 2, "ABAP|GATEWAY|ICMAN|IGS")
	kv.On("ListMap", p, p).Return(m, nil)

	tags := map[string]interface{}{
		"tag1": struct{}{},
	}
	kv.On("ListMap", "trento/v0/tags/sapsystems/HA1/", "trento/v0/tags/sapsystems/HA1/").Return(tags, nil)

	consulInst.On("KV").Return(kv)

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
	assert.Regexp(t, regexp.MustCompile("<td><a href=/sapsystems/HA1>HA1</a></td><td></td><td>.*<input.*value=tag1.*>.*</td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<td>HA1</td><td>ENQREP</td><td>10</td><td><a href=/clusters/e2f2eb50aef748e586a7baa85e0162cf>banana</a></td><td><a href=/hosts/netweaver01>netweaver01</a></td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<td>HA1</td><td>MESSAGESERVER|ENQUE</td><td>00</td><td><a href=/clusters/e2f2eb50aef748e586a7baa85e0162cf>banana</a></td><td><a href=/hosts/netweaver02>netweaver02</a></td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<td>HA1</td><td>ABAP|GATEWAY|ICMAN|IGS</td><td>01</td><td><a href=/clusters/e2f2eb50aef748e586a7baa85e0162cf>banana</a></td><td><a href=/hosts/netweaver03>netweaver03</a></td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<td>HA1</td><td>ABAP|GATEWAY|ICMAN|IGS</td><td>02</td><td><a href=/clusters/e2f2eb50aef748e586a7baa85e0162cf>banana</a></td><td><a href=/hosts/netweaver04>netweaver04</a></td>"), responseBody)
}

func TestSAPSystemHandler(t *testing.T) {
	nodes := []*consulApi.Node{
		{
			Node:    "test_host",
			Address: "192.168.10.10",
			Meta: map[string]string{
				"trento-ha-cluster":     "banana",
				"trento-ha-cluster-id":  "e2f2eb50aef748e586a7baa85e0162cf",
				"trento-cloud-provider": "azure",
				"trento-sap-systems":    "HA1",
				"trento-agent-version":  "0",
			},
		},
	}

	passHealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthPassing,
		},
	}

	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	kv := new(mocks.KV)
	health := new(mocks.Health)

	catalog.On("Nodes", mock.Anything).Return(nodes, nil, nil)
	consulInst.On("Catalog").Return(catalog)

	consulInst.On("WaitLock", mock.Anything).Return(nil)
	p := fmt.Sprintf(consul.KvHostsSAPSystemPath, "test_host")
	m := sapSystemsMap("ERS10", "test_host", 10, "ENQREP")
	kv.On("ListMap", p, p).Return(m, nil)
	consulInst.On("KV").Return(kv)

	health.On("Node", "test_host", mock.Anything).Return(passHealthChecks, nil, nil)
	consulInst.On("Health").Return(health)

	deps := DefaultDependencies()
	deps.consul = consulInst

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/sapsystems/HA1", nil)
	if err != nil {
		t.Fatal(err)
	}
	app.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
	responseBody := minifyHtml(resp.Body.String())

	assert.Contains(t, responseBody, "SAP System details")
	assert.Contains(t, responseBody, "HA1")
	// Layout/
	assert.Regexp(t, regexp.MustCompile("<tr><td>test_host</td><td>10</td><td>ENQREP</td><td>50013</td><td>50014</td><td>0.5</td><td><span.*primary.*>SAPControl-GREEN</span></td></tr>"), responseBody)
	// Host
	assert.Regexp(t, regexp.MustCompile("<tr><td>.*check_circle.*</td><td><a href=/hosts/test_host>test_host</a></td><td>192.168.10.10</td><td>azure</td><td><a href=/clusters/e2f2eb50aef748e586a7baa85e0162cf>banana</a></td><td><a href=/sapsystems/HA1>HA1</a></td><td>v0</td></tr>"), responseBody)
}

func TestSAPSystemHandler404Error(t *testing.T) {
	nodes := []*consulApi.Node{
		{
			Node: "test_host",
			Meta: map[string]string{
				"trento-ha-cluster":    "banana",
				"trento-ha-cluster-id": "e2f2eb50aef748e586a7baa85e0162cf",
			},
		},
	}

	consulInst := new(mocks.Client)
	kv := new(mocks.KV)
	catalog := new(mocks.Catalog)

	catalog.On("Nodes", mock.Anything).Return(nodes, nil, nil)
	consulInst.On("Catalog").Return(catalog)

	consulInst.On("WaitLock", mock.Anything).Return(nil)
	p := fmt.Sprintf(consul.KvHostsSAPSystemPath, "test_host")
	m := sapSystemsMap("ERS10", "test_host", 10, "ENQREP")
	kv.On("ListMap", p, p).Return(m, nil)
	consulInst.On("KV").Return(kv)

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
