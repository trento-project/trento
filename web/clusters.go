package web

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/trento-project/trento/internal"

	consulApi "github.com/hashicorp/consul/api"

	"github.com/gin-gonic/gin"

	"github.com/trento-project/trento/internal/cluster"
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

// HANARole parses the hana_prd_roles string and returns the HANA Role
// Possible values: master, slave
// e.g. 4:P:master1:master:worker:master returns master (last element)
func (node *Node) HANARole() string {
	if r, ok := node.Attributes["hana_prd_roles"]; ok {
		role := r[strings.LastIndex(r, ":")+1:]
		return strings.Title(role)
	}
	return "-"
}

// HANAHealthState parses the hana_prd_roles string and returns the SAPHanaSR Health state
// Possible values: 0-4
// 4 - SAP HANA database is up and OK. The cluster does interpret this as a correctly running database.
// 3 - SAP HANA database is up and in status info. The cluster does interpret this as a correctly running database.
// 2 - SAP HANA database is up and in status warning. The cluster does interpret this as a correctly running database.
// 1 - SAP HANA database is down. If the database should be up and is not down by intention, this could trigger a takeover.
// 0 - Internal Script Error – to be ignored.
// e.g. 4:P:master1:master:worker:master returns 4 (first element)
func (node *Node) HANAHealthState() string {
	if r, ok := node.Attributes["hana_prd_roles"]; ok {
		healthState := strings.SplitN(r, ":", 2)[0]
		return healthState
	}
	return "-"
}

// HANAStatus parses the hana_prd_roles string and returns the SAPHanaSR Health state
// Possible values: Primary, Secondary
// e.g. 4:P:master1:master:worker:master returns Primary (second element)
func (node *Node) HANAStatus() string {
	if r, ok := node.Attributes["hana_prd_roles"]; ok {
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

func NewNodes(c *cluster.Cluster, hl hosts.HostList) Nodes {
	// TODO: this factory is HANA specific,
	// eventually we will need to have different factory methods depending on the cluster type

	var nodes Nodes

	// TODO: remove plain resources grouping as in the future we'll need to distinguish between Cloned and Groups
	resources := c.Crmmon.Resources
	for _, g := range c.Crmmon.Groups {
		resources = append(resources, g.Resources...)
	}

	for _, c := range c.Crmmon.Clones {
		resources = append(resources, c.Resources...)
	}

	for _, n := range c.Crmmon.NodeAttributes.Nodes {
		node := &Node{Name: n.Name, Attributes: make(map[string]string)}

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

				if r.Agent == "ocf::heartbeat:IPaddr2" {
					for _, p := range c.Cib.Configuration.Resources.Primitives {
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
				node.Health = h.Health()
			}
		}

		nodes = append(nodes, node)
	}

	return nodes
}

func (nodes Nodes) HANASecondarySyncState() string {
	for _, n := range nodes {
		if n.HANAStatus() == "Secondary" {
			if s, ok := n.Attributes["hana_prd_sync_state"]; ok {
				return s
			}
		}
	}
	return "-"
}

func (nodes Nodes) GroupBySite() map[string]Nodes {
	nodesBySite := make(map[string]Nodes)

	for _, n := range nodes {
		if site, ok := n.Attributes["hana_prd_site"]; ok {
			nodesBySite[site] = append(nodesBySite[site], n)
		}
	}

	return nodesBySite
}

func (nodes Nodes) CriticalCount() int {
	var critical int

	for _, n := range nodes {
		if n.Health == consulApi.HealthCritical {
			critical += 1
		}
	}

	return critical
}

func (nodes Nodes) WarningCount() int {
	var warning int

	for _, n := range nodes {
		if n.Health == consulApi.HealthWarning {
			warning += 1
		}
	}

	return warning
}

func (nodes Nodes) PassingCount() int {
	var warning int

	for _, n := range nodes {
		if n.Health == consulApi.HealthPassing {
			warning += 1
		}
	}

	return warning
}

type ClustersRow struct {
	Id              string
	Name            string
	Health          string
	SIDs            []string
	Type            string
	ResourcesNumber int
	HostsNumber     int
}

type ClustersTable []*ClustersRow

func NewClustersTable(clusters map[string]*cluster.Cluster) ClustersTable {
	var clusterTable ClustersTable

	for id, c := range clusters {
		var health string
		// TODO: Cost-optimized has multiple SIDs
		var sids []string

		sids = append(sids, getHanaSID(c))

		clustersRow := &ClustersRow{
			Id:              id,
			Name:            c.Name,
			Health:          health,
			SIDs:            sids,
			Type:            detectClusterType(c),
			ResourcesNumber: c.Crmmon.Summary.Resources.Number,
			HostsNumber:     c.Crmmon.Summary.Nodes.Number,
		}
		clusterTable = append(clusterTable, clustersRow)
	}

	return clusterTable
}

func (t ClustersTable) Filter(name []string, health []string, sid []string, clusterType []string) ClustersTable {
	var filteredClustersTable ClustersTable

	for _, r := range t {
		nameFilter := len(name) == 0 || internal.Contains(name, r.Name)
		healthFilter := len(health) == 0 || internal.Contains(health, r.Health)

		sidFilter := false
		if len(sid) == 0 {
			sidFilter = true
		} else {
			for _, s := range sid {
				if internal.Contains(r.SIDs, s) {
					sidFilter = true
					break
				}
			}
		}

		clusterTypeFilter := len(clusterType) == 0 || internal.Contains(clusterType, r.Type)

		if nameFilter && healthFilter && sidFilter && clusterTypeFilter {
			filteredClustersTable = append(filteredClustersTable, r)
		}
	}
	return filteredClustersTable
}

func (t ClustersTable) GetAllSIDs() []string {
	var sids []string
	set := make(map[string]struct{})

	for _, r := range t {
		for _, s := range r.SIDs {
			_, ok := set[s]
			if !ok {
				set[s] = struct{}{}
				sids = append(sids, s)
			}
		}
	}

	return sids
}

func (t ClustersTable) GetAllClusterTypes() []string {
	var clusterTypes []string
	set := make(map[string]struct{})

	for _, r := range t {

		_, ok := set[r.Type]
		if !ok {
			set[r.Type] = struct{}{}
			clusterTypes = append(clusterTypes, r.Type)
		}

	}

	return clusterTypes
}

func NewClustersHealthContainer(t ClustersTable) *HealthContainer {
	h := &HealthContainer{}
	for _, r := range t {
		switch r.Health {
		case consulApi.HealthPassing:
			h.PassingCount += 1
		case consulApi.HealthWarning:
			h.WarningCount += 1
		case consulApi.HealthCritical:
			h.CriticalCount += 1
		}
	}
	return h
}

func NewClusterListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()

		healthFilter := query["health"]
		sidFilter := query["sid"]
		nameFilter := query["name"]
		clusterTypeFilter := query["type"]

		clusters, err := cluster.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		clustersTable := NewClustersTable(clusters)
		clustersTable = clustersTable.Filter(nameFilter, healthFilter, sidFilter, clusterTypeFilter)

		healthContainer := NewClustersHealthContainer(clustersTable)
		healthContainer.Layout = "horizontal"

		page := c.DefaultQuery("page", "1")
		perPage := c.DefaultQuery("per_page", "10")
		pagination := NewPaginationWithStrings(len(clustersTable), page, perPage)
		firstElem, lastElem := pagination.GetSliceNumbers()

		c.HTML(http.StatusOK, "clusters.html.tmpl", gin.H{
			"ClustersTable":   clustersTable[firstElem:lastElem],
			"AppliedFilters":  query,
			"Pagination":      pagination,
			"HealthContainer": healthContainer,
		})
	}
}

