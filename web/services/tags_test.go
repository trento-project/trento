package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/web/models"

	"gorm.io/gorm"
)

func setupTestTagsService() *gorm.DB {
	db := setupTestDatabase()
	db.AutoMigrate(models.Tag{})

	db = db.Exec("TRUNCATE TABLE tags")
	if db.Error != nil {
		panic(db.Error)
	}

	loadTagsFixtures(db)
	return db
}

func tearDownTestTagsService(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE tags")
}

func TestTagsService_GetAll(t *testing.T) {
	db := setupTestTagsService()
	defer tearDownTestTagsService(db)

	tagsService := NewTagsService(db)
	tags, _ := tagsService.GetAll()

	assert.ElementsMatch(t, []string{"tag1", "tag2", "tag3"}, tags)
}

func TestTagsService_GetAll_Filter(t *testing.T) {
	db := setupTestTagsService()
	defer tearDownTestTagsService(db)

	tagsService := NewTagsService(db)
	tags, _ := tagsService.GetAll(models.TagClusterResourceType, models.TagHostResourceType)

	assert.ElementsMatch(t, []string{"tag2", "tag3"}, tags)
}

func TestTagsService_GetAllByResource(t *testing.T) {
	db := setupTestTagsService()
	defer tearDownTestTagsService(db)

	tagsService := NewTagsService(db)
	tags, _ := tagsService.GetAllByResource(models.TagHostResourceType, "suse")

	assert.ElementsMatch(t, []string{"tag3"}, tags)
}

func TestTagsService_Create(t *testing.T) {
	db := setupTestTagsService()
	tx := db.Begin()

	defer func() {
		tx.Rollback()
		tearDownTestTagsService(db)
	}()

	tagsService := NewTagsService(tx)
	tagsService.Create("newtag", models.TagHostResourceType, "suse")

	expectedTag := models.Tag{Value: "newtag", ResourceType: models.TagHostResourceType, ResourceId: "suse"}

	var tags []models.Tag
	tx.Where(&expectedTag).Find(&tags)

	assert.ElementsMatch(t, []models.Tag{expectedTag}, tags)
}

func TestTagsService_Delete(t *testing.T) {
	db := setupTestTagsService()
	tx := db.Begin()

	defer func() {
		tx.Rollback()
		tearDownTestTagsService(db)
	}()

	tagsService := NewTagsService(tx)
	tagsService.Delete("tag1", models.TagSAPSystemResourceType, "HA1")

	var tags []models.Tag
	tx.Find(&tags)

	assert.Equal(t, 3, len(tags))
	assert.ElementsMatch(t, []models.Tag{
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
