package discover

import (
	"github.com/trento-project/trento/internal/consul"
)

type HanaClusterDiscover struct {
	cluster               ClusterDiscover
	SAPid                 string
	SAPHanaInstanceNumber string
}

func (hana HanaClusterDiscover) GetId() string {
	return hana.cluster.id
}

// check if the current node this trento agent is running on can be discovered
// by HanaClusterDiscover, with other words is a SAP Hana Pacemaker cluster
func (hana HanaClusterDiscover) ShouldDiscover(client consul.Client) bool {

	return len(hana.SAPid) > 0
}

// Create or Updating the given Consul Key-Value Path Store with a new value from the Agent
func (hana HanaClusterDiscover) storeDiscovery(cStorePath, cStoreValue string) error {
	return hana.cluster.storeDiscovery(cStorePath, cStoreValue)
}

// Execute one iteration of a discovery and store the result in the Consul KVStore.
func (hana HanaClusterDiscover) Discover() error {
	err := hana.cluster.Discover()
	if err != nil {
		return err
	}
	for _, master := range hana.cluster.cibConfig.Configuration.Resources.Masters {
		for _, attr := range master.Primitive.InstanceAttributes {
			if master.Primitive.Class == "ocf" &&
				master.Primitive.Type == "SAPHana" {
				switch attr.Name {
				case "SID":
					hana.SAPid = attr.Value
					hana.storeDiscovery("SAPID", hana.SAPid)
				case "InstanceNumber":
					hana.SAPHanaInstanceNumber = attr.Value
					hana.storeDiscovery("SAPHanaInstanceNumber", hana.SAPHanaInstanceNumber)
				}
			}
		}
	}
	return nil
}

// Discover module for handling SAP Hana clusters
func NewHanaClusterDiscover(client consul.Client) HanaClusterDiscover {
	r := HanaClusterDiscover{}
	r.cluster = NewClusterDiscover(client)
	return r
}
