package services

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
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
	prometheusService *prometheusService
}

func TestPrometheusServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PrometheusServiceTestSuite))
}

func (suite *PrometheusServiceTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())

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
	suite.prometheusService = NewPrometheusService(suite.tx)
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
			Labels:  map[string]string{"agentID": "1", "hostname": "host1"},
		},
		&models.PrometheusTargets{
			Targets: []string{"192.168.1.2:9100"},
			Labels:  map[string]string{"agentID": "2", "hostname": "host2"},
		},
		&models.PrometheusTargets{
			Targets: []string{"192.168.1.3:9100"},
			Labels:  map[string]string{"agentID": "3", "hostname": "host3"},
		},
	}, targets)
}
