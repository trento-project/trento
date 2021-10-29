package web

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/trento-project/trento/internal"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/internal/sapsystem"
	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

type SAPSystemRow struct {
	Id                  string
	SID                 string
	AttachedDatabaseSID string
	AttachedDatabaseId  string
	Profile             map[string]interface{}
	Tags                []string
	InstancesTable      []*InstanceRow
	HasDuplicatedSid    bool
}

type InstanceRow struct {
	Type              int
	SID               string
	InstanceNumber    string
	Features          string
	SystemReplication SystemReplication // Only for Database type
	Host              string
	ClusterName       string
	ClusterId         string
}

type SystemReplication map[string]interface{}

type SAPSystemsTable []*SAPSystemRow

var systemTypeToTag = map[int]string{
	sapsystem.Application: models.TagSAPSystemResourceType,
	sapsystem.Database:    models.TagDatabaseResourceType,
}

func NewSAPSystemsTable(sapSystemsList sapsystem.SAPSystemsList, hostsService services.HostsService,
	sapSystemsService services.SAPSystemsService, tagsService services.TagsService) (SAPSystemsTable, error) {
	var sapSystemsTable SAPSystemsTable
	sids := make(map[string]int)
	rowsBySID := make(map[string]*SAPSystemRow)

	for _, s := range sapSystemsList {

		sapSystem, ok := rowsBySID[s.Id]
		if !ok {
			sapsystemTags, err := tagsService.GetAllByResource(systemTypeToTag[s.Type], s.Id)
			if err != nil {
				return nil, err
			}

			sids[s.SID] += 1

			sapSystem = &SAPSystemRow{
				Id:      s.Id,
				SID:     s.SID,
				Tags:    sapsystemTags,
				Profile: s.Profile,
			}

			if s.Type == sapsystem.Application {
				attachedDatabases, err := sapSystemsService.GetAttachedDatabasesById(s.Id)
				if err != nil {
					return nil, err
				}

				if len(attachedDatabases) == 0 {
					continue
				}

				sapSystem.AttachedDatabaseSID = attachedDatabases[0].SID
				sapSystem.AttachedDatabaseId = attachedDatabases[0].Id

				// Add the database instances to the instaces table
				for index, database := range attachedDatabases {
					for _, data := range database.Instances {
						s.Instances[fmt.Sprint(index)] = data
					}
				}
			}

			rowsBySID[s.Id] = sapSystem
		}

		for _, i := range s.Instances {
			var features string
			var sid string
			var instanceNumber string
			var clusterName, clusterId string

			if p, ok := i.SAPControl.Properties["SAPSYSTEMNAME"]; ok {
				sid = p.Value
			}

			if p, ok := i.SAPControl.Properties["SAPSYSTEM"]; ok {
				instanceNumber = p.Value
			}

			for _, ci := range i.SAPControl.Instances {
				if instanceNumber == fmt.Sprintf("%02d", ci.InstanceNr) {
					features = ci.Features
				}
			}

			metadata, _ := hostsService.GetHostMetadata(i.Host)
			clusterName = metadata["trento-ha-cluster"]
			clusterId = metadata["trento-ha-cluster-id"]

			instance := &InstanceRow{
				Type:           i.Type,
				SID:            sid,
				InstanceNumber: instanceNumber,
				Features:       features,
				Host:           i.Host,
				ClusterName:    clusterName,
				ClusterId:      clusterId,
			}

			if i.Type == sapsystem.Database {
				// Cast to local struct to manage the data in this package
				instance.SystemReplication = SystemReplication(i.SystemReplication)
			}

			sapSystem.InstancesTable = append(sapSystem.InstancesTable, instance)
		}

		sort.Slice(sapSystem.InstancesTable, func(i, j int) bool {
			return sapSystem.InstancesTable[i].SID < sapSystem.InstancesTable[j].SID
		})
	}

	for _, row := range rowsBySID {
		sapSystemsTable = append(sapSystemsTable, row)
	}

	for _, s := range sapSystemsTable {
		if sids[s.SID] > 1 {
			s.HasDuplicatedSid = true
		}
	}

	sort.Slice(sapSystemsTable, func(i, j int) bool {
		return sapSystemsTable[i].SID < sapSystemsTable[j].SID
	})

	return sapSystemsTable, nil
}

func (t SAPSystemsTable) filter(sid []string, tags []string) SAPSystemsTable {
	var filteredSAPSystemsTable SAPSystemsTable

	for _, r := range t {
		if len(sid) > 0 && !internal.Contains(sid, r.SID) {
			continue
		}

		if len(tags) > 0 {
			tagFound := false
			for _, t := range tags {
				if internal.Contains(r.Tags, t) {
					tagFound = true
					break
				}
			}

			if !tagFound {
				continue
			}
		}

		filteredSAPSystemsTable = append(filteredSAPSystemsTable, r)
	}

	return filteredSAPSystemsTable
}

