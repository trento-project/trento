package datapipeline

import (
	"bytes"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockery --name=Projector --inpackage --filename=projector_mock.go
type Projector interface {
	Project(dataCollectedEvent *DataCollectedEvent) error
}

type ProjectorHandler func(dataCollectedEvent *DataCollectedEvent, db *gorm.DB) error

type projector struct {
	ID       string
	db       *gorm.DB
	handlers map[string]ProjectorHandler
}

func NewProjector(ID string, db *gorm.DB) *projector {
	return &projector{
		ID:       ID,
		db:       db,
		handlers: make(map[string]ProjectorHandler),
	}
}

// AddHandler registers a handler for a specific discovery type
func (p *projector) AddHandler(discoveryType string, handler ProjectorHandler) {
	p.handlers[discoveryType] = handler
}

// Project processes the data collected event and calls the registered handlers
// By updating the subscription with the LastProjectedEventID, it leverages the PostgresSQL implicit lock
// to enforce linearizability if a specific agent tries to use the same projector concurrently
func (p *projector) Project(dataCollectedEvent *DataCollectedEvent) error {
	handler, ok := p.handlers[dataCollectedEvent.DiscoveryType]

	if !ok {
		log.Debugf("Projector: %s is not interested in %s. Discarding event: %d", p.ID, dataCollectedEvent.DiscoveryType, dataCollectedEvent.ID)
		return nil
	}

	log.Infof("Projector: %s is interested in %s. Projecting event: %d", p.ID, dataCollectedEvent.DiscoveryType, dataCollectedEvent.ID)

	return p.db.Transaction(func(tx *gorm.DB) error {
		var subscription Subscription
		tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&Subscription{ProjectorID: p.ID, AgentID: dataCollectedEvent.AgentID}).First(&subscription)

		if subscription.LastProjectedEventID >= dataCollectedEvent.ID {
			log.Warnf("Projector: %s received an old event: %s %s %d. Skipping",
				p.ID, dataCollectedEvent.DiscoveryType,
				dataCollectedEvent.AgentID,
				dataCollectedEvent.ID)

			return nil
		}

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

func getPayloadDecoder(payload datatypes.JSON) *json.Decoder {
	data, _ := payload.MarshalJSON()
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	return decoder
}
