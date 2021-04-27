package discover

import (
	"fmt"
	"log"
	"strconv"

	// Reusing the Prometheus Ha Exporter cibadmin xml parser here
	"github.com/ClusterLabs/ha_cluster_exporter/collector/pacemaker/cib"

	consul "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	consul_internal "github.com/trento-project/trento/internal/consul"
)

const crm_monPath string = "/usr/sbin/crm_mon"
const cib_admPath string = "/usr/sbin/cibadmin"

//const crm_monPath string = "./test/fake/crm_mon.sh"
//const cib_admPath string = "./test/fake_cibadmin.sh"

// This Discover handles any Pacemaker Cluster type
type ClusterDiscover struct {
	host           Discover
	cibConfig      cib.Root
	clusterName    string
	stonithEnabled bool
	stonithType    consul_internal.ClusterStonithType
}

// check if the current node this trento agent is running on can be discovered
// by ClusterDiscover
func (cluster ClusterDiscover) ShouldDiscover(client consul.Client) bool {
	// ### Check if we have cibadmin available
	return true
}

// Create or Updating the given Consul Key-Value Path Store with a new value from the Agent
func (cluster ClusterDiscover) storeDiscovery(cStorePath, cStoreValue string) error {
	kvPath := fmt.Sprintf("%s/%s/%s", consul_internal.KvClustersPath, cluster.clusterName, cStorePath)

	_, err := cluster.host.client.KV().Put(&consul.KVPair{
		Key:   kvPath,
		Value: []byte(cStoreValue)}, nil)
	return err
}

// Execute one iteration of a discovery and store the result in the Consul KVStore.
func (cluster ClusterDiscover) Discover() error {
	err := cluster.host.Discover()
	if err != nil {
		return err
	}
	cluster.stonithType = consul_internal.ClusterStonithUnknown

	cibParser := cib.NewCibAdminParser(cib_admPath)

	cibConfig, err := cibParser.Parse()
	if err != nil {
		log.Printf("Failing to parse: %s", errors.Wrap(err, "cibadmin parser error"))
		return err
	}
	cluster.cibConfig = cibConfig

	for _, prop := range cluster.cibConfig.Configuration.CrmConfig.ClusterProperties {
		switch prop.Id {
		case "cib-bootstrap-options-cluster-name":
			cluster.clusterName = prop.Value
		case "cib-bootstrap-options-stonith-enabled":
			cluster.stonithEnabled = bool(prop.Value == "true")
		}
	}

	var foundVIP bool
	for _, primitive := range cluster.cibConfig.Configuration.Resources.Primitives {
		switch primitive.Type {
		case "external/sbd":
			cluster.stonithType = consul_internal.ClusterStonithSBD
		case "IPaddr2":
			if !foundVIP && primitive.Provider == "heartbeat" {
				for _, attr := range primitive.InstanceAttributes {
					switch attr.Name {
					case "ip":
						cluster.storeDiscovery("VIP_primary", string(attr.Value))
						foundVIP = true
					case "cidr_netmask":
						cluster.storeDiscovery("VIP_cidr_netmask", attr.Value)
					}
				}
			}
		}
	}

	cluster.storeDiscovery("stonith_enabled", strconv.FormatBool(cluster.stonithEnabled))
	cluster.storeDiscovery("stonith_type", strconv.FormatUint(uint64(cluster.stonithType), 10))
	return nil
}

func NewClusterDiscover(client consul.Client) ClusterDiscover {
	r := ClusterDiscover{}
	r.host = NewDiscover(client)
	return r
}
