package entities

import (
	"encoding/json"
	"time"

	"github.com/trento-project/trento/web/models"
	"gorm.io/datatypes"
)

type Check struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	Payload   datatypes.JSON
}

type CheckList []*Check

func (c *Check) ToModel() (*models.Check, error) {
	var check models.Check
	err := json.Unmarshal(c.Payload, &check)
	if err != nil {
		return nil, err
	}

	return &check, nil
}

func (c CheckList) ToModel() (models.ChecksCatalog, error) {
	var checksCatalog models.ChecksCatalog

	for _, checkRaw := range c {
		check, err := checkRaw.ToModel()
		if err != nil {
			return nil, err
		}
		checksCatalog = append(checksCatalog, check)
	}

	return checksCatalog, nil
}
