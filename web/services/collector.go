package services

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/web/datapipeline"
	"gorm.io/gorm"
)

//go:generate mockery --name=CollectorService --inpackage --filename=collector_mock.go
type CollectorService interface {
	StoreEvent(dataCollected *datapipeline.DataCollectedEvent) error
}

type collectorService struct {
	db                *gorm.DB
	projectorsChannel chan *datapipeline.DataCollectedEvent
}

func NewCollectorService(db *gorm.DB, projectorsChannel chan *datapipeline.DataCollectedEvent) *collectorService {
	return &collectorService{db: db, projectorsChannel: projectorsChannel}
}

// StoreEvent stores the event in the database and sends it to the projectors.
// If the last event is equal to the current one, it is not stored nor sent.
func (c *collectorService) StoreEvent(collectedData *datapipeline.DataCollectedEvent) error {
	return c.db.Transaction(func(tx *gorm.DB) error {
		var event datapipeline.DataCollectedEvent
		err := tx.
			Where(&datapipeline.DataCollectedEvent{
				AgentID:       collectedData.AgentID,
				DiscoveryType: collectedData.DiscoveryType,
				Payload:       collectedData.Payload,
			}).
			Last(&event).
			Error

		if errors.Is(gorm.ErrRecordNotFound, err) {
			err := tx.Create(collectedData).Error
			if err != nil {
				return err
			}
			c.projectorsChannel <- collectedData
			return nil
		}

		if err != nil {
			return err
		}

		log.Debugf("Event already exists. Agent ID: %s, DiscoveryType: %s ", collectedData.AgentID, collectedData.DiscoveryType)
		c.projectorsChannel <- &event

		return nil
	})
}
