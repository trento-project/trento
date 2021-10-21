package discovery

import (
	"os"

	"github.com/trento-project/trento/agent/collector"
	"github.com/trento-project/trento/internal/consul"
)

type Discovery interface {
	// Returns an arbitrary unique string identifier of the discovery, so that we can associate it to a Consul check ID
	GetId() string
	// Execute the discovery mechanism
	Discover() (string, error)
}

type BaseDiscovery struct {
	id              string
	client          consul.Client
	collectorClient collector.Client
	host            string
}

func (d BaseDiscovery) GetId() string {
	return d.id
}

// Execute one iteration of a discovery and store the result in the Consul KVStore.
func (d BaseDiscovery) Discover() (string, error) {
	d.host, _ = os.Hostname()
	return "Basic discovery example", nil
}

// NewDiscovery Return a new base discovery
func NewDiscovery(client consul.Client) BaseDiscovery {
	d := BaseDiscovery{}
	d.id = ""
	d.client = client
	d.host, _ = os.Hostname()
	return d
}

// NewDiscovery Return a new base discovery with the support for consul storage and data collector endpoint
func NewDiscoveryWithCollector(client consul.Client, collectorClient collector.Client) BaseDiscovery {
	d := NewDiscovery(client)
	d.collectorClient = collectorClient
	return d
}
