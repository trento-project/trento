package discover

import (
	"os"

	"github.com/trento-project/trento/internal/consul"
)

type Discovery func() (Discoverer, error)

type Discoverer interface {
	// Returns an arbitrary unique string identifier of the discovery, so that we can associate it to a Consul check ID
	GetId() string
	// Execute the discovery mechanism
	Discover() error
}

type Discover struct {
	id     string
	client consul.Client
	host   string
}

func (discover Discover) GetId() string {
	return discover.id
}

// Execute one iteration of a discovery and store the result in the Consul KVStore.
func (discover Discover) Discover() error {
	discover.host, _ = os.Hostname()
	return nil
}

// Return a Host Discover instance
func NewDiscover(client consul.Client) Discover {
	r := Discover{}
	r.id = ""
	r.client = client
	r.host, _ = os.Hostname()
	return r
}
