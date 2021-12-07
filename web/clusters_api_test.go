package web

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

type ClustersApiTestCase struct {
	suite.Suite
	mockClusterService *services.MockClustersService
	config             *Config
	deps               Dependencies
}

func TestClustersApiTestCase(t *testing.T) {
	suite.Run(t, new(ClustersApiTestCase))
}

func (suite *ClustersApiTestCase) SetupTest() {
	suite.mockClusterService = new(services.MockClustersService)
	suite.config = setupTestConfig()
	suite.deps = setupTestDependencies()
}

func (suite *ClustersApiTestCase) Test_EmptyClustersSettings() {
	suite.mockClusterService.On("GetAllClustersSettings").Return(models.ClustersSettings{}, nil)
	suite.deps.clustersService = suite.mockClusterService

	app, err := NewAppWithDeps(suite.config, suite.deps)
	if err != nil {
		suite.T().Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/internal/clusters/settings", nil)
	app.webEngine.ServeHTTP(resp, req)

	suite.Equal(200, resp.Code)
	suite.JSONEq(`[]`, resp.Body.String())
}

func (suite *ClustersApiTestCase) Test_ClustersSettingsWereFound() {
	mockedClustersSettings := mockedClustersSettings()
	suite.mockClusterService.On("GetAllClustersSettings").Return(mockedClustersSettings, nil)

	suite.deps.clustersService = suite.mockClusterService

	app, err := NewAppWithDeps(suite.config, suite.deps)
	if err != nil {
		suite.T().Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/internal/clusters/settings", nil)
	app.webEngine.ServeHTTP(resp, req)

	expectedJson, _ := json.Marshal(mockedClustersSettings)
	suite.Equal(200, resp.Code)
	suite.JSONEq(string(expectedJson), resp.Body.String())
}

func (suite *ClustersApiTestCase) Test_AnErrorOccurs() {
	suite.mockClusterService.On("GetAllClustersSettings").Return(nil, errors.New("KABOOM"))

	suite.deps.clustersService = suite.mockClusterService

	app, err := NewAppWithDeps(suite.config, suite.deps)
	if err != nil {
		suite.T().Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/internal/clusters/settings", nil)
	app.webEngine.ServeHTTP(resp, req)

	suite.Equal(500, resp.Code)
	suite.JSONEq(`{"error":"KABOOM"}`, resp.Body.String())
}

func mockedClustersSettings() models.ClustersSettings {
	return models.ClustersSettings{
		{
			ID:             "cluster1",
			SelectedChecks: []string{"A", "B", "C"},
			Hosts: []*models.ConnectionInfoAwareHost{
				{
					Name:    "host1",
					Address: "10.0.0.1",
					User:    "root",
				},
			},
		},
		{
			ID:             "cluster2",
			SelectedChecks: []string{},
			Hosts: []*models.ConnectionInfoAwareHost{
				{
					Name:    "host2",
					Address: "10.0.0.2",
					User:    "theuser",
				},
			},
		},
	}
}
