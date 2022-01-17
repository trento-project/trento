package cluster

import (
	"os"
	"strconv"
	"strings"

	"github.com/trento-project/trento/internal"

	// These packages were originally imported from github.com/ClusterLabs/ha_cluster_exporter/collector/pacemaker
	// Now we mantain our own fork
	"github.com/trento-project/trento/internal/cluster/cib"
	"github.com/trento-project/trento/internal/cluster/crmmon"
)

const (
	cibAdmPath             string = "/usr/sbin/cibadmin"
	crmmonAdmPath          string = "/usr/sbin/crm_mon"
	corosyncKeyPath        string = "/etc/corosync/authkey"
	clusterNameProperty    string = "cib-bootstrap-options-cluster-name"
	stonithEnabled         string = "cib-bootstrap-options-stonith-enabled"
	stonithResourceMissing string = "notconfigured"
	stonithAgent           string = "stonith:"
	sbdFencingAgentName    string = "external/sbd"
	clusterNameWordCount   int    = 1
)

type DiscoveryTools struct {
	CibAdmPath      string
	CrmmonAdmPath   string
	CorosyncKeyPath string
	SBDPath         string
	SBDConfigPath   string
}

type Cluster struct {
	Cib    cib.Root    `mapstructure:"cib,omitempty"`
	Crmmon crmmon.Root `mapstructure:"crmmon,omitempty"`
	SBD    SBD         `mapstructure:"sbd,omitempty"`
	Id     string      `mapstructure:"id"`
	Name   string      `mapstructure:"name"`
	DC     bool        `mapstructure:"dc"`
}

func NewCluster() (Cluster, error) {
	return NewClusterWithDiscoveryTools(&DiscoveryTools{
		CibAdmPath:      cibAdmPath,
		CrmmonAdmPath:   crmmonAdmPath,
		CorosyncKeyPath: corosyncKeyPath,
		SBDPath:         SBDPath,
		SBDConfigPath:   SBDConfigPath,
	})
}

func NewClusterWithDiscoveryTools(discoveryTools *DiscoveryTools) (Cluster, error) {
	var cluster = Cluster{}

	cibParser := cib.NewCibAdminParser(discoveryTools.CibAdmPath)

	cibConfig, err := cibParser.Parse()
	if err != nil {
		return cluster, err
	}

	cluster.Cib = cibConfig

	crmmonParser := crmmon.NewCrmMonParser(discoveryTools.CrmmonAdmPath)

	crmmonConfig, err := crmmonParser.Parse()
	if err != nil {
		return cluster, err
	}

	cluster.Crmmon = crmmonConfig

	// Set MD5-hashed key based on the corosync auth key
	cluster.Id, err = getCorosyncAuthkeyMd5(discoveryTools.CorosyncKeyPath)
	if err != nil {
		return cluster, err
	}

	cluster.Name = getName(cluster)

	if cluster.IsFencingSBD() {
		sbdData, err := NewSBD(cluster.Id, discoveryTools.SBDPath, discoveryTools.SBDConfigPath)
		if err != nil {
			return cluster, err
		}

		cluster.SBD = sbdData
	}

	cluster.DC = isDC(&cluster)

	return cluster, nil
}

func getCorosyncAuthkeyMd5(corosyncKeyPath string) (string, error) {
	kp, err := internal.Md5sumFile(corosyncKeyPath)
	return kp, err
}

func getName(c Cluster) string {
	// Handle not named clusters
	for _, prop := range c.Cib.Configuration.CrmConfig.ClusterProperties {
		if prop.Id == clusterNameProperty {
			return prop.Value
		}
	}

	return ""
}

func isDC(c *Cluster) bool {
	host, _ := os.Hostname()

	for _, nodes := range c.Crmmon.Nodes {
		if nodes.Name == host {
			return nodes.DC
		}
	}

	return false
}

func (c *Cluster) IsFencingEnabled() bool {
	for _, prop := range c.Cib.Configuration.CrmConfig.ClusterProperties {
		if prop.Id == stonithEnabled {
			b, err := strconv.ParseBool(prop.Value)
			if err != nil {
				return false
			}
			return b
		}
	}

	return false
}

func (c *Cluster) FencingResourceExists() bool {
	f := c.FencingType()

	return f != stonithResourceMissing
}

func (c *Cluster) FencingType() string {
	for _, resource := range c.Crmmon.Resources {
		if strings.HasPrefix(resource.Agent, stonithAgent) {
			return strings.Split(resource.Agent, ":")[1]
		}
	}
	return stonithResourceMissing
}

func (c *Cluster) IsFencingSBD() bool {
	f := c.FencingType()

	return f == sbdFencingAgentName
}
