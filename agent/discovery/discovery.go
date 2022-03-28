package discovery

import (
	"os"
	"time"

	"github.com/trento-project/trento/agent/discovery/collector"
)

type Discovery interface {
	// Returns an arbitrary unique string identifier of the discovery
	GetId() string
	// Execute the discovery mechanism
	Discover() (string, error)
	// Get interval
	GetInterval() time.Duration
}

type BaseDiscovery struct {
	id              string
	collectorClient collector.Client
	host            string
	interval        time.Duration
}

func (d BaseDiscovery) GetId() string {
	return d.id
}

// Execute one iteration of a discovery
func (d BaseDiscovery) Discover() (string, error) {
	d.host, _ = os.Hostname()
	return "Basic discovery example", nil
}

func (d BaseDiscovery) GetInterval() time.Duration {
	return d.interval
}

// NewDiscovery Return a new base discovery with the support for data collector endpoint
func NewDiscovery(collectorClient collector.Client, interval time.Duration) BaseDiscovery {
	d := BaseDiscovery{}
	d.id = ""
	d.collectorClient = collectorClient
	d.host, _ = os.Hostname()
	d.interval = interval
	return d
}
