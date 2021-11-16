package services

import (
	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/subscription"

	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
)

const (
	SlesIdentifier string = "SLES_SAP"
	Premium        string = "Premium"
	Free           string = "Free"
)

type SubscriptionData struct {
	Type            string
	SubscribedCount int
}

//go:generate mockery --name=SubscriptionsService --inpackage --filename=subscriptions_mock.go
type SubscriptionsService interface {
	GetSubscriptionData() (*SubscriptionData, error)
	GetHostSubscriptions(host string) ([]*models.SlesSubscription, error)
}

type subscriptionsService struct {
	db *gorm.DB
	consul consul.Client
}

func NewSubscriptionsService(db *gorm.DB, client consul.Client) SubscriptionsService {
	return &subscriptionsService{db: db, consul: client}
}

func (s *subscriptionsService) GetSubscriptionData() (*SubscriptionData, error) {
	query := &consulApi.QueryOptions{}
	consulNodes, _, err := s.consul.Catalog().Nodes(query)
	if err != nil {
		return nil, err
	}

	var subData = &SubscriptionData{Type: Free}

	for _, node := range consulNodes {
		subs, err := subscription.Load(s.consul, node.Node)
		if err != nil {
			log.Errorf("Couldn't get subscriptions data from node %s", node.Node)
			continue
		}

		for _, sub := range subs {
			if sub.Identifier == SlesIdentifier {
				subData.SubscribedCount += 1
				subData.Type = Premium
			}
		}
	}

	return subData, nil
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
