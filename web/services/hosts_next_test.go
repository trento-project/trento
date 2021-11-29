package services

import (
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
)

func hostsFixtures() []entities.Host {
	return []entities.Host{
		{
			AgentID:       "1",
			Name:          "host1",
			ClusterID:     "cluster_id_1",
			ClusterName:   "cluster_1",
			CloudProvider: "azure",
			IPAddresses:   pq.StringArray{"10.74.1.5"},
			SAPSystemInstances: []*entities.SAPSystemInstance{
				{
					AgentID:        "1",
					ID:             "sap_system_id_1",
					SID:            "DEV",
					InstanceNumber: "00",
				},
			},
			AgentVersion: "rolling1337",
			Heartbeat: &entities.HostHeartbeat{
				AgentID:   "1",
				UpdatedAt: time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC),
			},
			Tags: []*models.Tag{{
				Value:        "tag1",
				ResourceID:   "1",
				ResourceType: models.TagHostResourceType,
			}},
		},
		{
			AgentID:       "2",
			Name:          "host2",
			ClusterID:     "cluster_id_2",
			ClusterName:   "cluster_2",
			CloudProvider: "azure",
			IPAddresses:   pq.StringArray{"10.74.1.10"},
			SAPSystemInstances: []*entities.SAPSystemInstance{
				{
					AgentID:        "2",
					ID:             "sap_system_id_2",
					SID:            "QAS",
					InstanceNumber: "10",
				},
			},
			AgentVersion: "stable",
			Heartbeat: &entities.HostHeartbeat{
				AgentID:   "2",
				UpdatedAt: time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC),
			},
			Tags: []*models.Tag{{
				Value:        "tag2",
				ResourceID:   "2",
				ResourceType: models.TagHostResourceType,
			}},
		},
	}
}

type HostsNextServiceTestSuite struct {
	suite.Suite
	db               *gorm.DB
	tx               *gorm.DB
	hostsNextService *hostsNextService
}

func TestHostsNextServiceTestSuite(t *testing.T) {
	suite.Run(t, new(HostsNextServiceTestSuite))
}

func (suite *HostsNextServiceTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())

	suite.db.AutoMigrate(&entities.Host{}, &entities.HostHeartbeat{}, &entities.SAPSystemInstance{}, &models.Tag{})
	hosts := hostsFixtures()
	err := suite.db.Create(&hosts).Error
	suite.NoError(err)
}

func (suite *HostsNextServiceTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(&entities.Host{},
		&entities.HostHeartbeat{},
		&entities.SAPSystemInstance{},
		&models.Tag{})
}

func (suite *HostsNextServiceTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
	suite.hostsNextService = NewHostsNextService(suite.tx)
}

func (suite *HostsNextServiceTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func (suite *HostsNextServiceTestSuite) TestHostsNextService_GetAll() {
	timeSince = func(_ time.Time) time.Duration {
		return time.Duration(0)
	}

	hosts, err := suite.hostsNextService.GetAll(nil, nil)
	suite.NoError(err)

	suite.ElementsMatch(models.HostList{
		{
			ID:            "1",
			Name:          "host1",
			Health:        "passing",
			IPAddresses:   []string{"10.74.1.5"},
			CloudProvider: "azure",
			ClusterID:     "cluster_id_1",
			ClusterName:   "cluster_1",
			AgentVersion:  "rolling1337",
			SAPSystems: []*models.SAPSystem{
				{
					ID:  "sap_system_id_1",
					SID: "DEV",
					Instances: []*models.SAPSystemInstance{
						{
							InstanceNumber: "00",
						},
					},
				},
			},
			Tags: []string{"tag1"},
		},
		{
			ID:            "2",
			Name:          "host2",
			Health:        "passing",
			IPAddresses:   []string{"10.74.1.10"},
			CloudProvider: "azure",
			ClusterID:     "cluster_id_2",
			ClusterName:   "cluster_2",
			AgentVersion:  "stable",
			SAPSystems: []*models.SAPSystem{
				{
					ID:  "sap_system_id_2",
					SID: "QAS",
					Instances: []*models.SAPSystemInstance{
						{
							InstanceNumber: "10",
						},
					},
				},
			},
			Tags: []string{"tag2"},
		},
	}, hosts)
}

func (suite *HostsNextServiceTestSuite) TestHostsNextService_GetAll_Filters() {
	timeSince = func(_ time.Time) time.Duration {
		return time.Duration(0)
	}

	hosts, _ := suite.hostsNextService.GetAll(&HostsFilter{
		Tags:   []string{"tag1"},
		SIDs:   []string{"DEV"},
		Health: []string{"passing", "unknown"},
	}, nil)
	suite.Equal(1, len(hosts))
	suite.Equal("1", hosts[0].ID)
}

func (suite *HostsNextServiceTestSuite) TestHostsNextService_GetHostsCount() {
	count, err := suite.hostsNextService.GetCount()

	suite.NoError(err)
	suite.Equal(2, count)
}

func (suite *HostsNextServiceTestSuite) TestHostsNextService_GetAllTags() {
	hosts, _ := suite.hostsNextService.GetAllTags()
	suite.EqualValues([]string{"tag1", "tag2"}, hosts)
}

func (suite *HostsNextServiceTestSuite) TestHostsNextService_GetAllSIDs() {
	hosts, _ := suite.hostsNextService.GetAllSIDs()
	suite.ElementsMatch([]string{"DEV", "QAS"}, hosts)
}

func (suite *HostsNextServiceTestSuite) TestHostsNextService_Heartbeat() {
	err := suite.hostsNextService.Heartbeat("1")
	suite.NoError(err)

	var heartbeat entities.HostHeartbeat
	suite.tx.First(&heartbeat)
	suite.Equal("1", heartbeat.AgentID)
}

func (suite *HostsNextServiceTestSuite) TestHostsNextService_computeHealth() {
	host := hostsFixtures()[0]

	timeSince = func(_ time.Time) time.Duration {
		return time.Duration(0)
	}
	suite.Equal(models.HostHealthPassing, computeHealth(&host))

	timeSince = func(_ time.Time) time.Duration {
		return time.Duration(HeartbeatTreshold + 1)
	}
	suite.Equal(models.HostHealthCritical, computeHealth(&host))

	host.Heartbeat = nil
	suite.Equal(models.HostHealthUnknown, computeHealth(&host))

}
