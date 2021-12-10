package services

import (
	"gorm.io/gorm"

	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
)

const (
	SlesIdentifier string = "SLES_SAP"
)

//go:generate mockery --name=SubscriptionsService --inpackage --filename=subscriptions_mock.go
type SubscriptionsService interface {
	IsTrentoPremium() (bool, error)
	GetPremiumData() (*models.PremiumData, error)
	GetHostSubscriptions(host string) ([]*models.SlesSubscription, error)
}

type subscriptionsService struct {
	db *gorm.DB
}

func NewSubscriptionsService(db *gorm.DB) *subscriptionsService {
	return &subscriptionsService{db: db}
}

func (s *subscriptionsService) IsTrentoPremium() (bool, error) {
	premiumData, err := s.GetPremiumData()
	if err != nil {
		return false, err
	}

	return premiumData.IsPremium, nil
}

func (s *subscriptionsService) GetPremiumData() (*models.PremiumData, error) {
	var count int64
	result := s.db.Table("sles_subscriptions").Where("id", SlesIdentifier).Count(&count)
	if result.Error != nil {
		return nil, result.Error
	}

	premiumData := &models.PremiumData{
		IsPremium:     count > 0,
		Sles4SapCount: int(count),
	}

	return premiumData, nil
}

func (s *subscriptionsService) GetHostSubscriptions(id string) ([]*models.SlesSubscription, error) {
	var subEntities []*entities.SlesSubscription
	err := s.db.
		Where("agent_id", id).
		Find(&subEntities).
		Error

	if err != nil {
		return nil, err
	}

	var subModels []*models.SlesSubscription
	for _, sub := range subEntities {
		subModels = append(subModels, sub.ToModel())
	}

	return subModels, nil
}
