package services

import (
	"testing"

	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/web/models"
)

type HealthSummaryServiceTestSuite struct {
	suite.Suite
}

func TestHealthSummaryServiceestSuite(t *testing.T) {
	suite.Run(t, new(HealthSummaryServiceTestSuite))
}

func (suite *HealthSummaryServiceTestSuite) TestGetHealthSummary() {
	sapSystemsService := new(MockSAPSystemsService)
	clustersService := new(MockClustersService)
	hostsService := new(MockHostsService)

	sapSystemsService.On("GetAllApplications", mock.Anything, mock.Anything).Return(models.SAPSystemList{
		{
			ID:     "application_id",
			SID:    "HA1",
			Type:   models.SAPSystemTypeApplication,
			Health: models.SAPSystemHealthPassing,
			Instances: []*models.SAPSystemInstance{
				{
					HostID:    "netweaver01",
					ClusterID: "netweaver_cluster",
				},
				{
					HostID:    "netweaver02",
					ClusterID: "netweaver_cluster",
				},
			},
			AttachedDatabase: &models.SAPSystem{
				ID:     "database_id",
				SID:    "PRD",
				Health: models.SAPSystemHealthPassing,
				Instances: []*models.SAPSystemInstance{
					{
						HostID:    "hana01",
						ClusterID: "hana_cluster",
					},
				},
			},
		},
	}, nil)

	clustersService.On("GetAll", mock.Anything, mock.Anything).Return(models.ClusterList{
		{
			ID:     "hana_cluster",
			Health: models.CheckCritical,
		},
	}, nil)

	hostsService.On("GetAll", mock.Anything, mock.Anything).Return(models.HostList{
		{
			ID:     "host_id",
			Health: models.HostHealthWarning,
		},
		{
			ID:     "netweaver01",
			Health: models.HostHealthPassing,
		}}, nil)

	healthSummaryService := NewHealthSummaryService(sapSystemsService, clustersService, hostsService)
	healthSummary, _ := healthSummaryService.GetHealthSummary()

	suite.EqualValues(models.HealthSummary{{
		ID: "application_id", SID: "HA1",
		SAPSystemHealth: models.HealthSummaryHealthPassing,
		ClustersHealth:  models.HealthSummaryHealthCritical,
		DatabaseHealth:  models.HealthSummaryHealthPassing,
		HostsHealth:     models.HealthSummaryHealthWarning,
	}}, healthSummary)
}
