package web

import (
	"net/http"
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

func TestClustersListNextHandler(t *testing.T) {
	clustersList := models.ClusterList{
		{
			ID:                "47d1190ffb4f781974c8356d7f863b03",
			Name:              "hana_cluster",
			ClusterType:       ClusterTypeScaleUp,
			SIDs:              []string{"PRD"},
			ResourcesNumber:   5,
			HostsNumber:       3,
			Tags:              []string{"tag1"},
			Health:            models.CheckPassing,
			HasDuplicatedName: false,
		},
		{
			ID:                "a615a35f65627be5a757319a0741127f",
			Name:              "other_cluster",
			ClusterType:       ClusterTypeUnknown,
			SIDs:              []string{},
			Tags:              []string{"tag1"},
			Health:            models.CheckCritical,
			HasDuplicatedName: false,
		},
		{
			ID:                "e2f2eb50aef748e586a7baa85e0162cf",
			Name:              "netweaver_cluster",
			ClusterType:       ClusterTypeUnknown,
			SIDs:              []string{},
			ResourcesNumber:   10,
			HostsNumber:       2,
			Tags:              []string{"tag1"},
			Health:            models.CheckCritical,
			HasDuplicatedName: true,
		},
		{
			ID:                "e27d313a674375b2066777a89ee346b9",
			Name:              "netweaver_cluster",
			ClusterType:       ClusterTypeUnknown,
			SIDs:              []string{},
			Tags:              []string{"tag1"},
			Health:            models.CheckUndefined,
			HasDuplicatedName: true,
		},
	}

	mockClusterService := new(services.MockClustersService)
	mockClusterService.On("GetAll", mock.Anything).Return(clustersList, nil)

	deps := setupTestDependencies()
	deps.clustersService = mockClusterService

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/clusters-next", nil)
	if err != nil {
		t.Fatal(err)
	}

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
	assert.Contains(t, minified, "Clusters")
	assert.Regexp(t, regexp.MustCompile("<td .*>.*check_circle.*</td><td>.*hana_cluster.*</td><td>.*47d1190ffb4f781974c8356d7f863b03.*</td><td>HANA scale-up</td><td>PRD</td><td>3</td><td>5</td><td><input.*value=tag1.*></td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td .*>.*error.*</td><td>.*other_cluster.*</td><td>.*a615a35f65627be5a757319a0741127f.*</td><td>Unknown</td><td></td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td .*>.*error.*</td><td>.*duplicated.*netweaver_cluster.*</td><td>.*e2f2eb50aef748e586a7baa85e0162cf.*</td><td>Unknown</td><td></td><td>2</td><td>10</td><td><input.*value=tag1.*></td>"), minified)
	assert.Regexp(t, regexp.MustCompile("<td .*>.*fiber_manual_record.*</td><td>.*duplicated.*netweaver_cluster.*</td><td>.*e27d313a674375b2066777a89ee346b9.*</td><td>Unknown</td><td></td>"), minified)
}
