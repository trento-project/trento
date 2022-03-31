package discovery

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/agent/discovery/collector"
	"github.com/trento-project/trento/internal/subscription"
)

const SubscriptionDiscoveryId string = "subscription_discovery"

type SubscriptionDiscovery struct {
	id              string
	collectorClient collector.Client
	host            string
	interval        time.Duration
}

func NewSubscriptionDiscovery(collectorClient collector.Client, config DiscoveriesConfig) (Discovery, error) {
	if config.DiscoveriesPeriodsConfig.Subscription < 1 {
		return nil, fmt.Errorf("invalid interval %s", config.DiscoveriesPeriodsConfig.Subscription)
	}

	r := SubscriptionDiscovery{}
	r.id = SubscriptionDiscoveryId
	r.collectorClient = collectorClient
	r.host, _ = os.Hostname()
	r.interval = config.DiscoveriesPeriodsConfig.Subscription

	return r, nil
}

func (d SubscriptionDiscovery) GetId() string {
	return d.id
}

func (d SubscriptionDiscovery) GetInterval() time.Duration {
	return d.interval
}

func (d SubscriptionDiscovery) Discover() (string, error) {
	subsData, err := subscription.NewSubscriptions()
	if err != nil {
		return "", err
	}

	err = d.collectorClient.Publish(d.id, subsData)
	if err != nil {
		log.Debugf("Error while sending subscription discovery to data collector: %s", err)
		return "", err
	}

	return fmt.Sprintf("Subscription (%d entries) discovered", len(subsData)), nil
}
