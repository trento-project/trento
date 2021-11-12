package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/internal/cluster/cib"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

type Node struct {
	Name       string
	Attributes map[string]string
	Resources  []*Resource
	Ip         string
	Health     string
	VirtualIps []string
	SID        string
}

type Resource struct {
	Id        string
	Type      string
	Role      string
	Status    string
	FailCount int
}

type Nodes []*Node

const (
	ClusterTypeScaleUp  = "HANA scale-up"
	ClusterTypeScaleOut = "HANA scale-out"
	ClusterTypeUnknown  = "Unknown"
)

func detectClusterType(c *cluster.Cluster) string {
	var hasSapHanaTopology, hasSAPHanaController, hasSAPHana bool

	for _, c := range c.Crmmon.Clones {
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
		return ClusterTypeScaleUp
	case hasSapHanaTopology && hasSAPHanaController:
		return ClusterTypeScaleOut
	default:
		return ClusterTypeUnknown
	}
}

func stoppedResources(c *cluster.Cluster) []*Resource {
	var stoppedResources []*Resource

	for _, r := range c.Crmmon.Resources {
		if r.NodesRunningOn == 0 && !r.Active {
			resource := &Resource{
				Id: r.Id,
			}
			stoppedResources = append(stoppedResources, resource)
		}
	}

	return stoppedResources
}

func getHanaSID(c *cluster.Cluster) string {
	for _, r := range c.Cib.Configuration.Resources.Clones {
		if r.Primitive.Type == "SAPHanaTopology" {
			for _, a := range r.Primitive.InstanceAttributes {
				if a.Name == "SID" {
					return a.Value
				}
			}
		}
	}

	return ""
}

func (node *Node) GetHanaAttribute(attributeName string) (string, bool) {
	hanaAttributeName := fmt.Sprintf("hana_%s_%s", strings.ToLower(node.SID), attributeName)
	value, ok := node.Attributes[hanaAttributeName]

	return value, ok
}

// HANAHealthState parses the hana_<SID>_roles string and returns the SAPHanaSR Health state
// Possible values: 0-4
// 4 - SAP HANA database is up and OK. The cluster does interpret this as a correctly running database.
// 3 - SAP HANA database is up and in status info. The cluster does interpret this as a correctly running database.
// 2 - SAP HANA database is up and in status warning. The cluster does interpret this as a correctly running database.
// 1 - SAP HANA database is down. If the database should be up and is not down by intention, this could trigger a takeover.
// 0 - Internal Script Error – to be ignored.
// e.g. 4:P:master1:master:worker:master returns 4 (first element)
func (node *Node) HANAHealthState() string {
	if r, ok := node.GetHanaAttribute("roles"); ok {
		healthState := strings.SplitN(r, ":", 2)[0]
		return healthState
	}
	return "-"
}

// HANAStatus parses the hana_<SID>_roles string and returns the SAPHanaSR Health state
// Possible values: Primary, Secondary
// e.g. 4:P:master1:master:worker:master returns Primary (second element)
func (node *Node) HANAStatus() string {
	if r, ok := node.GetHanaAttribute("roles"); ok {
		status := strings.SplitN(r, ":", 3)[1]

		switch status {
		case "P":
			return "Primary"
		case "S":
			return "Secondary"
		}
	}
	return "-"
}

func NewNodes(s services.ChecksService, c *cluster.Cluster, hl hosts.HostList) Nodes {
	// TODO: this factory is HANA specific,
	// eventually we will need to have different factory methods depending on the cluster type

	var nodes Nodes
	sid := getHanaSID(c)

	// TODO: remove plain resources grouping as in the future we'll need to distinguish between Cloned and Groups
	resources := c.Crmmon.Resources
	for _, g := range c.Crmmon.Groups {
		resources = append(resources, g.Resources...)
	}

	for _, c := range c.Crmmon.Clones {
		resources = append(resources, c.Resources...)
	}

	for _, n := range c.Crmmon.NodeAttributes.Nodes {
		node := &Node{
			Name:       n.Name,
			Attributes: make(map[string]string),
			SID:        sid,
		}

		for _, a := range n.Attributes {
			node.Attributes[a.Name] = a.Value
		}

		for _, r := range resources {
			if r.Node.Name == n.Name {
				resource := &Resource{
					Id:   r.Id,
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
							node.VirtualIps = append(node.VirtualIps, p.InstanceAttributes[0].Value)
							break
						}
					}
				}

				for _, nh := range c.Crmmon.NodeHistory.Nodes {
					if nh.Name == n.Name {
						for _, rh := range nh.ResourceHistory {
							if rh.Name == resource.Id {
								resource.FailCount = rh.FailCount
								break
							}
						}
					}
				}

				node.Resources = append(node.Resources, resource)
			}
		}

		for _, h := range hl {
			if h.Name() == node.Name {
				node.Ip = h.Address
				cData, _ := s.GetAggregatedChecksResultByHost(c.Id)
				if _, ok := cData[node.Name]; ok {
					node.Health = cData[node.Name].String()
				}
			}
		}

		nodes = append(nodes, node)
	}

	return nodes
}

