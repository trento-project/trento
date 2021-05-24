package discover

import (
	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
)

const ClusterDiscoverId string = "discover_cluster"

// This Discover handles any Pacemaker Cluster type
type ClusterDiscover struct {
	id      string
	host    Discover
	Cluster cluster.Cluster
}

func NewClusterDiscover(client consul.Client) ClusterDiscover {
	r := ClusterDiscover{}
	r.id = ClusterDiscoverId
	r.host = NewDiscover(client)
	return r
}

func (c ClusterDiscover) GetId() string {
	return c.id
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

func storeClusterMetadata(client consul.Client, clusterName string) error {
	metadata := hosts.Metadata{
		Cluster: clusterName,
	}
	err := metadata.Store(client)
	if err != nil {
		return err
	}

	return nil
}
