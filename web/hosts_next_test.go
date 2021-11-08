package web

import (
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

func TestHostsListNextHandler(t *testing.T) {
	hostList := models.HostList{
		{
			ID:            "1",
			Name:          "host1",
			IPAddresses:   []string{"192.168.1.1"},
			CloudProvider: "azure",
			SIDs:          []string{"PRD"},
			AgentVersion:  "v1",
			Tags:          []string{"tag1"},
		},
		{
			ID:            "2",
			Name:          "host2",
			IPAddresses:   []string{"192.168.1.2"},
			CloudProvider: "aws",
			SIDs:          []string{"QAS"},
			AgentVersion:  "v1",
			Tags:          []string{"tag2"},
		},
		{
			ID:            "1",
			Name:          "host3",
			IPAddresses:   []string{"192.168.1.3"},
			CloudProvider: "gcp",
			SIDs:          []string{"DEV"},
			AgentVersion:  "v1",
			Tags:          []string{"tag3"},
		},
	}

	mockHostsService := new(services.MockHostsNextService)
	mockHostsService.On("GetAll", mock.Anything).Return(hostList, nil)
	mockHostsService.On("GetAllSIDs", mock.Anything).Return([]string{"PRD", "QAS", "DEV"}, nil)
	mockHostsService.On("GetAllTags", mock.Anything).Return([]string{"tag1", "tag2", "tag3"}, nil)

	deps := setupTestDependencies()
	deps.hostsNextService = mockHostsService

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/hosts-next", nil)

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

	// TODO: test sap systems link and health
	assert.Regexp(t, regexp.MustCompile("<select name=sids.*>.*PRD.*QAS.*DEV.*</select>"), minified)
	assert.Regexp(t, regexp.MustCompile("</td><td>.*host1.*</td><td>192.168.1.1</td><td>.*azure.*</td><td>.*PRD.*</td><td>v1</td><td>.*<input.*value=tag1.*>.*</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>.*host2.*</td><td>192.168.1.2</td><td>.*aws.*</td><td>.*QAS.*</td><td>v1</td><td>.*<input.*value=tag2.*>.*</td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td>.*host3.*</td><td>192.168.1.3</td><td>.*gcp.*</td><td>.*DEV.*</td><td>v1</td><td>.*<input.*value=tag3.*>.*</td>"), minified)
}
