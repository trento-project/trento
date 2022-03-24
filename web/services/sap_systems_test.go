package services

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/internal/sapsystem/sapcontrol"
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
			Status:         string(sapcontrol.STATECOLOR_RED),
			DBHost:         "dbhost_1",
			DBName:         "tenant",
			DBAddress:      "192.168.1.10",
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
			Status:                  string(sapcontrol.STATECOLOR_GREEN),
			Tenants:                 pq.StringArray{"tenant"},
			SystemReplication:       "Primary",
			SystemReplicationStatus: "SOK",
			Host: &entities.Host{
				AgentID:     "2",
				Name:        "dbhost_1",
				ClusterID:   "cluster_id_2",
				ClusterName: "dbcluster",
				IPAddresses: pq.StringArray{"192.168.1.10"},
			},
			Tags: []*models.Tag{
				{
					Value:        "tag2",
					ResourceID:   "sap_system_2",
					ResourceType: models.TagDatabaseResourceType,
				},
			},
		},
		{
			ID:                      "sap_system_2",
			AgentID:                 "3",
			SID:                     "PRD",
			Type:                    "database",
			InstanceNumber:          "11",
			Features:                "features",
			Status:                  string(sapcontrol.STATECOLOR_YELLOW),
			Tenants:                 pq.StringArray{"tenant"},
			SystemReplication:       "Secondary",
			SystemReplicationStatus: "SOK",
			Host: &entities.Host{
				AgentID:     "3",
				Name:        "dbhost_2",
				ClusterID:   "cluster_id_2",
				ClusterName: "dbcluster",
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
			ID:        "sap_system_1",
			SID:       "HA1",
			Type:      models.SAPSystemTypeApplication,
			DBHost:    "dbhost_1",
			DBName:    "tenant",
			DBAddress: "192.168.1.10",
			Health:    models.SAPSystemHealthCritical,
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
					Status:         string(sapcontrol.STATECOLOR_RED),
				},
			},
			AttachedDatabase: &models.SAPSystem{
				ID:     "sap_system_2",
				SID:    "PRD",
				Type:   models.SAPSystemTypeDatabase,
				Health: models.SAPSystemHealthWarning,
				Instances: []*models.SAPSystemInstance{
					{
						HostID:                  "2",
						Hostname:                "dbhost_1",
						ClusterID:               "cluster_id_2",
						ClusterName:             "dbcluster",
						Features:                "features",
						InstanceNumber:          "10",
						SystemReplication:       "Primary",
						SystemReplicationStatus: "SOK",
						SID:                     "PRD",
						Type:                    models.SAPSystemTypeDatabase,
						Status:                  string(sapcontrol.STATECOLOR_GREEN),
					},
					{
						HostID:                  "3",
						Hostname:                "dbhost_2",
						ClusterID:               "cluster_id_2",
						ClusterName:             "dbcluster",
						Features:                "features",
						InstanceNumber:          "11",
						SystemReplication:       "Secondary",
						SystemReplicationStatus: "SOK",
						SID:                     "PRD",
						Type:                    models.SAPSystemTypeDatabase,
						Status:                  string(sapcontrol.STATECOLOR_YELLOW),
					},
				},
			},
			Tags: []string{"tag1"},
		},
	}, applications)
}

func (suite *SAPSystemsServiceTestSuite) TestSAPSystemsService_getAllByType_Pagination() {
	applications, err := suite.sapSystemsService.getAllByType(
		models.SAPSystemTypeDatabase, models.TagSAPSystemResourceType,
		&SAPSystemFilter{
			SIDs: []string{"PRD"}, Tags: []string{},
		}, &Page{Number: 1, Size: 1})
	suite.NoError(err)
	suite.Equal(1, len(applications))
	suite.Equal(2, len(applications[0].Instances))
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
			ID:     "sap_system_2",
			SID:    "PRD",
			Type:   models.SAPSystemTypeDatabase,
			Health: models.SAPSystemHealthWarning,
			Instances: []*models.SAPSystemInstance{
				{
					Features:                "features",
					InstanceNumber:          "10",
					HostID:                  "2",
					Hostname:                "dbhost_1",
					ClusterID:               "cluster_id_2",
					ClusterName:             "dbcluster",
					SystemReplication:       "Primary",
					SystemReplicationStatus: "SOK",
					SID:                     "PRD",
					Type:                    models.SAPSystemTypeDatabase,
					Status:                  string(sapcontrol.STATECOLOR_GREEN),
				},
				{
					Features:                "features",
					InstanceNumber:          "11",
					HostID:                  "3",
					Hostname:                "dbhost_2",
					ClusterID:               "cluster_id_2",
					ClusterName:             "dbcluster",
					SystemReplication:       "Secondary",
					SystemReplicationStatus: "SOK",
					SID:                     "PRD",
					Type:                    models.SAPSystemTypeDatabase,
					Status:                  string(sapcontrol.STATECOLOR_YELLOW),
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
