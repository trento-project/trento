package models

const (
	TagHostResourceType      = "hosts"
	TagClusterResourceType   = "clusters"
	TagSAPSystemResourceType = "sapsystems"
	TagDatabaseResourceType  = "databases"
)

type Tag struct {
	Value        string `gorm:"primaryKey"`
	ResourceID   string `gorm:"primaryKey"`
	ResourceType string `gorm:"primaryKey"`
}
