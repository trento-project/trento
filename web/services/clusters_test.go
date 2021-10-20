package services

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
)

type ClustersServiceTestSuite struct {
	suite.Suite
	db              *gorm.DB
	tx              *gorm.DB
	clustersService *clustersService
	tagsService     *MockTagsService
	checksService   *MockChecksService
}

func TestClustersServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ClustersServiceTestSuite))
}

func (suite *ClustersServiceTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase()

	suite.db.AutoMigrate(models.Cluster{})
	loadClustersFixtures(suite.db)
}

func (suite *ClustersServiceTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(models.Cluster{})
}

func (suite *ClustersServiceTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
	suite.checksService = new(MockChecksService)
	suite.tagsService = new(MockTagsService)
	suite.clustersService = NewClustersService(suite.tx, suite.checksService, suite.tagsService)
}

func (suite *ClustersServiceTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func loadClustersFixtures(db *gorm.DB) {
	db.Create(&models.Cluster{
		ID:              "1",
		Name:            "cluster1",
		ClusterType:     models.ClusterTypeScaleUp,
		SIDs:            []string{"DEV"},
		ResourcesNumber: 10,
		HostsNumber:     2,
	})
	db.Create(&models.Cluster{
		ID:              "2",
		Name:            "cluster2",
		ClusterType:     models.ClusterTypeScaleOut,
		SIDs:            []string{"QAS"},
		ResourcesNumber: 11,
		HostsNumber:     2,
	})
	db.Create(&models.Cluster{
		ID:              "3",
		Name:            "cluster3",
		ClusterType:     models.ClusterTypeUnknown,
		SIDs:            []string{"PRD", "PRD2"},
		ResourcesNumber: 3,
		HostsNumber:     5,
	})
}

func (suite *ClustersServiceTestSuite) TestClustersService_GetAll() {
	suite.checksService.On("GetAggregatedChecksResultByCluster", "1").Return(&AggregatedCheckData{PassingCount: 1}, nil)
	suite.checksService.On("GetAggregatedChecksResultByCluster", "2").Return(&AggregatedCheckData{WarningCount: 1}, nil)
	suite.checksService.On("GetAggregatedChecksResultByCluster", "3").Return(&AggregatedCheckData{CriticalCount: 1}, nil)

	suite.tagsService.On("GetAllByResource", models.TagClusterResourceType, "1").Return([]string{"tag1", "tag2"}, nil)
	suite.tagsService.On("GetAllByResource", models.TagClusterResourceType, "2").Return([]string{"tag3", "tag4"}, nil)
	suite.tagsService.On("GetAllByResource", models.TagClusterResourceType, "3").Return([]string{"tag5", "tag6"}, nil)

	clusters, _ := suite.clustersService.GetAll(map[string][]string{})

	suite.ElementsMatch(models.ClusterList{
		&models.Cluster{
			ID:              "1",
			Name:            "cluster1",
			ClusterType:     models.ClusterTypeScaleUp,
			SIDs:            []string{"DEV"},
			ResourcesNumber: 10,
			HostsNumber:     2,
			Health:          models.CheckPassing,
			Tags:            []string{"tag1", "tag2"},
		},
		&models.Cluster{
			ID:              "2",
			Name:            "cluster2",
			ClusterType:     models.ClusterTypeScaleOut,
			SIDs:            []string{"QAS"},
			ResourcesNumber: 11,
			HostsNumber:     2,
			Health:          models.CheckWarning,
			Tags:            []string{"tag3", "tag4"},
		},
		&models.Cluster{
			ID:              "3",
			Name:            "cluster3",
			ClusterType:     models.ClusterTypeUnknown,
			SIDs:            []string{"PRD", "PRD2"},
			ResourcesNumber: 3,
			HostsNumber:     5,
			Health:          models.CheckCritical,
			Tags:            []string{"tag5", "tag6"},
		},
	}, clusters)
}
