package cluster

import (
	"os"
	"testing"

	"github.com/ClusterLabs/ha_cluster_exporter/collector/pacemaker/cib"
	"github.com/ClusterLabs/ha_cluster_exporter/collector/pacemaker/crmmon"

	"github.com/stretchr/testify/assert"
)

func TestClusterName(t *testing.T) {
	root := new(cib.Root)

	crmConfig := struct {
		ClusterProperties []cib.Attribute `xml:"cluster_property_set>nvpair"`
	}{
		ClusterProperties: []cib.Attribute{
			cib.Attribute{
				Id:    "cib-bootstrap-options-cluster-name",
				Value: "cluster_name",
			},
		},
	}

	root.Configuration.CrmConfig = crmConfig

	c := Cluster{
		Cib: *root,
		Crmmon: crmmon.Root{
			Version: "1.2.3",
			Nodes: []crmmon.Node{
				crmmon.Node{
					Name: "othernode",
					DC:   false,
				},
				crmmon.Node{
					Name: "yetanothernode",
					DC:   true,
				},
			},
		},
	}

	assert.Equal(t, "cluster_name", c.Name())
}

func TestIsDC(t *testing.T) {
	host, _ := os.Hostname()
	root := new(cib.Root)

	crmConfig := struct {
		ClusterProperties []cib.Attribute `xml:"cluster_property_set>nvpair"`
	}{
		ClusterProperties: []cib.Attribute{
			cib.Attribute{
				Id:    "cib-bootstrap-options-cluster-name",
				Value: "cluster_name",
			},
		},
	}

	root.Configuration.CrmConfig = crmConfig

	c := Cluster{
		Cib: *root,
		Crmmon: crmmon.Root{
			Version: "1.2.3",
			Nodes: []crmmon.Node{
				crmmon.Node{
					Name: "othernode",
					DC:   false,
				},
				crmmon.Node{
					Name: host,
					DC:   true,
				},
			},
		},
	}

	assert.Equal(t, true, c.IsDc())

	c = Cluster{
		Cib: *root,
		Crmmon: crmmon.Root{
			Version: "1.2.3",
			Nodes: []crmmon.Node{
				crmmon.Node{
					Name: "othernode",
					DC:   true,
				},
				crmmon.Node{
					Name: host,
					DC:   false,
				},
			},
		},
	}

	assert.Equal(t, false, c.IsDc())
}
