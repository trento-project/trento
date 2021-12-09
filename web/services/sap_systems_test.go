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

func sapSystemsFixtures() entities.SAPSystemInstances {
	return entities.SAPSystemInstances{
		{
			ID:             "sap_system_1",
			AgentID:        "1",
			SID:            "HA1",
			Type:           "application",
			InstanceNumber: "00",
			Features:       "features",
			DBHost:         "dbhost",
			DBName:         "tenant",
			Host: &entities.Host{
				AgentID:     "1",
				Name:        "apphost",
				ClusterID:   "cluster_id_1",
				ClusterName: "appcluster",
			},
			Tags: []*models.Tag{
				{
					Value:        "tag1",
					ResourceID:   "sap_system_1",
					ResourceType: models.TagSAPSystemResourceType,
				},
			},
		},
		{
			ID:                      "sap_system_2",
			AgentID:                 "2",
			SID:                     "PRD",
			Type:                    "database",
			InstanceNumber:          "10",
			Features:                "features",
			Tenants:                 pq.StringArray{"tenant"},
			SystemReplication:       "Primary",
			SystemReplicationStatus: "SOK",
			Host: &entities.Host{
				AgentID:     "2",
				Name:        "dbhost",
				ClusterID:   "cluster_id_2",
				ClusterName: "dbcluster",
			},
			Tags: []*models.Tag{
				{
					Value:        "tag2",
					ResourceID:   "sap_system_2",
					ResourceType: models.TagDatabaseResourceType,
				},
			},
		},
	}
}

type SAPSystemsServiceTestSuite struct {
	suite.Suite
	db                *gorm.DB
	tx                *gorm.DB
	sapSystemsService *sapSystemsService
}

func TestSAPSystemsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SAPSystemsServiceTestSuite))
}

func (suite *SAPSystemsServiceTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())

	suite.db.AutoMigrate(&entities.SAPSystemInstance{}, &entities.Host{}, &models.Tag{})
	sapSystemInstances := sapSystemsFixtures()
	err := suite.db.Create(&sapSystemInstances).Error
	suite.NoError(err)
}

func (suite *SAPSystemsServiceTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(&entities.SAPSystemInstance{},
		&entities.Host{},
		&models.Tag{})
}

func (suite *SAPSystemsServiceTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
	suite.sapSystemsService = NewSAPSystemsService(suite.tx)
}

