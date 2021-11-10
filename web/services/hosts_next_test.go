package services

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
)

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

	suite.db.AutoMigrate(&entities.Host{}, &models.Tag{})
	loadHostsFixtures(suite.db)
}

func (suite *HostsNextServiceTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(&entities.Host{}, &models.Tag{})
}

func (suite *HostsNextServiceTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
	suite.hostsNextService = NewHostsNextService(suite.tx)
}

func (suite *HostsNextServiceTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func (suite *HostsNextServiceTestSuite) TestHostsNextService_GetAll() {
	hosts, _ := suite.hostsNextService.GetAll(map[string][]string{})

	suite.ElementsMatch(models.HostList{
		{
			ID:            "1",
			Name:          "host1",
			Health:        "",
			IPAddresses:   []string{"10.74.1.5"},
			CloudProvider: "azure",
			ClusterID:     "cluster_id_1",
			ClusterName:   "cluster_1",
			SIDs:          []string{"DEV"},
			AgentVersion:  "rolling1337",
			Tags:          []string{"tag1"},
		},
		{
			ID:            "2",
			Name:          "host2",
			Health:        "",
			IPAddresses:   []string{"10.74.1.10"},
			CloudProvider: "azure",
			ClusterID:     "cluster_id_2",
			ClusterName:   "cluster_2",
			SIDs:          []string{"QAS"},
			AgentVersion:  "stable",
			Tags:          []string{"tag2"},
		},
	}, hosts)
}

func (suite *HostsNextServiceTestSuite) TestHostsNextService_GetAll_Filters() {
	hosts, _ := suite.hostsNextService.GetAll(map[string][]string{
		"tags": {"tag1"},
		"sids": {"DEV"},
	})
	suite.Equal(1, len(hosts))
	suite.Equal("1", hosts[0].ID)
}

func (suite *HostsNextServiceTestSuite) TestHostsNextService_GetAllTags() {
	hosts, _ := suite.hostsNextService.GetAllTags()
	suite.EqualValues([]string{"tag1", "tag2"}, hosts)
}

func (suite *HostsNextServiceTestSuite) TestHostsNextService_GetAllSIDs() {
	hosts, _ := suite.hostsNextService.GetAllSIDs()
	suite.ElementsMatch([]string{"DEV", "QAS"}, hosts)
}

func loadHostsFixtures(db *gorm.DB) {
	db.Create(&[]entities.Host{
		{
			AgentID:       "1",
			Name:          "host1",
			ClusterID:     "cluster_id_1",
			ClusterName:   "cluster_1",
			CloudProvider: "azure",
			IPAddresses:   pq.StringArray{"10.74.1.5"},
			SIDs:          pq.StringArray{"DEV"},
			AgentVersion:  "rolling1337",
			Tags: []models.Tag{{
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
			SIDs:          pq.StringArray{"QAS"},
			AgentVersion:  "stable",
			Tags: []models.Tag{{
				Value:        "tag2",
				ResourceID:   "2",
				ResourceType: models.TagHostResourceType,
			}},
		},
	})
}
