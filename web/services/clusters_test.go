package services

import (
	"encoding/json"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	araMocks "github.com/trento-project/trento/web/services/ara/mocks"
	"gorm.io/gorm"
)

type ClustersServiceTestSuite struct {
	suite.Suite
	db              *gorm.DB
	tx              *gorm.DB
	clustersService *clustersService
	checksService   *MockChecksService
}

func TestClustersServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ClustersServiceTestSuite))
}

func (suite *ClustersServiceTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())

	suite.db.AutoMigrate(entities.Cluster{}, entities.Host{}, models.Tag{}, models.SelectedChecks{}, models.ConnectionSettings{})
	loadClustersFixtures(suite.db)
}

func (suite *ClustersServiceTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(entities.Cluster{}, entities.Host{}, models.Tag{}, models.SelectedChecks{}, models.ConnectionSettings{})
}

func (suite *ClustersServiceTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
	suite.checksService = new(MockChecksService)
	suite.clustersService = NewClustersService(suite.tx, suite.checksService)
}

func (suite *ClustersServiceTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func loadClustersFixtures(db *gorm.DB) {
	details := &entities.HANAClusterDetails{
		Nodes: []*entities.HANAClusterNode{
			{
				Name: "host1",
			},
		},
	}

	detailsJSON, _ := json.Marshal(details)

	err := db.Create(&entities.Cluster{
		ID:              "1",
		Name:            "cluster1",
		ClusterType:     models.ClusterTypeHANAScaleUp,
		SID:             "DEV",
		ResourcesNumber: 10,
		HostsNumber:     2,
		Tags: []*models.Tag{
			{
				ResourceID:   "1",
				ResourceType: models.TagClusterResourceType,
				Value:        "tag1",
			},
		},
		Hosts: []*entities.Host{
			{
				AgentID:     "1",
				SSHAddress:  "10.74.2.10",
				ClusterID:   "1",
				Name:        "host1",
				IPAddresses: []string{"10.74.1.10"},
			},
		},
		Details: detailsJSON,
	}).Error
	if err != nil {
		panic(err)
	}
	db.Create(&entities.Cluster{
		ID:              "2",
		Name:            "cluster2",
		ClusterType:     models.ClusterTypeHANAScaleOut,
		SID:             "QAS",
		ResourcesNumber: 11,
		HostsNumber:     2,
		Tags: []*models.Tag{
			{
				ResourceID:   "2",
				ResourceType: models.TagClusterResourceType,
				Value:        "tag2",
			},
		},
		Hosts: []*entities.Host{
			{
				AgentID:     "2",
				SSHAddress:  "10.74.2.11",
				ClusterID:   "2",
				Name:        "host2",
				IPAddresses: pq.StringArray{"10.74.1.11"},
			},
		},
	})

	azureCloudData := &entities.AzureCloudData{}
	azureCloudData.AdminUsername = "cloudadmin"
	cloudData, _ := json.Marshal(&azureCloudData)

	db.Create(&entities.Cluster{
		ID:              "3",
		Name:            "cluster3",
		ClusterType:     models.ClusterTypeUnknown,
		SID:             "PRD",
		ResourcesNumber: 3,
		HostsNumber:     5,
		Tags: []*models.Tag{
			{
				ResourceID:   "3",
				ResourceType: models.TagClusterResourceType,
				Value:        "tag3",
			},
		},
		Hosts: []*entities.Host{
			{
				AgentID:       "3",
				SSHAddress:    "10.74.2.12",
				ClusterID:     "3",
				Name:          "host3",
				IPAddresses:   pq.StringArray{"10.74.1.12"},
				CloudProvider: "azure",
				CloudData:     cloudData,
			},
		},
	})
}

func (suite *ClustersServiceTestSuite) TestClustersService_GetAll() {
	suite.checksService.On("GetAggregatedChecksResultByCluster", "1").Return(&AggregatedCheckData{PassingCount: 1}, nil)
	suite.checksService.On("GetAggregatedChecksResultByCluster", "2").Return(&AggregatedCheckData{WarningCount: 1}, nil)
	suite.checksService.On("GetAggregatedChecksResultByCluster", "3").Return(&AggregatedCheckData{CriticalCount: 1}, nil)

	clusters, _ := suite.clustersService.GetAll(nil, nil)

	suite.ElementsMatch(models.ClusterList{
		&models.Cluster{
			ID:              "1",
			Name:            "cluster1",
			ClusterType:     models.ClusterTypeHANAScaleUp,
			SID:             "DEV",
			ResourcesNumber: 10,
			HostsNumber:     2,
			Health:          models.CheckPassing,
			PassingCount:    1,
			WarningCount:    0,
			CriticalCount:   0,
			Tags:            []string{"tag1"},
		},
		&models.Cluster{
			ID:              "2",
			Name:            "cluster2",
			ClusterType:     models.ClusterTypeHANAScaleOut,
			SID:             "QAS",
			ResourcesNumber: 11,
			HostsNumber:     2,
			Health:          models.CheckWarning,
			PassingCount:    0,
			WarningCount:    1,
			CriticalCount:   0,
			Tags:            []string{"tag2"},
		},
		&models.Cluster{
			ID:              "3",
			Name:            "cluster3",
			ClusterType:     models.ClusterTypeUnknown,
			SID:             "PRD",
			ResourcesNumber: 3,
			HostsNumber:     5,
			Health:          models.CheckCritical,
			PassingCount:    0,
			WarningCount:    0,
			CriticalCount:   1,
			Tags:            []string{"tag3"},
		},
	}, clusters)
}

func (suite *ClustersServiceTestSuite) TestClustersService_GetAll_Filter() {
	suite.checksService.On("GetAggregatedChecksResultByCluster", "1").Return(&AggregatedCheckData{PassingCount: 1}, nil)

	clusters, _ := suite.clustersService.GetAll(&ClustersFilter{
		Name:        []string{"cluster1"},
		SIDs:        []string{"DEV"},
		ClusterType: []string{models.ClusterTypeHANAScaleUp},
		Health:      []string{models.CheckPassing},
		Tags:        []string{"tag1"},
	}, nil)

	suite.Equal(1, len(clusters))
	suite.Equal(clusters[0].ID, "1")
}
func (suite *ClustersServiceTestSuite) TestClustersService_GetByID() {
	suite.checksService.On("GetAggregatedChecksResultByCluster", "1").Return(&AggregatedCheckData{PassingCount: 1}, nil)
	suite.checksService.On("GetAggregatedChecksResultByHost", "1").Return(map[string]*AggregatedCheckData{
		"host1": {PassingCount: 1},
	}, nil)

	cluster, err := suite.clustersService.GetByID("1")

	suite.NoError(err)
	suite.Equal("cluster1", cluster.Name)
	suite.EqualValues(&models.HANAClusterDetails{
		Nodes: []*models.HANAClusterNode{
			{
				HostID:      "1",
				Name:        "host1",
				Health:      models.CheckPassing,
				IPAddresses: []string{"10.74.1.10"},
			},
		},
	}, cluster.Details.(*models.HANAClusterDetails))
}
func (suite *ClustersServiceTestSuite) TestClustersService_GetByID_NotFound() {
	cluster, err := suite.clustersService.GetByID("not_there")

	suite.NoError(err)
	suite.Nil(cluster)
}
func (suite *ClustersServiceTestSuite) TestClustersService_GetClustersCount() {
	count, err := suite.clustersService.GetCount()

	suite.NoError(err)
	suite.Equal(3, count)
}

func (suite *ClustersServiceTestSuite) TestClustersService_GetAllClusterTypes() {
	clusterTypes, _ := suite.clustersService.GetAllClusterTypes()
	suite.ElementsMatch(
		[]string{models.ClusterTypeHANAScaleUp,
			models.ClusterTypeHANAScaleOut,
			models.ClusterTypeUnknown}, clusterTypes)
}

func (suite *ClustersServiceTestSuite) TestClustersService_GetAllTags() {
	tags, _ := suite.clustersService.GetAllTags()
	suite.ElementsMatch([]string{"tag1", "tag2", "tag3"}, tags)
}

func (suite *ClustersServiceTestSuite) TestClustersService_GetAllSIDs() {
	sids, _ := suite.clustersService.GetAllSIDs()
	suite.ElementsMatch([]string{"DEV", "QAS", "PRD"}, sids)
}

func (suite *ClustersServiceTestSuite) TestClustersService_GetAllClustersSettingsReturnsNoSettings() {
	mockAra := new(araMocks.AraService)

	tx := suite.tx.Raw("TRUNCATE TABLE clusters")
	checksService := NewChecksService(mockAra, tx)
	suite.clustersService = NewClustersService(tx, checksService)

	clustersSettings, err := suite.clustersService.GetAllClustersSettings()
	suite.NoError(err)
	suite.Empty(clustersSettings)

	tx.Rollback()
}

func (suite *ClustersServiceTestSuite) TestClustersService_GetAllClustersSettingsReturnsExpectedSettings() {
	mockAra := new(araMocks.AraService)
	checksService := NewChecksService(mockAra, suite.db)
	suite.clustersService = NewClustersService(suite.db, checksService)

	loadClustersSettingsFixtures(suite.db)

	clustersSettings, err := suite.clustersService.GetAllClustersSettings()
	suite.NoError(err)
	suite.NotEmpty(clustersSettings)
	suite.Len(clustersSettings, 3)

	suite.EqualValues(expectedClustersSettingsFixtures(), clustersSettings)
}

func loadClustersSettingsFixtures(db *gorm.DB) {
	db.Create(&models.SelectedChecks{
		ID:             "1",
		SelectedChecks: []string{"A", "B", "C"},
	})
	db.Create(&models.SelectedChecks{
		ID:             "2",
		SelectedChecks: []string{},
	})
	db.Create(&models.ConnectionSettings{
		ID:   "1",
		Node: "host1",
		User: "theuser",
	})
	db.Create(&models.ConnectionSettings{
		ID:   "2",
		Node: "host2",
		User: "root",
	})
	db.Create(&models.ConnectionSettings{
		ID:   "3",
		Node: "host3",
		User: "",
	})
}

func expectedClustersSettingsFixtures() models.ClustersSettings {
	return models.ClustersSettings{
		{
			ID:             "1",
			SelectedChecks: []string{"A", "B", "C"},
			Hosts: []*models.HostConnection{
				{
					Name:    "host1",
					Address: "10.74.2.10",
					User:    "theuser",
				},
			},
		},
		{
			ID:             "2",
			SelectedChecks: []string{},
			Hosts: []*models.HostConnection{
				{
					Name:    "host2",
					Address: "10.74.2.11",
					User:    "root",
				},
			},
		},
		{
			ID:             "3",
			SelectedChecks: []string{},
			Hosts: []*models.HostConnection{
				{
					Name:    "host3",
					Address: "10.74.2.12",
					User:    "cloudadmin",
				},
			},
		},
	}
}
