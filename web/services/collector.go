package services

import (
	"time"

	"github.com/trento-project/trento/web/datapipeline"
	"gorm.io/gorm"
)

//go:generate mockery --name=CollectorService --inpackage --filename=collector_mock.go
type CollectorService interface {
	StoreEvent(dataCollected *datapipeline.DataCollectedEvent) error
	PruneEvents(olderThan time.Duration) error
}

type collectorService struct {
	db                *gorm.DB
	projectorsChannel chan *datapipeline.DataCollectedEvent
}

func NewCollectorService(db *gorm.DB, projectorsChannel chan *datapipeline.DataCollectedEvent) *collectorService {
	return &collectorService{db: db, projectorsChannel: projectorsChannel}
}

func (c *collectorService) StoreEvent(collectedData *datapipeline.DataCollectedEvent) error {
	if err := c.db.Create(collectedData).Error; err != nil {
		return err
	}
	c.projectorsChannel <- collectedData

	return nil
}

func (c *collectorService) PruneEvents(olderThan time.Duration) error {
	return c.db.Delete(datapipeline.DataCollectedEvent{}, "created_at < ?", time.Now().Add(-olderThan)).Error
}
