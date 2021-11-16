package entities

import (
	"github.com/trento-project/trento/web/models"
)

type SlesSubscription struct {
	AgentID            string `gorm:"primaryKey"`
	ID                 string `gorm:"primaryKey"`
	Version            string
	Type               string
	Arch               string
	Status             string
	StartsAt           string
	ExpiresAt          string
	SubscriptionStatus string
}

func (s *SlesSubscription) ToModel() *models.SlesSubscription {
	return &models.SlesSubscription{
		ID:                 s.ID,
		Version:            s.Version,
		Type:               s.Type,
		Arch:               s.Arch,
		Status:             s.Status,
		StartsAt:           s.StartsAt,
		ExpiresAt:          s.ExpiresAt,
		SubscriptionStatus: s.SubscriptionStatus,
	}
}
