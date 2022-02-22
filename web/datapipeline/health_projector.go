package datapipeline

import (
	"encoding/json"
	"errors"

	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ProjectHealth(db *gorm.DB, healthID, healthType, healthValue string) error {
	var healthState entities.HealthState
	var partialHealths map[string]string

	err := db.Where("id = ?", healthID).First(&healthState).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			healthState.ID = healthID
			partialHealths = make(map[string]string)
		} else {
			return err
		}
	} else {
		err = json.Unmarshal(healthState.PartialHealths, &partialHealths)
		if err != nil {
			return err
		}
	}

	partialHealths[healthType] = healthValue

	partialHealthsJson, _ := json.Marshal(partialHealths)
	healthState.Health = computeOverallHealth(partialHealths)
	healthState.PartialHealths = (datatypes.JSON)(partialHealthsJson)

	return db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(healthState).Error
}

func computeOverallHealth(partialHealths map[string]string) string {
	health := models.HealthSummaryHealthPassing
	for _, pHealth := range partialHealths {
		switch {
		case pHealth == models.HealthSummaryHealthCritical:
			health = models.HealthSummaryHealthCritical
		case health != models.HealthSummaryHealthCritical && pHealth == models.HealthSummaryHealthWarning:
			health = models.HealthSummaryHealthWarning
		case health == models.HealthSummaryHealthPassing && pHealth == models.HealthSummaryHealthUnknown:
			health = models.HealthSummaryHealthUnknown
		}
	}

	return health
}
