package web

import (
	"fmt"
	"net/http/httptest"
	"regexp"
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	consulMocks "github.com/trento-project/trento/internal/consul/mocks"
	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

func hostListFixture() models.HostList {
	return models.HostList{
		{
			ID:            "1",
			Name:          "host1",
			IPAddresses:   []string{"192.168.1.1"},
			CloudProvider: "azure",
			SAPSystems: []*models.SAPSystem{
				{
					ID:   "sap_system_id_1",
					SID:  "PRD",
					Type: "database",
				},
			},
			AgentVersion: "v1",
			Tags:         []string{"tag1"},
			Health:       "passing",
		},
		{
			ID:            "2",
			Name:          "host2",
			IPAddresses:   []string{"192.168.1.2"},
			CloudProvider: "aws",
			SAPSystems: []*models.SAPSystem{
				{
					ID:   "sap_system_id_2",
					SID:  "QAS",
					Type: "application",
				},
			},
			AgentVersion: "v1",
			Tags:         []string{"tag2"},
			Health:       "warning",
		},
		{
			ID:            "1",
			Name:          "host3",
			IPAddresses:   []string{"192.168.1.3"},
			CloudProvider: "gcp",
			SAPSystems: []*models.SAPSystem{
				{
					ID:   "sap_system_id_3",
					SID:  "DEV",
					Type: "application",
				},
			},
			AgentVersion: "v1",
			Tags:         []string{"tag3"},
			Health:       "critical",
		},
	}
}

func TestNewHostsHealthContainer(t *testing.T) {
	hCont := NewHostsHealthContainer(hostListFixture())

	expectedHealth := &HealthContainer{
		PassingCount:  1,
		WarningCount:  1,
		CriticalCount: 1,
	}

	assert.Equal(t, expectedHealth, hCont)
}

func TestHostListNextHandler(t *testing.T) {

	mockHostsService := new(services.MockHostsService)
	mockHostsService.On("GetAll", mock.Anything, mock.Anything).Return(hostListFixture(), nil)
	mockHostsService.On("GetCount").Return(3, nil)
	mockHostsService.On("GetAllSIDs", mock.Anything).Return([]string{"PRD", "QAS", "DEV"}, nil)
	mockHostsService.On("GetAllTags", mock.Anything).Return([]string{"tag1", "tag2", "tag3"}, nil)

	deps := setupTestDependencies()
	deps.hostsService = mockHostsService

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/hosts", nil)

	app.webEngine.ServeHTTP(resp, req)

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

	assert.Regexp(t, regexp.MustCompile("<select name=sids.*>.*PRD.*QAS.*DEV.*</select>"), minified)
	assert.Regexp(t, regexp.MustCompile(".*check_circle.*<td>.*host1.*</td><td>192.168.1.1</td><td>.*azure.*</td><td>.*databases/sap_system_id_1.*PRD.*</td><td>v1</td><td>.*<input.*value=tag1.*>.*</td>"), minified)
	assert.Regexp(t, regexp.MustCompile(".*warning.*<td>.*host2.*</td><td>192.168.1.2</td><td>.*aws.*</td><td>.*sapsystems/sap_system_id_2.*QAS.*</td><td>v1</td><td>.*<input.*value=tag2.*>.*</td>"), minified)
	assert.Regexp(t, regexp.MustCompile(".*error.*<td>.*host3.*</td><td>192.168.1.3</td><td>.*gcp.*</td><td>.*sapsystems/sap_system_id_3.*DEV.*</td><td>v1</td><td>.*<input.*value=tag3.*>.*</td>"), minified)
}

func TestApiHostHeartbeat(t *testing.T) {
	agentID := "agent_id"

	mockHostsService := new(services.MockHostsService)
	mockHostsService.On("Heartbeat", agentID).Return(nil)

	deps := setupTestDependencies()
	deps.hostsService = mockHostsService

	app, err := NewAppWithDeps(setupTestConfig(), deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	url := fmt.Sprintf("/api/hosts/%s/heartbeat", agentID)
	req := httptest.NewRequest("POST", url, nil)

	app.collectorEngine.ServeHTTP(resp, req)

	assert.Equal(t, 204, resp.Code)
}

func TestHostHandler(t *testing.T) {
	consulInst := new(consulMocks.Client)
	catalog := new(consulMocks.Catalog)
	health := new(consulMocks.Health)
	kv := new(consulMocks.KV)
	subscriptionsMocks := new(services.MockSubscriptionsService)

	consulInst.On("Catalog").Return(catalog)
	consulInst.On("Health").Return(health)
	consulInst.On("KV").Return(kv)

	node := &consulApi.Node{
		Node:       "test_host",
		Datacenter: "dc1",
		Address:    "192.168.1.1",
		Meta: map[string]string{
			"trento-sap-systems":    "sys1",
			"trento-sap-systems-id": "123456",
			"trento-agent-version":  "1",
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
			CheckID: "trentoAgent",
			Status:  consulApi.HealthPassing,
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

	subscriptionsList := []*models.SlesSubscription{
		&models.SlesSubscription{
			ID:                 "SLES_SAP",
			Version:            "15.2",
			Arch:               "x64_84",
			Status:             "Registered",
			StartsAt:           "2019-03-20 09:55:32 UTC",
			ExpiresAt:          "2024-03-20 09:55:32 UTC",
			SubscriptionStatus: "ACTIVE",
			Type:               "internal",
		},
		&models.SlesSubscription{
			ID:      "sle-module-desktop-applications",
			Version: "15.2",
			Arch:    "x64_84",
			Status:  "Registered",
		},
	}

	subscriptionsMocks.On("GetHostSubscriptions", "test_host").Return(subscriptionsList, nil)

	deps := setupTestDependencies()
	deps.consul = consulInst
	deps.subscriptionsService = subscriptionsMocks

	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/hosts/test_host", nil)
	req.Header.Set("Accept", "text/html")

	app.webEngine.ServeHTTP(resp, req)

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
	assert.Regexp(t, regexp.MustCompile("<a.*sapsystems/123456.*>sys1</a>"), minified)
	assert.Regexp(t, regexp.MustCompile("<dd.*>v1</dd>"), minified)
	assert.Regexp(t, regexp.MustCompile("<span.*>passing</span>"), minified)
	// Subscriptions
	assert.Regexp(t, regexp.MustCompile(
		"<td>SLES_SAP</td><td>x64_84</td><td>15.2</td><td>internal</td><td>Registered</td>"+
			"<td>ACTIVE</td><td>2019-03-20 09:55:32 UTC</td><td>2024-03-20 09:55:32 UTC</td>"), minified)
	assert.Regexp(t, regexp.MustCompile(
		"<td>sle-module-desktop-applications</td><td>x64_84</td><td>15.2</td><td></td>"+
			"<td>Registered</td><td></td><td></td><td></td>"), minified)
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
	consulInst := new(consulMocks.Client)
	catalog := new(consulMocks.Catalog)
	health := new(consulMocks.Health)
	kv := new(consulMocks.KV)
	subscriptionsMocks := new(services.MockSubscriptionsService)

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
	subscriptionsMocks.On(
		"GetHostSubscriptions", "test_host").Return([]*models.SlesSubscription{}, nil)

	deps := setupTestDependencies()
	deps.consul = consulInst
	deps.subscriptionsService = subscriptionsMocks

	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/hosts/test_host", nil)
	req.Header.Set("Accept", "text/html")

	app.webEngine.ServeHTTP(resp, req)

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
	consulInst := new(consulMocks.Client)
	catalog := new(consulMocks.Catalog)
	catalog.On("Node", "foobar", (*consulApi.QueryOptions)(nil)).Return((*consulApi.CatalogNode)(nil), nil, nil)
	consulInst.On("Catalog").Return(catalog)

	deps := setupTestDependencies()
	deps.consul = consulInst

	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/hosts/foobar", nil)
	req.Header.Set("Accept", "text/html")

	app.webEngine.ServeHTTP(resp, req)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.Code)
	assert.Contains(t, resp.Body.String(), "Not Found")
}
