package services

import (
	"testing"
	"time"

	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"

	prometheusV1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prometheusModel "github.com/prometheus/common/model"
	prometheusInternal "github.com/trento-project/trento/internal/prometheus"
)

func targetsFixtures() []entities.Host {
	return []entities.Host{
		{
			AgentID:    "1",
			Name:       "host1",
			SSHAddress: "192.168.1.1",
		},
		{
			AgentID:    "2",
			Name:       "host2",
			SSHAddress: "192.168.1.2",
		},
		{
			AgentID:    "3",
			Name:       "host3",
			SSHAddress: "192.168.1.3",
		},
	}
}

type PrometheusServiceTestSuite struct {
	suite.Suite
	db                *gorm.DB
	tx                *gorm.DB
	prometheusApi     *prometheusInternal.MockPrometheusAPI
	prometheusService *prometheusService
}

func TestPrometheusServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PrometheusServiceTestSuite))
}

func (suite *PrometheusServiceTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())
	suite.prometheusApi = new(prometheusInternal.MockPrometheusAPI)

	suite.db.AutoMigrate(&entities.Host{}, &entities.HostHeartbeat{}, &entities.SAPSystemInstance{}, &models.Tag{})
	hosts := targetsFixtures()
	err := suite.db.Create(&hosts).Error
	suite.NoError(err)
}

func (suite *PrometheusServiceTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(&entities.Host{})
}

func (suite *PrometheusServiceTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
	suite.prometheusService = NewPrometheusService(suite.tx, suite.prometheusApi)
}

func (suite *PrometheusServiceTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func (suite *PrometheusServiceTestSuite) TestPrometheusService_GetHttpSDTargets() {
	targets, err := suite.prometheusService.GetHttpSDTargets()
	suite.NoError(err)

	suite.ElementsMatch(models.PrometheusTargetsList{
		&models.PrometheusTargets{
			Targets: []string{"192.168.1.1:9100"},
			Labels:  map[string]string{"agentID": "1", "hostname": "host1", "exporter_name": "Node Exporter"},
		},
		&models.PrometheusTargets{
			Targets: []string{"192.168.1.2:9100"},
			Labels:  map[string]string{"agentID": "2", "hostname": "host2", "exporter_name": "Node Exporter"},
		},
		&models.PrometheusTargets{
			Targets: []string{"192.168.1.3:9100"},
			Labels:  map[string]string{"agentID": "3", "hostname": "host3", "exporter_name": "Node Exporter"},
		},
	}, targets)
}

func (suite *PrometheusServiceTestSuite) TestPrometheusService_Query() {
	cTime := time.Now()
	expectedResult := prometheusModel.Vector{
		&prometheusModel.Sample{
			Metric: prometheusModel.Metric{
				"exporter_name": "some exporter",
				"job":           "some job",
			},
			Value:     1,
			Timestamp: 1234567,
		},
	}
	suite.prometheusApi.On("Query", mock.Anything, "some nice query", cTime).Return(
		expectedResult, prometheusV1.Warnings{}, nil,
	)

	result, err := suite.prometheusService.Query("some nice query", cTime)
	suite.NoError(err)

	suite.Equal(expectedResult, result)
}
