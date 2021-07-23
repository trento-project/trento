package projectors

import (
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const retryAfter = 30 * time.Second

// Subscription persists the last seen index of the projector
type Subscription struct {
	ID            int
	Projector     string `gorm:"uniqueIndex"`
	LastSeenIndex uint64
}

func getOrCreateSubscription(db *gorm.DB, projector string) (*Subscription, error) {
	var subscription Subscription
	var err error

	if err = db.Where(&Subscription{Projector: projector}).First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			subscription.Projector = projector
			err = db.Create(&subscription).Error
		}
	}

	return &subscription, err
}

func updateLastSeenIndex(db *gorm.DB, subscription *Subscription, lastSeenIndex uint64) error {
	return db.Model(subscription).Update("last_seen_index", lastSeenIndex).Error
}

type Projector struct {
	handler      ProjectorHandler
	db           *gorm.DB
	subscription *Subscription
}

// NewProjector returns a new projector
func NewProjector(handler ProjectorHandler, db *gorm.DB) *Projector {
	return &Projector{
		handler: handler,
		db:      db,
	}
}

// Run gets the projector subscription from the database by using the handler name
// Spawns a go routine that runs the projector transaction every time data is received,
// by calling the handler Query function
func (p *Projector) Run() {
	var err error
	name := p.handler.GetName()
	p.subscription, err = getOrCreateSubscription(p.db, name)

	if err != nil {
		log.Fatal("Error retrieving hosts projectors subscription: ", err)
	}
	log.Debugf("Projector %s subscription found, last seen index: %v", name, p.subscription.LastSeenIndex)

	go func() {
		for {
			lastSeenIndex := p.subscription.LastSeenIndex
			data, lastIndex, err := p.handler.Query(lastSeenIndex)

			if err != nil {
				log.Errorf("Projector %s: error while querying: %e", name, err)
				time.Sleep(retryAfter)
				continue
			}

			if lastIndex > p.subscription.LastSeenIndex {
				log.Debugf("Changes detected, projecting %s. New index: %v", name, lastIndex)

				if err := p.Project(data, lastIndex); err != nil {
					p.subscription.LastSeenIndex = lastSeenIndex
					log.Errorf("Projector %s: error while projecting: %e", name, err)
					time.Sleep(retryAfter)
				}
			}
		}
	}()
}

// Project calls the project handler Project function in a transaction and updates the lastSeenIndex
// returns an error otherwise
func (p *Projector) Project(data interface{}, lastSeenIndex uint64) error {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if err := p.handler.Project(tx, data); err != nil {
			return err
		}

		if err := updateLastSeenIndex(tx, p.subscription, lastSeenIndex); err != nil {
			return err
		}
		return nil
	})

	return err
}
