package cluster

import (
	"os"

	// Reusing the Prometheus Ha Exporter cibadmin xml parser here
	"github.com/ClusterLabs/ha_cluster_exporter/collector/pacemaker/cib"
	"github.com/ClusterLabs/ha_cluster_exporter/collector/pacemaker/crmmon"
)

const (
	cibAdmPath          string = "/usr/sbin/cibadmin"
	crmmonAdmPath       string = "/usr/sbin/crm_mon"
	clusterNameProperty string = "cib-bootstrap-options-cluster-name"
)

type Cluster struct {
	Cib    cib.Root    `mapstructure:"cib,omitempty"`
	Crmmon crmmon.Root `mapstructure:"crmmon,omitempty"`
}

func NewCluster() (Cluster, error) {
	var cluster = Cluster{}

	cibParser := cib.NewCibAdminParser(cibAdmPath)

	cibConfig, err := cibParser.Parse()
	if err != nil {
		return cluster, err
	}

	cluster.Cib = cibConfig

	crmmonParser := crmmon.NewCrmMonParser(crmmonAdmPath)

	crmmonConfig, err := crmmonParser.Parse()
	if err != nil {
		return cluster, err
	}

	cluster.Crmmon = crmmonConfig

	return cluster, nil
}

func (c *Cluster) Name() string {
	// Handle not named clusters
	for _, prop := range c.Cib.Configuration.CrmConfig.ClusterProperties {
		if prop.Id == clusterNameProperty {
			return prop.Value
		}
	}

	return ""
}

func (c *Cluster) IsDc() bool {
	host, _ := os.Hostname()

	for _, nodes := range c.Crmmon.Nodes {
		if nodes.Name == host {
			return nodes.DC
		}
	}

	return false
}
