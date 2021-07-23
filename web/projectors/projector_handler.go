package projectors

import "gorm.io/gorm"

//go:generate mockery --all

// ProjectorHandler is the projector handler interface
// GetName returns the projector name (used in the subscription)
// Query gets data from the source
// Project projects data to the database
type ProjectorHandler interface {
	GetName() string
	Query(lastSeenIndex uint64) (data interface{}, lastIndex uint64, err error)
	Project(db *gorm.DB, data interface{}) error
}
