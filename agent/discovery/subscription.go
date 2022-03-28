package discovery

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/agent/discovery/collector"
	"github.com/trento-project/trento/internal/subscription"
)

const SubscriptionDiscoveryId string = "subscription_discovery"

type SubscriptionDiscovery struct {
	id        string
	discovery BaseDiscovery
}

func NewSubscriptionDiscovery(collectorClient collector.Client, interval time.Duration) SubscriptionDiscovery {
	r := SubscriptionDiscovery{}
	r.id = SubscriptionDiscoveryId
	r.discovery = NewDiscovery(collectorClient, interval)
	return r
}

func (d SubscriptionDiscovery) GetId() string {
	return d.id
}

func (d SubscriptionDiscovery) GetInterval() time.Duration {
	return d.discovery.interval
}

func (d SubscriptionDiscovery) Discover() (string, error) {
	subsData, err := subscription.NewSubscriptions()
	if err != nil {
		return "", err
	}

	err = d.discovery.collectorClient.Publish(d.id, subsData)
	if err != nil {
		log.Debugf("Error while sending subscription discovery to data collector: %s", err)
		return "", err
	}

	return fmt.Sprintf("Subscription (%d entries) discovered", len(subsData)), nil
}
