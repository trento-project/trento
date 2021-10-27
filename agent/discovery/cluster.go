package discovery

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/agent/collector"
	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
)

const ClusterDiscoveryId string = "ha_cluster_discovery"

// This Discover handles any Pacemaker Cluster type
type ClusterDiscovery struct {
	id        string
	discovery BaseDiscovery
	Cluster   cluster.Cluster
}

func NewClusterDiscovery(consulClient consul.Client, collectorClient collector.Client) ClusterDiscovery {
	d := ClusterDiscovery{}
	d.id = ClusterDiscoveryId
	d.discovery = NewDiscovery(consulClient, collectorClient)
	return d
}

func (c ClusterDiscovery) GetId() string {
	return c.id
}

// Execute one iteration of a discovery and store the result in the Consul KVStore.
func (d ClusterDiscovery) Discover() (string, error) {
	cluster, err := cluster.NewCluster()
	if err != nil {
		return "No HA cluster discovered on this host", nil
	}

	d.Cluster = cluster

	err = d.Cluster.Store(d.discovery.consulClient)
	if err != nil {
		return "", err
	}

	err = storeClusterMetadata(d.discovery.consulClient, cluster.Name, cluster.Id)
	if err != nil {
		return "", err
	}

	err = d.discovery.collectorClient.Publish(d.id, cluster)
	if err != nil {
		log.Debugf("Error while sending cluster discovery to data collector: %s", err)
		return "", err
	}

	return fmt.Sprintf("Cluster with name: %s successfully discovered", cluster.Name), nil
}

func storeClusterMetadata(client consul.Client, clusterName string, clusterId string) error {
	metadata := hosts.Metadata{
		Cluster:   clusterName,
		ClusterId: clusterId,
	}
	err := metadata.Store(client)
	if err != nil {
		return err
	}

	return nil
}
