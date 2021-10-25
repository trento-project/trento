package services

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/models"

	"gorm.io/gorm"
)

type TagsServiceTestSuite struct {
	suite.Suite
	db          *gorm.DB
	tx          *gorm.DB
	tagsService *tagsService
}

func TestTagsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TagsServiceTestSuite))
}

func (suite *TagsServiceTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase()

	suite.db.AutoMigrate(models.Tag{})
	loadTagsFixtures(suite.db)
}

func (suite *TagsServiceTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(models.Tag{})
}

func (suite *TagsServiceTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
	suite.tagsService = NewTagsService(suite.tx)
}

func (suite *TagsServiceTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func loadTagsFixtures(db *gorm.DB) {
	db.Create(&models.Tag{
		ResourceType: models.TagSAPSystemResourceType,
		ResourceId:   "HA1",
		Value:        "tag1",
	})
	db.Create(&models.Tag{
		ResourceType: models.TagSAPSystemResourceType,
		ResourceId:   "HA2",
		Value:        "tag1",
	})
	db.Create(&models.Tag{
		ResourceType: models.TagClusterResourceType,
		ResourceId:   "cluster_id",
		Value:        "tag2",
	})
	db.Create(&models.Tag{
		ResourceType: models.TagHostResourceType,
		ResourceId:   "suse",
		Value:        "tag3",
	})
}

func (suite *TagsServiceTestSuite) TestTagsService_GetAll() {
	tags, _ := suite.tagsService.GetAll()

	suite.ElementsMatch([]string{"tag1", "tag2", "tag3"}, tags)
}

func (suite *TagsServiceTestSuite) TestTagsService_GetAll_Filter() {
	tags, _ := suite.tagsService.GetAll(models.TagClusterResourceType, models.TagHostResourceType)

	suite.ElementsMatch([]string{"tag2", "tag3"}, tags)
}

func (suite *TagsServiceTestSuite) TestTagsService_GetAllByResource() {
	tags, _ := suite.tagsService.GetAllByResource(models.TagHostResourceType, "suse")

	suite.ElementsMatch([]string{"tag3"}, tags)
}

func (suite *TagsServiceTestSuite) TestTagsService_Create() {
	suite.tagsService.Create("newtag", models.TagHostResourceType, "suse")

	expectedTag := models.Tag{Value: "newtag", ResourceType: models.TagHostResourceType, ResourceId: "suse"}

	var tags []models.Tag
	suite.tx.Where(&expectedTag).Find(&tags)

	suite.ElementsMatch([]models.Tag{expectedTag}, tags)
}

func (suite *TagsServiceTestSuite) TestTagsService_Delete() {
	suite.tagsService.Delete("tag1", models.TagSAPSystemResourceType, "HA1")

	var tags []models.Tag
	suite.tx.Find(&tags)

	suite.Equal(3, len(tags))
	suite.ElementsMatch([]models.Tag{
		{
			ResourceType: models.TagSAPSystemResourceType,
			ResourceId:   "HA2",
			Value:        "tag1",
		},
		{
			ResourceType: models.TagClusterResourceType,
			ResourceId:   "cluster_id",
			Value:        "tag2",
		},
		{
			ResourceType: models.TagHostResourceType,
			ResourceId:   "suse",
			Value:        "tag3",
		}}, tags)
}
