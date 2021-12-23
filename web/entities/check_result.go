package entities

import (
	"encoding/json"
	"time"

	"github.com/trento-project/trento/web/models"
	"gorm.io/datatypes"
)

type ChecksResult struct {
	ID        int64
	CreatedAt time.Time
	GroupID   string
	Payload   datatypes.JSON
}

func (c *ChecksResult) ToModel() (*models.ChecksResult, error) {
	var checkResult models.ChecksResult
	checkResult.ID = c.GroupID
	err := json.Unmarshal(c.Payload, &checkResult)

	return &checkResult, err
}
