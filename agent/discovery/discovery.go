package discovery

import (
	"time"

	"github.com/trento-project/trento/agent/discovery/collector"
)

type DiscoveriesPeriodConfig struct {
	Cluster      time.Duration
	SAPSystem    time.Duration
	Cloud        time.Duration
	Host         time.Duration
	Subscription time.Duration
}

type DiscoveriesConfig struct {
	SSHAddress               string
	DiscoveriesPeriodsConfig *DiscoveriesPeriodConfig
	CollectorConfig          *collector.Config
}

type Discovery interface {
	// Returns an arbitrary unique string identifier of the discovery
	GetId() string
	// Execute the discovery mechanism
	Discover() (string, error)
	// Get interval
	GetInterval() time.Duration
}
