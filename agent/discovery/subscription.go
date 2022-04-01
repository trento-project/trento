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
const SubscriptionDiscoveryMinPeriod time.Duration = 20 * time.Second

type SubscriptionDiscovery struct {
	id              string
	collectorClient collector.Client
	host            string
	interval        time.Duration
}

func NewSubscriptionDiscovery(collectorClient collector.Client, config DiscoveriesConfig) Discovery {
	d := SubscriptionDiscovery{}
	d.id = SubscriptionDiscoveryId
	d.collectorClient = collectorClient
	d.host, _ = os.Hostname()
	d.interval = config.DiscoveriesPeriodsConfig.Subscription

	return d
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
