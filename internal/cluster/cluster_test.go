package cluster

import (
	"os"
	"testing"

	"github.com/ClusterLabs/ha_cluster_exporter/collector/pacemaker/cib"
	"github.com/ClusterLabs/ha_cluster_exporter/collector/pacemaker/crmmon"

	"github.com/stretchr/testify/assert"
)

func TestClusterId(t *testing.T) {
	root := new(cib.Root)

	c := Cluster{
		Cib:  *root,
		Name: "sculpin",
		Id:   "47d1190ffb4f781974c8356d7f863b03",
	}

	authkey, _ := getCorosyncAuthkeyMd5("../../test/authkey")

	assert.Equal(t, c.Id, authkey)
}

func TestClusterAlias(t *testing.T) {
	root := new(cib.Root)

	c := Cluster{
		Cib:  *root,
		Name: "sculpin",
		Id:   "47d1190ffb4f781974c8356d7f863b03",
	}

	name, _ := getName(c.Id)

	assert.Equal(t, c.Name, name)
}

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
		Name: "cluster_name",
	}

	assert.Equal(t, "cluster_name", c.Name)
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

func TestIsFencingEnabled(t *testing.T) {
	root := new(cib.Root)

	crmConfig := struct {
		ClusterProperties []cib.Attribute `xml:"cluster_property_set>nvpair"`
	}{
		ClusterProperties: []cib.Attribute{
			cib.Attribute{
				Id:    "cib-bootstrap-options-stonith-enabled",
				Value: "true",
			},
		},
	}

	root.Configuration.CrmConfig = crmConfig

	c := Cluster{
		Cib: *root,
	}

	assert.Equal(t, true, c.IsFencingEnabled())

	crmConfig = struct {
		ClusterProperties []cib.Attribute `xml:"cluster_property_set>nvpair"`
	}{
		ClusterProperties: []cib.Attribute{
			cib.Attribute{
				Id:    "cib-bootstrap-options-stonith-enabled",
				Value: "false",
			},
		},
	}

	root.Configuration.CrmConfig = crmConfig

	c = Cluster{
		Cib: *root,
	}

	assert.Equal(t, false, c.IsFencingEnabled())
}

func TestFencingType(t *testing.T) {
	c := Cluster{
		Crmmon: crmmon.Root{
			Version: "1.2.3",
			Resources: []crmmon.Resource{
				crmmon.Resource{
					Agent: "stonith:myfencing",
				},
			},
		},
	}

	assert.Equal(t, "myfencing", c.FencingType())

	c = Cluster{
		Crmmon: crmmon.Root{
			Version: "1.2.3",
			Resources: []crmmon.Resource{
				crmmon.Resource{
					Agent: "notstonith:myfencing",
				},
			},
		},
	}

	assert.Equal(t, "notconfigured", c.FencingType())
}

func TestFencingResourceExists(t *testing.T) {
	c := Cluster{
		Crmmon: crmmon.Root{
			Version: "1.2.3",
			Resources: []crmmon.Resource{
				crmmon.Resource{
					Agent: "stonith:myfencing",
				},
			},
		},
	}

	assert.Equal(t, true, c.FencingResourceExists())

	c = Cluster{
		Crmmon: crmmon.Root{
			Version: "1.2.3",
			Resources: []crmmon.Resource{
				crmmon.Resource{
					Agent: "notstonith:myfencing",
				},
			},
		},
	}

	assert.Equal(t, false, c.FencingResourceExists())
}

func TestIsFencingSBD(t *testing.T) {
	c := Cluster{
		Crmmon: crmmon.Root{
			Version: "1.2.3",
			Resources: []crmmon.Resource{
				crmmon.Resource{
					Agent: "stonith:external/sbd",
				},
			},
		},
	}

	assert.Equal(t, true, c.IsFencingSBD())

	c = Cluster{
		Crmmon: crmmon.Root{
			Version: "1.2.3",
			Resources: []crmmon.Resource{
				crmmon.Resource{
					Agent: "stonith:other",
				},
			},
		},
	}

	assert.Equal(t, false, c.IsFencingSBD())
}
