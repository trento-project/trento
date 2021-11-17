package services

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
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

	suite.db.AutoMigrate(entities.Cluster{}, models.Tag{})
	loadClustersFixtures(suite.db)
}

func (suite *ClustersServiceTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(entities.Cluster{}, models.Tag{})
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
	db.Create(&entities.Cluster{
		ID:              "1",
		Name:            "cluster1",
		ClusterType:     models.ClusterTypeScaleUp,
		SIDs:            []string{"DEV"},
		ResourcesNumber: 10,
		HostsNumber:     2,
		Tags: []models.Tag{
			{
				ResourceID:   "1",
				ResourceType: models.TagClusterResourceType,
				Value:        "tag1",
			},
		},
	})
	db.Create(&entities.Cluster{
		ID:              "2",
		Name:            "cluster2",
		ClusterType:     models.ClusterTypeScaleOut,
		SIDs:            []string{"QAS"},
		ResourcesNumber: 11,
		HostsNumber:     2,
		Tags: []models.Tag{
			{
				ResourceID:   "2",
				ResourceType: models.TagClusterResourceType,
				Value:        "tag2",
			},
		},
	})
	db.Create(&entities.Cluster{
		ID:              "3",
		Name:            "cluster3",
		ClusterType:     models.ClusterTypeUnknown,
		SIDs:            []string{"PRD", "PRD2"},
		ResourcesNumber: 3,
		HostsNumber:     5,
		Tags: []models.Tag{
			{
				ResourceID:   "3",
				ResourceType: models.TagClusterResourceType,
				Value:        "tag3",
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
			ClusterType:     models.ClusterTypeScaleUp,
			SIDs:            []string{"DEV"},
			ResourcesNumber: 10,
			HostsNumber:     2,
			Health:          models.CheckPassing,
			Tags:            []string{"tag1"},
		},
		&models.Cluster{
			ID:              "2",
			Name:            "cluster2",
			ClusterType:     models.ClusterTypeScaleOut,
			SIDs:            []string{"QAS"},
			ResourcesNumber: 11,
			HostsNumber:     2,
			Health:          models.CheckWarning,
			Tags:            []string{"tag2"},
		},
		&models.Cluster{
			ID:              "3",
			Name:            "cluster3",
			ClusterType:     models.ClusterTypeUnknown,
			SIDs:            []string{"PRD", "PRD2"},
			ResourcesNumber: 3,
			HostsNumber:     5,
			Health:          models.CheckCritical,
			Tags:            []string{"tag3"},
		},
	}, clusters)
}

func (suite *ClustersServiceTestSuite) TestClustersService_GetAll_Filter() {
	suite.checksService.On("GetAggregatedChecksResultByCluster", "1").Return(&AggregatedCheckData{PassingCount: 1}, nil)

	clusters, _ := suite.clustersService.GetAll(&ClustersFilter{
		Name:        []string{"cluster1"},
		SIDs:        []string{"DEV"},
		ClusterType: []string{models.ClusterTypeScaleUp},
		Health:      []string{models.CheckPassing},
		Tags:        []string{"tag1"},
	}, nil)

	suite.Equal(1, len(clusters))
	suite.Equal(clusters[0].ID, "1")
}

func (suite *ClustersServiceTestSuite) TestClustersService_GetClustersCount() {
	count, err := suite.clustersService.GetCount()

	suite.NoError(err)
	suite.Equal(3, count)
}

func (suite *ClustersServiceTestSuite) TestClustersService_GetAllClusterTypes() {
	clusterTypes, _ := suite.clustersService.GetAllClusterTypes()
	suite.ElementsMatch(
		[]string{models.ClusterTypeScaleUp,
			models.ClusterTypeScaleOut,
			models.ClusterTypeUnknown}, clusterTypes)
}

func (suite *ClustersServiceTestSuite) TestClustersService_GetAllTags() {
	tags, _ := suite.clustersService.GetAllTags()
	suite.ElementsMatch([]string{"tag1", "tag2", "tag3"}, tags)
}

func (suite *ClustersServiceTestSuite) TestClustersService_GetAllSIDs() {
	sids, _ := suite.clustersService.GetAllSIDs()
	suite.ElementsMatch([]string{"DEV", "QAS", "PRD", "PRD2"}, sids)
}
