package services

import (
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/subscription"
)

//go:generate mockery --name=SubscriptionsService

type SubscriptionsService interface {
	//GetSubscriptionType() (string, error)
	GetHostSubscriptions(host string) (subscription.Subscriptions, error)
}

type subscriptionsService struct {
	consul consul.Client
}

func NewSubscriptionsService(client consul.Client) SubscriptionsService {
	return &subscriptionsService{consul: client}
}

func (s *subscriptionsService) GetHostSubscriptions(host string) (subscription.Subscriptions, error) {
	sub, err := subscription.Load(s.consul, host)
	if err != nil {
		return nil, err
	}

	return sub, nil
}
