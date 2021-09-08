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
	"github.com/trento-project/trento/internal/hosts"
)

func TestNewHostsHealthContainer(t *testing.T) {
	consulInst := new(mocks.Client)
	health := new(mocks.Health)
	consulInst.On("Health").Return(health)

	host1 := hosts.NewHost(consulApi.Node{Node: "node1"}, consulInst)
	host2 := hosts.NewHost(consulApi.Node{Node: "node2"}, consulInst)
	host3 := hosts.NewHost(consulApi.Node{Node: "node3"}, consulInst)
	host4 := hosts.NewHost(consulApi.Node{Node: "node4"}, consulInst)
	host5 := hosts.NewHost(consulApi.Node{Node: "node5"}, consulInst)
	host6 := hosts.NewHost(consulApi.Node{Node: "node6"}, consulInst)

	nodes := hosts.HostList{
		&host1, &host2, &host3, &host4, &host5, &host6,
	}

	passHealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthPassing,
		},
	}

	warningHealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthCritical,
		},
	}

	criticalHealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthWarning,
		},
	}

	health.On("Node", "node1", (*consulApi.QueryOptions)(nil)).Return(passHealthChecks, nil, nil)
	health.On("Node", "node2", (*consulApi.QueryOptions)(nil)).Return(warningHealthChecks, nil, nil)
	health.On("Node", "node3", (*consulApi.QueryOptions)(nil)).Return(criticalHealthChecks, nil, nil)
	health.On("Node", "node4", (*consulApi.QueryOptions)(nil)).Return(passHealthChecks, nil, nil)
	health.On("Node", "node5", (*consulApi.QueryOptions)(nil)).Return(warningHealthChecks, nil, nil)
	health.On("Node", "node6", (*consulApi.QueryOptions)(nil)).Return(criticalHealthChecks, nil, nil)

	hCont := NewHostsHealthContainer(nodes)

	expectedHealth := &HealthContainer{
		PassingCount:  2,
		WarningCount:  2,
		CriticalCount: 2,
	}

	assert.Equal(t, expectedHealth, hCont)
}

