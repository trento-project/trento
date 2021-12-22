package entities

import (
	"time"

	"gorm.io/datatypes"
)

type ChecksResult struct {
	ID        int64
	CreatedAt time.Time
	GroupID   string
	Payload   datatypes.JSON
}
