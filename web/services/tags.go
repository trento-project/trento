package services

import (
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
)

//go:generate mockery --name=TagsService --inpackage --filename=tags_mock.go
type TagsService interface {
	GetAll(resourceTypeFilter ...string) ([]string, error)
	GetAllByResource(resourceType string, resourceId string) ([]string, error)
	Create(value string, resourceType string, resourceId string) error
	Delete(value string, resourceType string, resourceId string) error
}

type tagsService struct {
	db *gorm.DB
}

func NewTagsService(db *gorm.DB) *tagsService {
	return &tagsService{db: db}
}

func (r *tagsService) GetAll(resourceTypeFilter ...string) ([]string, error) {
	db := r.db
	for _, f := range resourceTypeFilter {
		db = db.Or("resource_type", f)
	}

	return getTags(db)
}

func (r *tagsService) GetAllByResource(resourceType string, resourceId string) ([]string, error) {
	db := r.db.Where("resource_type", resourceType)
	db = db.Where("resource_id", resourceId)

	return getTags(db)
}

func (r *tagsService) Create(value string, resourceType string, resourceId string) error {
	tag := models.Tag{
		Value:        value,
		ResourceType: resourceType,
		ResourceID:   resourceId,
	}

	result := r.db.Create(&tag)

	return result.Error
}

func (r *tagsService) Delete(value string, resourceType string, resourceId string) error {
	tag := models.Tag{
		Value:        value,
		ResourceType: resourceType,
		ResourceID:   resourceId,
	}

	result := r.db.Delete(&tag)

	return result.Error
}

func getTags(db *gorm.DB) ([]string, error) {
	var tags []models.Tag
	result := db.Distinct("value").Find(&tags)

	if result.Error != nil {
		return nil, result.Error
	}

	var tagStrings []string
	for _, t := range tags {
		tagStrings = append(tagStrings, t.Value)
	}

	return tagStrings, nil
}
