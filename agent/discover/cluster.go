package discover

import (
	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/environments"
)

const ClusterDiscoverId string = "discover_cluster"

// This Discover handles any Pacemaker Cluster type
type ClusterDiscover struct {
	id      string
	host    Discover
	Cluster cluster.Cluster
}

func (cluster ClusterDiscover) GetId() string {
	return cluster.id
}

// check if the current node this trento agent is running on can be discovered
// by ClusterDiscover
func (cluster ClusterDiscover) ShouldDiscover(client consul.Client) bool {
	// ### Check if we have cibadmin available
	return true
}

// Create or Updating the given Consul Key-Value Path Store with a new value from the Agent
func (cluster ClusterDiscover) storeDiscovery(cStorePath, cStoreValue string) error {
	return nil
}

// Execute one iteration of a discovery and store the result in the Consul KVStore.
func (discover ClusterDiscover) Discover() error {
	cluster, err := cluster.NewCluster()

	if err != nil {
		return err
	}

	discover.Cluster = cluster

	err = discover.Cluster.Store(discover.host.client)
	if err != nil {
		return err
	}

	err = storeClusterMetadata(discover.host.client, cluster.Name())
	if err != nil {
		return err
	}

	return nil
}

func NewClusterDiscover(client consul.Client) ClusterDiscover {
	r := ClusterDiscover{}
	r.id = ClusterDiscoverId
	r.host = NewDiscover(client)
	return r
}

func storeClusterMetadata(client consul.Client, clusterName string) error {
	metadata := environments.NewMetadata()
	metadata.Cluster = clusterName
	err := metadata.Store(client)
	if err != nil {
		return err
	}

	return nil
}