func getChecksCatalogModalData(clusterId string, client consul.Client, s services.ChecksService) (map[string]map[string]*models.Check, error) {
	checksCatalog, err := s.GetChecksCatalogByGroup()
	if err != nil {
		return checksCatalog, err
	}

	selectedChecks, err := cluster.GetCheckSelection(client, clusterId)
	if err != nil {
		return checksCatalog, err
	}

	for _, checkList := range checksCatalog {
		for _, check := range checkList {
			if internal.Contains(strings.Split(selectedChecks, ","), check.ID) {
				check.Selected = true
			}
		}
	}

	return checksCatalog, nil
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

		checksCatalog, errCatalog := s.GetChecksCatalog()
		checksCatalogModalData, errCatalogByGroup := getChecksCatalogModalData(clusterId, client, s)
		checksRestult, errResult := s.GetChecksResultByCluster(clusterItem.Name)
		if errCatalog != nil || errCatalogByGroup != nil || errResult != nil {
			StoreAlert(c, AlertCatalogNotFound())
		}

		nodes := NewNodes(clusterItem, hosts)

		hContainer := &HealthContainer{
			CriticalCount: nodes.CriticalCount(),
			WarningCount:  nodes.WarningCount(),
			PassingCount:  nodes.PassingCount(),
			Layout:        "vertical",
		}

		c.HTML(http.StatusOK, "cluster_hana.html.tmpl", gin.H{
			"Cluster":            clusterItem,
			"Nodes":              nodes,
			"StoppedResources":   stoppedResources(clusterItem),
			"ClusterType":        clusterType,
			"HealthContainer":    hContainer,
			"ChecksCatalog":      checksCatalog,
			"ChecksCatalogModal": checksCatalogModalData,
			"ChecksResult":       checksRestult,
			"Alerts":             GetAlerts(c),
		})
	}
}

type checkSelectionForm struct {
	Ids []string `form:"ids[]"`
}

func NewSaveChecksHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterId := c.Param("id")

		var selectedChecks checkSelectionForm
		c.ShouldBind(&selectedChecks)

		err := cluster.StoreCheckSelection(
			client, clusterId, strings.Join(selectedChecks.Ids, ","))

		var newAlert *Alert
		if err == nil {
			newAlert = &Alert{
				Type:  "success",
				Title: "Check selection saved",
				Text:  "The cluster checks selection has been saved correctly.",
			}
		} else {
			newAlert = &Alert{
				Type:  "danger",
				Title: "Error saving check selection",
				Text:  "The cluster checks selection couldn't be saved.",
			}
		}
		StoreAlert(c, newAlert)

		c.Redirect(http.StatusFound, path.Join("/clusters", clusterId))
	}
}
