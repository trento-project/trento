package datapipeline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/internal/cluster/cib"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewClustersProjector(db *gorm.DB) *projector {
	clusterProjector := NewProjector("clusters", db)
	clusterProjector.AddHandler(ClusterDiscovery, clustersProjector_ClusterDiscoveryHandler)

	return clusterProjector
}

// TODO: this is a temporary solution, this code needs to be abstracted in the projector.Project() method
func clustersProjector_ClusterDiscoveryHandler(event *DataCollectedEvent, db *gorm.DB) error {
	data, _ := event.Payload.MarshalJSON()
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	var cluster cluster.Cluster
	if err := dec.Decode(&cluster); err != nil {
		log.Errorf("can't decode data: %s", err)
		return err
	}

	clusterListReadModel, err := transformClusterData(&cluster)
	if err != nil {
		log.Errorf("can't transform data: %s", err)
		return err
	}

	return db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(clusterListReadModel).Error
}

// transformClusterData transforms the cluster data into the read model
func transformClusterData(cluster *cluster.Cluster) (*entities.Cluster, error) {
	if cluster.Id == "" {
		return nil, fmt.Errorf("no cluster ID found")
	}

	clusterDetail, _ := parseClusterDetails(cluster)
	log.Debugf("%s", clusterDetail)

	return &entities.Cluster{
		ID:          cluster.Id,
		Name:        cluster.Name,
		ClusterType: detectClusterType(cluster),
		// TODO: Cost-optimized has multiple SIDs, we will need to implement this in the future
		SIDs:            getClusterSIDs(cluster),
		ResourcesNumber: cluster.Crmmon.Summary.Resources.Number,
		HostsNumber:     cluster.Crmmon.Summary.Nodes.Number,
		Details:         (datatypes.JSON)(clusterDetail),
	}, nil
}

// detectClusterType returns the cluster type based on the cluster resources
func detectClusterType(cluster *cluster.Cluster) string {
	var hasSapHanaTopology, hasSAPHanaController, hasSAPHana bool

	for _, c := range cluster.Crmmon.Clones {
		for _, r := range c.Resources {
			switch r.Agent {
			case "ocf::suse:SAPHanaTopology":
				hasSapHanaTopology = true
			case "ocf::suse:SAPHana":
				hasSAPHana = true
			case "ocf::suse:SAPHanaController":
				hasSAPHanaController = true
			}
		}
	}

	switch {
	case hasSapHanaTopology && hasSAPHana:
		return models.ClusterTypeHANAScaleUp
	case hasSapHanaTopology && hasSAPHanaController:
		return models.ClusterTypeHANAScaleOut
	default:
		return models.ClusterTypeUnknown
	}
}

// getClusterSIDs returns the SIDs of the cluster
// TODO: HANA scale-out has multiple SIDs, we will need to implement this in the future
func getClusterSIDs(c *cluster.Cluster) []string {
	var sids []string
	for _, r := range c.Cib.Configuration.Resources.Clones {
		if r.Primitive.Type == "SAPHanaTopology" {
			for _, a := range r.Primitive.InstanceAttributes {
				if a.Name == "SID" && a.Value != "" {
					sids = append(sids, a.Value)
				}
			}
		}
	}

	return sids
}

// parseClusterDetails parses the cluster details depending on the cluster type
func parseClusterDetails(c *cluster.Cluster) (json.RawMessage, error) {
	switch detectClusterType(c) {
	case models.ClusterTypeHANAScaleUp, models.ClusterTypeHANAScaleOut:
		return parseHANAClusterDetails(c)
	default:
		return json.RawMessage{}, nil
	}
}

// parseHANAClusterDetails parses the HANA cluster details
func parseHANAClusterDetails(c *cluster.Cluster) (json.RawMessage, error) {
	var sid string

	sids := getClusterSIDs(c)
	if len(sids) > 0 {
		sid = sids[0]
	}
	nodes := parseClusterNodes(c)

	var systemReplicationMode, systemReplicationOperationMode, secondarySyncState, srHealthState string
	if len(nodes) > 0 {
		systemReplicationMode, _ = parseHANAAttribute(nodes[0], "srmode", sid)
		systemReplicationOperationMode, _ = parseHANAAttribute(nodes[0], "op_mode", sid)
		secondarySyncState = parseHANASecondarySyncState(nodes, sid)
		srHealthState = parseHANAHealthState(nodes[0], sid)
	}

	dateLayout := "Mon Jan 2 15:04:05 2006"
	cibLastWritten, _ := time.Parse(dateLayout, c.Crmmon.Summary.LastChange.Time)

	clusterDetail := &entities.HANAClusterDetails{
		SystemReplicationMode:          systemReplicationMode,
		SecondarySyncState:             secondarySyncState,
		SystemReplicationOperationMode: systemReplicationOperationMode,
		SRHealthState:                  srHealthState,
		CIBLastWritten:                 cibLastWritten,
		StonithType:                    parseClusterFencingType(c),
		StoppedResources:               parseClusterStoppedResources(c),
		Nodes:                          nodes,
		SBDDevices:                     parseSBDDevices(c),
	}

	return json.Marshal(clusterDetail)
}

