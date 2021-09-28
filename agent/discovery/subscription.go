package discovery

import (
	"fmt"

	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/subscription"
)

const SubscriptionDiscoveryId string = "subscription_discovery"

type SubscriptionDiscovery struct {
	id        string
	discovery BaseDiscovery
}

func NewSubscriptionDiscovery(client consul.Client) SubscriptionDiscovery {
	r := SubscriptionDiscovery{}
	r.id = SubscriptionDiscoveryId
	r.discovery = NewDiscovery(client)
	return r
}

func (d SubscriptionDiscovery) GetId() string {
	return d.id
}

func (d SubscriptionDiscovery) Discover() (string, error) {
	subsData, err := subscription.NewSubscriptions()
	if err != nil {
		return "", err
	}

	err = subsData.Store(d.discovery.client)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Subscription (%d entries) discovered", len(subsData)), nil
}
