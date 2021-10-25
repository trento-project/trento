package services

import (
	log "github.com/sirupsen/logrus"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/subscription"
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
	GetHostSubscriptions(host string) (subscription.Subscriptions, error)
}

type subscriptionsService struct {
	consul consul.Client
}

func NewSubscriptionsService(client consul.Client) SubscriptionsService {
	return &subscriptionsService{consul: client}
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

func (s *subscriptionsService) GetHostSubscriptions(host string) (subscription.Subscriptions, error) {
	subs, err := subscription.Load(s.consul, host)
	if err != nil {
		return nil, err
	}

	return subs, nil
}
