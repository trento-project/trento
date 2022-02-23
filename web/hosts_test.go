package web

import (
	"fmt"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
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
			CloudData: models.AzureCloudData{
				VMName:          "host1",
				ResourceGroup:   "carbonara-resourcegroup",
				Location:        "southern-rome",
				VMSize:          "extra-large",
				DataDisksNumber: 8,
				Offer:           "sales",
				SKU:             "skeks",
				AdminUsername:   "toor",
			},
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
	assert.Regexp(t, regexp.MustCompile(".*check_circle.*<td .*>.*host1.*</td><td>192.168.1.1</td><td>.*azure.*</td><td>.*databases/sap_system_id_1.*PRD.*</td><td>v1</td><td .*>.*<input.*value=tag1.*>.*</td>"), minified)
	assert.Regexp(t, regexp.MustCompile(".*warning.*<td .*>.*host2.*</td><td>192.168.1.2</td><td>.*aws.*</td><td>.*sapsystems/sap_system_id_2.*QAS.*</td><td>v1</td><td .*>.*<input.*value=tag2.*>.*</td>"), minified)
	assert.Regexp(t, regexp.MustCompile(".*error.*<td .*>.*host3.*</td><td>192.168.1.3</td><td>.*gcp.*</td><td>.*sapsystems/sap_system_id_3.*DEV.*</td><td>v1</td><td .*>.*<input.*value=tag3.*>.*</td>"), minified)
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
	subscriptionsMocks := new(services.MockSubscriptionsService)
	mockHostsService := new(services.MockHostsService)

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

	exportersState := map[string]string{
		"Node exporter":  "passing",
		"Other exporter": "critical",
	}

	subscriptionsMocks.On("GetHostSubscriptions", "2").Return(subscriptionsList, nil)
	subscriptionsMocks.On("IsTrentoPremium").Return(true, nil)
	mockHostsService.On("GetByID", "2").Return(hostListFixture()[1], nil)
	mockHostsService.On("GetExportersState", "host2").Return(exportersState, nil)

	deps := setupTestDependencies()
	deps.subscriptionsService = subscriptionsMocks
	deps.hostsService = mockHostsService

	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/hosts/2", nil)
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

	assert.Regexp(t, regexp.MustCompile("<span.*>host2</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<a.*sapsystems/sap_system_id_2.*>QAS</a>"), minified)
	assert.Regexp(t, regexp.MustCompile("<span.*>v1</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>Trento agent</td><td><span.*>not running</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>Node exporter</td><td><span.*>running</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>Other exporter</td><td><span.*>not running</span>"), minified)

	// Subscriptions
	assert.Regexp(t, regexp.MustCompile(
		"<td>SLES_SAP</td><td>x64_84</td><td>15.2</td><td>internal</td><td>Registered</td>"+
			"<td>ACTIVE</td><td>2019-03-20 09:55:32 UTC</td><td>2024-03-20 09:55:32 UTC</td>"), minified)
	assert.Regexp(t, regexp.MustCompile(
		"<td>sle-module-desktop-applications</td><td>x64_84</td><td>15.2</td><td></td>"+
			"<td>Registered</td><td></td><td></td><td></td>"), minified)
}

func TestHostHandlerAzure(t *testing.T) {
	subscriptionsMocks := new(services.MockSubscriptionsService)
	mockHostsService := new(services.MockHostsService)

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

	subscriptionsMocks.On("GetHostSubscriptions", "1").Return(subscriptionsList, nil)
	subscriptionsMocks.On("IsTrentoPremium").Return(true, nil)
	mockHostsService.On("GetByID", "1").Return(hostListFixture()[0], nil)
	mockHostsService.On("GetExportersState", "host1").Return(make(map[string]string), nil)

	deps := setupTestDependencies()
	deps.subscriptionsService = subscriptionsMocks
	deps.hostsService = mockHostsService

	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/hosts/1", nil)
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

	assert.Regexp(t, regexp.MustCompile("<span.*>host1</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<a.*sapsystems/sap_system_id_1.*>PRD</a>"), minified)
	assert.Regexp(t, regexp.MustCompile("<span.*>v1</span>"), minified)
	assert.Regexp(t, regexp.MustCompile("<span.*>running</span>"), minified)

	// Subscriptions
	assert.Regexp(t, regexp.MustCompile(
		"<td>SLES_SAP</td><td>x64_84</td><td>15.2</td><td>internal</td><td>Registered</td>"+
			"<td>ACTIVE</td><td>2019-03-20 09:55:32 UTC</td><td>2024-03-20 09:55:32 UTC</td>"), minified)
	assert.Regexp(t, regexp.MustCompile(
		"<td>sle-module-desktop-applications</td><td>x64_84</td><td>15.2</td><td></td>"+
			"<td>Registered</td><td></td><td></td><td></td>"), minified)
}

func TestHostHandler404Error(t *testing.T) {
	subscriptionsMocks := new(services.MockSubscriptionsService)
	mockHostsService := new(services.MockHostsService)

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

	subscriptionsMocks.On("GetHostSubscriptions", "foobar").Return(subscriptionsList, nil)
	subscriptionsMocks.On("IsTrentoPremium").Return(true, nil)
	mockHostsService.On("GetByID", "foobar").Return(nil, nil)
	deps := setupTestDependencies()
	deps.subscriptionsService = subscriptionsMocks
	deps.hostsService = mockHostsService

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
