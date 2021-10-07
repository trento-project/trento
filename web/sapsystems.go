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
	SID            string
	Tags           []string
	InstancesTable []*InstanceRow
}

type InstanceRow struct {
	SID            string
	InstanceNumber string
	Features       string
	Host           string
	ClusterName    string
	ClusterId      string
}

type SAPSystemsTable []*SAPSystemRow

var systemTypeToTag = map[int]string{
	sapsystem.Application: models.TagSAPSystemResourceType,
	sapsystem.Database:    models.TagDatabaseResourceType,
}

func NewSAPSystemsTable(sapSystemsList sapsystem.SAPSystemsList, hostsService services.HostsService, tagsService services.TagsService) (SAPSystemsTable, error) {
	var sapSystemsTable SAPSystemsTable
	rowsBySID := make(map[string]*SAPSystemRow)

	for _, s := range sapSystemsList {

		sapSystem, ok := rowsBySID[s.SID]
		if !ok {
			sapsystemTags, err := tagsService.GetAllByResource(systemTypeToTag[s.Type], s.SID)
			if err != nil {
				return nil, err
			}

			sapSystem = &SAPSystemRow{
				SID:  s.SID,
				Tags: sapsystemTags,
			}
			rowsBySID[s.SID] = sapSystem
		}

		for _, i := range s.Instances {
			var features string
			var instanceNumber string
			var clusterName, clusterId string

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
				SID:            s.SID,
				InstanceNumber: instanceNumber,
				Features:       features,
				Host:           i.Host,
				ClusterName:    clusterName,
				ClusterId:      clusterId,
			}

			sapSystem.InstancesTable = append(sapSystem.InstancesTable, instance)
		}
	}

	for _, row := range rowsBySID {
		sapSystemsTable = append(sapSystemsTable, row)
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

func NewSAPSystemListHandler(client consul.Client, hostsService services.HostsService, sapSystemsService services.SAPSystemsService, tagsService services.TagsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()
		sidFilter := query["sid"]
		tagsFilter := query["tags"]

		saps, err := sapSystemsService.GetSAPSystemsByType(sapsystem.Application)
		if err != nil {
			_ = c.Error(err)
			return
		}

		sapSystemsTable, err := NewSAPSystemsTable(saps, hostsService, tagsService)
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

		sapDatabasesTable, err := NewSAPSystemsTable(saps, hostsService, tagsService)
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

		sid := c.Param("sid")

		systemList, err = sapSystemsService.GetSAPSystemsBySid(sid)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if len(systemList) == 0 {
			_ = c.Error(NotFoundError("could not find system"))
			return
		}

		systemHosts, err = hostsService.GetHostsBySid(sid)
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
