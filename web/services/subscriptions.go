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

func (s *subscriptionsService) GetHostSubscriptions(host string) ([]*models.SlesSubscription, error) {
	// Get the agent id by host name. This should be removed once the host page uses the agent id
	// to go the host details page
	var hostEntity *entities.Host
	result := s.db.Where("name", host).Find(&hostEntity)
	if result.Error != nil {
		return nil, result.Error
	}
	agent_id := hostEntity.ToModel().ID

	var subEntities []*entities.SlesSubscription
	result = s.db.Where("agent_id", agent_id).Find(&subEntities)
	if result.Error != nil {
		return nil, result.Error
	}

	var subModels []*models.SlesSubscription
	for _, sub := range subEntities {
		subModels = append(subModels, sub.ToModel())
	}

	return subModels, nil
}