// parseClusterNodes parses the cluster nodes from the crmmon/cib data
func parseClusterNodes(c *cluster.Cluster) []*entities.HANAClusterNode {
	var nodes []*entities.HANAClusterNode

	// TODO: remove plain resources grouping as in the future we'll need to distinguish between Cloned and Groups
	resources := c.Crmmon.Resources
	for _, g := range c.Crmmon.Groups {
		resources = append(resources, g.Resources...)
	}

	for _, c := range c.Crmmon.Clones {
		resources = append(resources, c.Resources...)
	}

	for _, n := range c.Crmmon.NodeAttributes.Nodes {
		node := &entities.HANAClusterNode{
			Name:       n.Name,
			Attributes: make(map[string]string),
		}

		for _, a := range n.Attributes {
			node.Attributes[a.Name] = a.Value
		}

		for _, r := range resources {
			if r.Node == nil {
				continue
			}
			if r.Node.Name == n.Name {
				resource := &entities.ClusterResource{
					ID:   r.Id,
					Type: r.Agent,
					Role: r.Role,
				}

				switch {
				case r.Active:
					resource.Status = "active"
				case r.Blocked:
					resource.Status = "blocked"
				case r.Failed:
					resource.Status = "failed"
				case r.FailureIgnored:
					resource.Status = "failure_ignored"
				case r.Orphaned:
					resource.Status = "orphaned"
				}

				var primitives []cib.Primitive
				primitives = append(primitives, c.Cib.Configuration.Resources.Primitives...)

				for _, g := range c.Cib.Configuration.Resources.Groups {
					primitives = append(primitives, g.Primitives...)
				}

				if r.Agent == "ocf::heartbeat:IPaddr2" {
					for _, p := range primitives {
						if r.Id == p.Id {
							if len(p.InstanceAttributes) > 0 {
								node.VirtualIPs = append(node.VirtualIPs, p.InstanceAttributes[0].Value)
								break
							}
						}
					}
				}

				for _, nh := range c.Crmmon.NodeHistory.Nodes {
					if nh.Name == n.Name {
						for _, rh := range nh.ResourceHistory {
							if rh.Name == resource.ID {
								resource.FailCount = rh.FailCount
								break
							}
						}
					}
				}

				node.Resources = append(node.Resources, resource)
			}
		}

		// TODO: refactor, remove array of sids and use only one sid
		sids := getClusterSIDs(c)

		if len(sids) > 0 {
			node.Site, _ = parseHANAAttribute(node, "site", sids[0])
			node.HANAStatus = parseHANAStatus(node, sids[0])
		}
		nodes = append(nodes, node)
	}

	return nodes
}

// parseHANAAttribute returns an HANA attribute value
func parseHANAAttribute(node *entities.HANAClusterNode, attributeName string, sid string) (string, bool) {
	hanaAttributeName := fmt.Sprintf("hana_%s_%s", strings.ToLower(sid), attributeName)
	value, ok := node.Attributes[hanaAttributeName]

	return value, ok
}

// parseHANASecondarySyncState returns the secondary sync state of the HANA cluster
func parseHANASecondarySyncState(nodes []*entities.HANAClusterNode, sid string) string {
	for _, n := range nodes {
		if parseHANAStatus(n, sid) == models.HANAStatusSecondary {
			if s, ok := parseHANAAttribute(n, "sync_state", sid); ok {
				return s
			}
		}
	}
	return ""
}

// parseHANAStatus parses the hana_<SID>_roles string and returns the SAPHanaSR Health state
// Possible values: Primary, Secondary
// e.g. 4:P:master1:master:worker:master returns Primary (second element)
func parseHANAStatus(node *entities.HANAClusterNode, sid string) string {
	if r, ok := parseHANAAttribute(node, "roles", sid); ok {
		status := strings.SplitN(r, ":", 3)[1]

		switch status {
		case "P":
			return models.HANAStatusPrimary
		case "S":
			return models.HANAStatusSecondary
		}
	}
	return ""
}

// parseHANAHealthState returns the SAPHanaSR Health state
func parseHANAHealthState(node *entities.HANAClusterNode, sid string) string {
	if r, ok := parseHANAAttribute(node, "roles", sid); ok {
		healthState := strings.SplitN(r, ":", 2)[0]
		return healthState
	}
	return ""
}

// parseClusterFencingType returns the cluster fencing type,
// or an empty string if the fencing is not configured
func parseClusterFencingType(c *cluster.Cluster) string {
	for _, resource := range c.Crmmon.Resources {
		if strings.HasPrefix(resource.Agent, "stonith:") {
			return strings.Split(resource.Agent, ":")[1]
		}
	}

	return ""
}

// parseClusterStoppedResources returns all the stopped resources in a cluster
func parseClusterStoppedResources(c *cluster.Cluster) []*entities.ClusterResource {
	var stoppedResources []*entities.ClusterResource

	for _, r := range c.Crmmon.Resources {
		if r.NodesRunningOn == 0 && !r.Active {
			resource := &entities.ClusterResource{
				ID: r.Id,
			}
			stoppedResources = append(stoppedResources, resource)
		}
	}

	return stoppedResources
}

// parseSBDDevices returns a slice of SBD devices
func parseSBDDevices(c *cluster.Cluster) []*entities.SBDDevice {
	var sbdDevices []*entities.SBDDevice
	for _, s := range c.SBD.Devices {
		sbdDevice := &entities.SBDDevice{
			Device: s.Device,
		}
		sbdDevices = append(sbdDevices, sbdDevice)
	}

	return sbdDevices
}