func (suite *SAPSystemsServiceTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func (suite *SAPSystemsServiceTestSuite) TestSAPSystemsService_GetAllApplications() {
	applications, err := suite.sapSystemsService.GetAllApplications(nil, nil)
	suite.NoError(err)
	suite.Equal(1, len(applications))

	suite.EqualValues(models.SAPSystemList{
		{
			ID:     "sap_system_1",
			SID:    "HA1",
			Type:   models.SAPSystemTypeApplication,
			DBHost: "dbhost",
			DBName: "tenant",
			Instances: []*models.SAPSystemInstance{
				{
					Features:       "features",
					InstanceNumber: "00",
					HostID:         "1",
					Hostname:       "apphost",
					ClusterID:      "cluster_id_1",
					ClusterName:    "appcluster",
					SID:            "HA1",
					Type:           models.SAPSystemTypeApplication,
				},
			},
			AttachedDatabase: &models.SAPSystem{
				ID:   "sap_system_2",
				SID:  "PRD",
				Type: models.SAPSystemTypeDatabase,
				Instances: []*models.SAPSystemInstance{
					{
						Features:                "features",
						InstanceNumber:          "10",
						SystemReplication:       "Primary",
						SystemReplicationStatus: "SOK",
						SID:                     "PRD",
						Type:                    models.SAPSystemTypeDatabase,
					},
				},
			},
			Tags: []string{"tag1"},
		},
	}, applications)
}

func (suite *SAPSystemsServiceTestSuite) TestSAPSystemsService_GetAllApplications_Filter() {
	applications, err := suite.sapSystemsService.GetAllApplications(&SAPSystemFilter{
		SIDs: []string{"HA1"}, Tags: []string{"tag1"},
	}, nil)
	suite.NoError(err)
	suite.Equal(1, len(applications))
	suite.Equal("HA1", applications[0].SID)
}

func (suite *SAPSystemsServiceTestSuite) TestSAPSystemsService_GetApplicationsCount() {
	count, err := suite.sapSystemsService.GetApplicationsCount()
	suite.NoError(err)
	suite.Equal(1, count)
}

func (suite *SAPSystemsServiceTestSuite) TestSAPSystemsService_GetAllApplicationsTags() {
	tags, err := suite.sapSystemsService.GetAllApplicationsTags()
	suite.NoError(err)
	suite.Equal([]string{"tag1"}, tags)
}
func (suite *SAPSystemsServiceTestSuite) TestSAPSystemsService_GetAllApplicationsSIDs() {
	sids, err := suite.sapSystemsService.GetAllApplicationsSIDs()
	suite.NoError(err)
	suite.Equal([]string{"HA1"}, sids)
}
func (suite *SAPSystemsServiceTestSuite) TestSAPSystemsService_GetAllDatabases() {
	databases, err := suite.sapSystemsService.GetAllDatabases(nil, nil)
	suite.NoError(err)
	suite.Equal(1, len(databases))

	suite.EqualValues(models.SAPSystemList{
		{
			ID:   "sap_system_2",
			SID:  "PRD",
			Type: models.SAPSystemTypeDatabase,
			Instances: []*models.SAPSystemInstance{
				{
					Features:                "features",
					InstanceNumber:          "10",
					HostID:                  "2",
					Hostname:                "dbhost",
					ClusterID:               "cluster_id_2",
					ClusterName:             "dbcluster",
					SystemReplication:       "Primary",
					SystemReplicationStatus: "SOK",
					SID:                     "PRD",
					Type:                    models.SAPSystemTypeDatabase,
				},
			},
			Tags: []string{"tag2"},
		},
	}, databases)
}

func (suite *SAPSystemsServiceTestSuite) TestSAPSystemsService_GetAllDatabases_Filter() {
	databases, err := suite.sapSystemsService.GetAllDatabases(&SAPSystemFilter{
		SIDs: []string{"PRD"}, Tags: []string{"tag2"},
	}, nil)
	suite.NoError(err)
	suite.Equal(1, len(databases))
	suite.Equal("PRD", databases[0].SID)
}

func (suite *SAPSystemsServiceTestSuite) TestSAPSystemsService_GetByID() {
	sapSystem, err := suite.sapSystemsService.GetByID("sap_system_1")
	suite.NoError(err)

	suite.Equal("sap_system_1", sapSystem.ID)
	suite.Equal("HA1", sapSystem.SID)
}

func (suite *SAPSystemsServiceTestSuite) TestSAPSystemsService_GetByID_NotFound() {
	sapSystem, err := suite.sapSystemsService.GetByID("not_found")
	suite.NoError(err)
	suite.Nil(sapSystem)
}

func (suite *SAPSystemsServiceTestSuite) TestSAPSystemsService_GetDatabasesCount() {
	count, err := suite.sapSystemsService.GetDatabasesCount()
	suite.NoError(err)
	suite.Equal(1, count)
}

func (suite *SAPSystemsServiceTestSuite) TestSAPSystemsService_GetAllDatabasesTags() {
	tags, err := suite.sapSystemsService.GetAllDatabasesTags()
	suite.NoError(err)
	suite.Equal([]string{"tag2"}, tags)
}
func (suite *SAPSystemsServiceTestSuite) TestSAPSystemsService_GetAllDatabasesSIDs() {
	sids, err := suite.sapSystemsService.GetAllDatabasesSIDs()
	suite.NoError(err)
	suite.Equal([]string{"PRD"}, sids)
}
