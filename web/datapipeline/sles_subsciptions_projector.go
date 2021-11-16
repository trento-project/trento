package datapipeline

import (
	log "github.com/sirupsen/logrus"

	"github.com/trento-project/trento/internal/subscription"
	"github.com/trento-project/trento/web/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewSlesSubscriptionsProjector(db *gorm.DB) *projector {
	subsProjector := NewProjector("sles_subscriptions", db)

	subsProjector.AddHandler(SubscriptionDiscovery, subsProjector_SubscriptionDiscoveryHandler)

	return subsProjector
}

func subsProjector_SubscriptionDiscoveryHandler(dataCollectedEvent *DataCollectedEvent, db *gorm.DB) error {
	decoder := getPayloadDecoder(dataCollectedEvent.Payload)

	var discoveredSubscriptions subscription.Subscriptions
	if err := decoder.Decode(&discoveredSubscriptions); err != nil {
		log.Errorf("can't decode data: %s", err)
		return err
	}

	var subEntities []entities.SlesSubscription

	for _, subscription := range discoveredSubscriptions {
		subEntity := entities.SlesSubscription{
			AgentID:            dataCollectedEvent.AgentID,
			ID:                 subscription.Identifier,
			Version:            subscription.Version,
			Type:               subscription.Type,
			Arch:               subscription.Arch,
			Status:             subscription.Status,
			StartsAt:           subscription.StartsAt,
			ExpiresAt:          subscription.ExpiresAt,
			SubscriptionStatus: subscription.SubscriptionStatus,
		}

		subEntities = append(subEntities, subEntity)
	}

	return db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&subEntities).Error
}