func (t SAPSystemsTable) GetAllSIDs() []string {
	var sids []string
	set := make(map[string]struct{})

	for _, r := range t {
		_, ok := set[r.SID]
		if !ok {
			set[r.SID] = struct{}{}
			sids = append(sids, r.SID)
		}
	}

	return sids
}

func (t SAPSystemsTable) GetAllTags() []string {
	var tags []string
	set := make(map[string]struct{})

	for _, r := range t {
		for _, tag := range r.Tags {
			_, ok := set[tag]
			if !ok {
				set[tag] = struct{}{}
				tags = append(tags, tag)
			}
		}
	}

	return tags
}

func (r InstanceRow) IsDatabase() bool {
	return bool(r.Type == sapsystem.Database)
}

// Find mode information at: https://help.sap.com/viewer/4e9b18c116aa42fc84c7dbfd02111aba/2.0.04/en-US/aefc55a27003440792e34ece2125dc89.html
func (s SystemReplication) GetReplicationMode() string {
	localSite, ok := s["local_site_id"]
	if !ok {
		return ""
	}

	var mode string

	site, ok := s["site"]
	if !ok {
		return ""
	}

	for siteId, site := range site.(map[string]interface{}) {
		if siteId == localSite {
			mode = fmt.Sprintf("%v", site.(map[string]interface{})["REPLICATION_MODE"])
			break
		}
	}

	switch mode {
	case "PRIMARY":
		return "Primary"
	case "":
		return ""
	default: // SYNC, SYNCMEM, ASYNC, UNKNOWN
		return "Secondary"
	}
}

// Find status information at: https://help.sap.com/viewer/4e9b18c116aa42fc84c7dbfd02111aba/2.0.04/en-US/aefc55a27003440792e34ece2125dc89.html
func (s SystemReplication) GetReplicationStatus() string {
	status, ok := s["overall_replication_status"]
	if !ok {
		return ""
	}

	status = fmt.Sprintf("%v", status)

	switch status {
	case "ACTIVE":
		return "SOK"
	case "ERROR":
		return "SFAIL"
	default: // UNKNOWN, INITIALIZING, SYNCING
		return ""
	}
}

func NewSAPSystemListHandler(client consul.Client, hostsService services.HostsService,
	sapSystemsService services.SAPSystemsService, tagsService services.TagsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()
		sidFilter := query["sid"]
		tagsFilter := query["tags"]

		saps, err := sapSystemsService.GetSAPSystemsByType(sapsystem.Application)
		if err != nil {
			_ = c.Error(err)
			return
		}

		sapSystemsTable, err := NewSAPSystemsTable(saps, hostsService, sapSystemsService, tagsService)
		if err != nil {
			_ = c.Error(err)
			return
		}

		sapSystemsTable = sapSystemsTable.filter(sidFilter, tagsFilter)

		c.HTML(http.StatusOK, "sapsystems.html.tmpl", gin.H{
			"Title":           "SAP Systems",
			"ResourcePath":    "sapsystems",
			"SAPSystemsTable": sapSystemsTable,
		})
	}
}

func NewHanaDatabaseListHandler(
	client consul.Client, hostsService services.HostsService,
	sapSystemsService services.SAPSystemsService, tagsService services.TagsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()
		sidFilter := query["sid"]
		tagsFilter := query["tags"]

		saps, err := sapSystemsService.GetSAPSystemsByType(sapsystem.Database)
		if err != nil {
			_ = c.Error(err)
			return
		}

		sapDatabasesTable, err := NewSAPSystemsTable(saps, hostsService, sapSystemsService, tagsService)
		if err != nil {
			_ = c.Error(err)
			return
		}

		sapDatabasesTable = sapDatabasesTable.filter(sidFilter, tagsFilter)

		c.HTML(http.StatusOK, "sapsystems.html.tmpl", gin.H{
			"Title":           "HANA Databases",
			"ResourcePath":    "databases",
			"SAPSystemsTable": sapDatabasesTable,
		})
	}
}

func NewSAPResourceHandler(hostsService services.HostsService, sapSystemsService services.SAPSystemsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var systemList sapsystem.SAPSystemsList
		var systemHosts hosts.HostList
		var err error

		id := c.Param("id")

		systemList, err = sapSystemsService.GetSAPSystemsById(id)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if len(systemList) == 0 {
			_ = c.Error(NotFoundError("could not find system"))
			return
		}

		systemHosts, err = hostsService.GetHostsBySystemId(id)
		if err != nil {
			_ = c.Error(err)
			return
		}

		// We will send the 1st entry by now, as only use the layout, which is repeated among all the
		// SAP instances within a System. It does not resolve the HANA SR scenario in any case
		c.HTML(http.StatusOK, "sapsystem.html.tmpl", gin.H{
			"SAPSystem": systemList[0],
			"Hosts":     systemHosts,
		})
	}
}
