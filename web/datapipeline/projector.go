package datapipeline

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProjectorHandler func(dataCollectedEvent *DataCollectedEvent, db *gorm.DB) error

type Projector struct {
	ID       string
	db       *gorm.DB
	handlers map[string]ProjectorHandler
}

func NewProjector(ID string, db *gorm.DB) *Projector {
	return &Projector{
		ID:       ID,
		db:       db,
		handlers: make(map[string]ProjectorHandler),
	}
}

// AddHandler registers a handler for a specific discovery type
func (p *Projector) AddHandler(discoveryType string, handler ProjectorHandler) {
	p.handlers[discoveryType] = handler
}

// Project processes the data collected event and calls the registered handlers
// By updating the subscription with the LastProjectedEventID, it leverages the PostgresSQL implicit lock
// to enforce linearizability if a specific agent tries to use the same projector concurrently
func (p *Projector) Project(dataCollectedEvent *DataCollectedEvent) error {
	handler, ok := p.handlers[dataCollectedEvent.DiscoveryType]

	if !ok {
		log.Infof("Projector: %s is not interested in %s. Discarding event: %d", p.ID, dataCollectedEvent.DiscoveryType, dataCollectedEvent.ID)
		return nil
	}

	log.Infof("Projector: %s is interested in %s. Projecting event: %d", p.ID, dataCollectedEvent.DiscoveryType, dataCollectedEvent.ID)
	return p.db.Transaction(func(tx *gorm.DB) error {
		tx.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&Subscription{
			ProjectorID:          p.ID,
			AgentID:              dataCollectedEvent.AgentID,
			LastProjectedEventID: dataCollectedEvent.ID,
		})

		err := handler(dataCollectedEvent, tx)
		if err != nil {
			return err
		}

		return nil
	})
}
