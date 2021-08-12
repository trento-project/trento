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
	"github.com/trento-project/trento/internal/tags"
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

func NewSAPSystemsTable(sapSystemsList sapsystem.SAPSystemsList, hostList hosts.HostList, client consul.Client) (SAPSystemsTable, error) {
	var sapSystemsTable SAPSystemsTable
	rowsBySID := make(map[string]*SAPSystemRow)

	for _, s := range sapSystemsList {
		if s.Type != sapsystem.Application {
			continue
		}

		sapSystem, ok := rowsBySID[s.SID]
		if !ok {
			t := tags.NewTags(client, "sapsystems", s.SID)
			sapsystemTags, err := t.GetAll()
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

			for _, h := range hostList {
				if i.Host == h.Name() {
					clusterName = h.TrentoMeta()["trento-ha-cluster"]
					clusterId = h.TrentoMeta()["trento-ha-cluster-id"]
				}
			}

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

func NewSAPSystemListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sapSystemsList sapsystem.SAPSystemsList
		query := c.Request.URL.Query()
		sidFilter := query["sid"]
		tagsFilter := query["tags"]

		hostList, err := hosts.Load(client, "", nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		for _, h := range hostList {
			sapSystems, err := h.GetSAPSystems()
			if err != nil {
				_ = c.Error(err)
				return
			}

			for _, s := range sapSystems {
				sapSystemsList = append(sapSystemsList, s)
			}
		}

		sapSystemsTable, err := NewSAPSystemsTable(sapSystemsList, hostList, client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		sapSystemsTable = sapSystemsTable.filter(sidFilter, tagsFilter)

		c.HTML(http.StatusOK, "sapsystems.html.tmpl", gin.H{
			"SAPSystemsTable": sapSystemsTable,
		})
	}
}

func NewSAPSystemHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var system *sapsystem.SAPSystem
		var systemHosts hosts.HostList

		sid := c.Param("sid")

		hostList, err := hosts.Load(client, "", nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		for _, h := range hostList {
			sapSystems, err := h.GetSAPSystems()

			if err != nil {
				_ = c.Error(err)
				return
			}

			for _, s := range sapSystems {
				if s.SID == sid {
					if system == nil {
						system = s
					}
					systemHosts = append(systemHosts, h)
				}
			}
		}

		if system == nil {
			_ = c.Error(NotFoundError("could not find system"))
			return
		}
		c.HTML(http.StatusOK, "sapsystem.html.tmpl", gin.H{
			"SAPSystem": system,
			"Hosts":     systemHosts,
		})
	}
}