func TestHostsListHandler(t *testing.T) {
	nodes := []*consulApi.Node{
		{
			Node:       "foo",
			Datacenter: "dc1",
			Address:    "192.168.1.1",
			Meta: map[string]string{
				"trento-sap-systems":    "sys1",
				"trento-cloud-provider": "azure",
				"trento-agent-version":  "1",
			},
		},
		{
			Node:       "bar",
			Datacenter: "dc",
			Address:    "192.168.1.2",
			Meta: map[string]string{
				"trento-sap-systems":    "sys2",
				"trento-cloud-provider": "aws",
				"trento-agent-version":  "1",
			},
		},
		{
			Node:       "buzz",
			Datacenter: "dc",
			Address:    "192.168.1.3",
			Meta: map[string]string{
				"trento-sap-systems":    "sys3",
				"trento-cloud-provider": "gcp",
				"trento-agent-version":  "1",
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

	buzzHealthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Status: consulApi.HealthWarning,
		},
	}

	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	health := new(mocks.Health)
	kv := new(mocks.KV)

	consulInst.On("Catalog").Return(catalog)
	consulInst.On("Health").Return(health)
	consulInst.On("KV").Return(kv)

	query := &consulApi.QueryOptions{Filter: ""}
	catalog.On("Nodes", (*consulApi.QueryOptions)(query)).Return(nodes, nil, nil)

	health.On("Node", "foo", (*consulApi.QueryOptions)(nil)).Return(fooHealthChecks, nil, nil)
	health.On("Node", "bar", (*consulApi.QueryOptions)(nil)).Return(barHealthChecks, nil, nil)
	health.On("Node", "buzz", (*consulApi.QueryOptions)(nil)).Return(buzzHealthChecks, nil, nil)

	kv.On("ListMap", "trento/v0/tags/hosts/foo/", "trento/v0/tags/hosts/foo/").Return(map[string]interface{}{
		"tag1": struct{}{},
	}, nil)
	kv.On("ListMap", "trento/v0/tags/hosts/bar/", "trento/v0/tags/hosts/bar/").Return(nil, nil)
	kv.On("ListMap", "trento/v0/tags/hosts/buzz/", "trento/v0/tags/hosts/buzz/").Return(nil, nil)

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
	assert.Regexp(t, regexp.MustCompile("<div.*alert-success.*<i.*check_circle.*</i>.*Passing.*1"), minified)
	assert.Regexp(t, regexp.MustCompile("<div.*alert-warning.*<i.*warning.*</i>.*Warning.*1"), minified)
	assert.Regexp(t, regexp.MustCompile("<div.*alert-danger.*<i.*error.*</i>.*Critical.*1"), minified)
	assert.Regexp(t, regexp.MustCompile("<td.*<i.*success.*check_circle.*</i></td><td>.*foo.*</td><td>192.168.1.1</td><td>.*azure.*</td><td>.*sys1.*</td><td>v1</td><td>.*<input.*value=tag1.*>.*</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<select name=trento-sap-systems.*>.*sys1.*sys2.*sys3.*</select>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td.*<i.*success.*check_circle.*</i></td><td>.*foo.*</td><td>192.168.1.1</td><td>.*azure.*</td><td>.*sys1.*</td><td>v1</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td.*<i.*danger.*error.*</i></td><td>.*bar.*</td><td>192.168.1.2</td><td>.*aws.*</td><td>.*sys2.*</td><td>v1</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td.*<i.*warning.*warning.*</i></td><td>.*buzz.*</td><td>192.168.1.3</td><td>.*gcp.*</td><td>.*sys3.*</td><td>v1</td>"), minified)
}

func TestHostHandler(t *testing.T) {
	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	health := new(mocks.Health)
	kv := new(mocks.KV)

	consulInst.On("Catalog").Return(catalog)
	consulInst.On("Health").Return(health)
	consulInst.On("KV").Return(kv)

	node := &consulApi.Node{
		Node:       "test_host",
		Datacenter: "dc1",
		Address:    "192.168.1.1",
		Meta: map[string]string{
			"trento-sap-systems":   "sys1",
			"trento-agent-version": "1",
		},
	}

	sapSystemMap := map[string]interface{}{
		"sys1": map[string]interface{}{
			"sid":  "sys1",
			"type": 1,
			"instances": map[string]interface{}{
				"HDB00": map[string]interface{}{
					"name": "HDB00",
					"type": 1,
					"host": "test_host",
					"sapcontrol": map[string]interface{}{
						"properties": map[string]interface{}{
							"INSTANCE_NAME": map[string]interface{}{
								"Value": "HDB00",
							},
							"SAPSYSTEMNAME": map[string]interface{}{
								"Value": "PRD",
							},
							"SAPSYSTEM": map[string]interface{}{
								"Value": "00",
							},
						},
						"processes": map[string]interface{}{
							"proc1": map[string]interface{}{
								"Name":       "proc1",
								"Dispstatus": "SAPControl-GREEN",
								"Textstatus": "Green",
							},
							"proc2": map[string]interface{}{
								"Name":       "proc2",
								"Dispstatus": "SAPControl-YELLOW",
								"Textstatus": "Yellow",
							},
							"proc3": map[string]interface{}{
								"Name":       "proc3",
								"Dispstatus": "SAPControl-RED",
								"Textstatus": "Red",
							},
							"proc4": map[string]interface{}{
								"Name":       "proc4",
								"Dispstatus": "SAPControl-GRAY",
								"Textstatus": "Gray",
							},
						},
					},
				},
			},
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
	kv.On("ListMap", sapsystemPath, sapsystemPath).Return(sapSystemMap, nil)

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
	assert.Regexp(t, regexp.MustCompile("<a.*sapsystems.*>sys1</a>"), minified)
	assert.Regexp(t, regexp.MustCompile("<dd.*>v1</dd>"), minified)
	assert.Regexp(t, regexp.MustCompile("<span.*>passing</span>"), minified)
	// SAP Instance
	assert.Regexp(t, regexp.MustCompile("<dt.*>Name</dt><dd.*>HDB00</dd>"), minified)
	assert.Regexp(t, regexp.MustCompile("<dt.*>SID</dt><dd.*>PRD</dd>"), minified)
	assert.Regexp(t, regexp.MustCompile("<dt.*>Instance number</dt><dd.*>00</dd>"), minified)
	// Processes
	assert.Regexp(t, regexp.MustCompile("<td>proc1</td>.*<span.*primary.*>Green</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>proc2</td>.*<span.*warning.*>Yellow</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>proc3</td>.*<span.*danger.*>Red</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>proc4</td>.*<span.*secondary.*>Gray</span>"), minified)
}

func TestHostHandlerAzure(t *testing.T) {
	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
	health := new(mocks.Health)
	kv := new(mocks.KV)

	consulInst.On("Catalog").Return(catalog)
	consulInst.On("Health").Return(health)
	consulInst.On("KV").Return(kv)

	node := &consulApi.Node{
		Node:       "test_host",
		Datacenter: "dc1",
		Address:    "192.168.1.1",
		Meta: map[string]string{
			"trento-sap-systems": "sys1",
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
	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
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

	consulInst := new(mocks.Client)
	catalog := new(mocks.Catalog)
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
