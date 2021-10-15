package datapipeline

import "time"

// Subscription is a cursor of a projector to the stream of the events
type Subscription struct {
	LastProjectedEventID int64
	AgentID              string `gorm:"primaryKey"`
	ProjectorID          string `gorm:"primaryKey"`
	UpdatedAt            time.Time
}