func (nodes Nodes) HANASecondarySyncState() string {
	for _, n := range nodes {
		if n.HANAStatus() == "Secondary" {
			if s, ok := n.GetHanaAttribute("sync_state"); ok {
				return s
			}
		}
	}
	return "-"
}

func (nodes Nodes) HANASystemReplicationMode() string {
	if len(nodes) > 0 {
		if srmode, ok := nodes[0].GetHanaAttribute("srmode"); ok {
			return srmode
		}
	}
	return "-"
}

func (nodes Nodes) HANASystemReplicationOperationMode() string {
	if len(nodes) > 0 {
		if srmode, ok := nodes[0].GetHanaAttribute("op_mode"); ok {
			return srmode
		}
	}
	return "-"
}

func (nodes Nodes) GroupBySite() map[string]Nodes {
	nodesBySite := make(map[string]Nodes)

	for _, n := range nodes {
		if site, ok := n.GetHanaAttribute("site"); ok {
			nodesBySite[site] = append(nodesBySite[site], n)
		}
	}

	return nodesBySite
}

func NewClustersHealthContainer(clusterList models.ClusterList) *HealthContainer {
	h := &HealthContainer{}
	for _, c := range clusterList {
		switch c.Health {
		case models.CheckPassing:
			h.PassingCount += 1
		case models.CheckWarning:
			h.WarningCount += 1
		case models.CheckCritical:
			h.CriticalCount += 1
		}
	}
	return h
}

func NewClusterListHandler(clustersService services.ClustersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()

		clusterList, err := clustersService.GetAll(query)
		if err != nil {
			_ = c.Error(err)
			return
		}

		filterClusterTypes, err := clustersService.GetAllClusterTypes()
		if err != nil {
			_ = c.Error(err)
			return
		}

		filterSIDs, err := clustersService.GetAllSIDs()
		if err != nil {
			_ = c.Error(err)
			return
		}

		filterTags, err := clustersService.GetAllTags()
		if err != nil {
			_ = c.Error(err)
			return
		}

		healthContainer := NewClustersHealthContainer(clusterList)
		healthContainer.Layout = "horizontal"

		page := c.DefaultQuery("page", "1")
		perPage := c.DefaultQuery("per_page", "10")
		pagination := NewPaginationWithStrings(len(clusterList), page, perPage) // not nice here
		firstElem, lastElem := pagination.GetSliceNumbers()

		c.HTML(http.StatusOK, "clusters.html.tmpl", gin.H{
			"ClustersTable":      clusterList[firstElem:lastElem],
			"AppliedFilters":     query,
			"FilterClusterTypes": filterClusterTypes,
			"FilterSIDs":         filterSIDs,
			"FilterTags":         filterTags,
			"Pagination":         pagination,
			"HealthContainer":    healthContainer,
		})
	}
}

func NewClusterHandler(client consul.Client, s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterId := c.Param("id")

		clusters, err := cluster.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		clusterItem, ok := clusters[clusterId]
		if !ok {
			_ = c.Error(NotFoundError("could not find cluster"))
			return
		}

		filterQuery := fmt.Sprintf("Meta[\"trento-ha-cluster-id\"] == \"%s\"", clusterId)
		hosts, err := hosts.Load(client, filterQuery, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		clusterType := detectClusterType(clusterItem)
		if clusterType == ClusterTypeUnknown {
			c.HTML(http.StatusOK, "cluster_generic.html.tmpl", gin.H{
				"Cluster": clusterItem,
				"Hosts":   hosts,
			})
			return
		}

		nodes := NewNodes(s, clusterItem, hosts)
		// It returns an empty aggretaged data in case of error
		aCheckData, _ := s.GetAggregatedChecksResultByCluster(clusterId)

		hContainer := &HealthContainer{
			CriticalCount: aCheckData.CriticalCount,
			WarningCount:  aCheckData.WarningCount,
			PassingCount:  aCheckData.PassingCount,
			Layout:        "vertical",
		}

		c.HTML(http.StatusOK, "cluster_hana.html.tmpl", gin.H{
			"Cluster":          clusterItem,
			"SID":              getHanaSID(clusterItem),
			"Nodes":            nodes,
			"StoppedResources": stoppedResources(clusterItem),
			"ClusterType":      clusterType,
			"HealthContainer":  hContainer,
			"Alerts":           GetAlerts(c),
		})
	}
}
