package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trento-project/trento/internal"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/internal/sapsystem"
)

type SAPSystemRow struct {
	SID            string
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

func NewSAPSystemsTable(sapSystemsList sapsystem.SAPSystemsList, hostList hosts.HostList) SAPSystemsTable {
	var sapSystemsTable SAPSystemsTable
	rowsBySID := make(map[string]*SAPSystemRow)

	for _, s := range sapSystemsList {
		if s.Type != sapsystem.Application {
			continue
		}

		sapSystem, ok := rowsBySID[s.SID]
		if !ok {
			sapSystem = &SAPSystemRow{
				SID: s.SID,
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

	return sapSystemsTable
}

func (t SAPSystemsTable) filter(sid []string) SAPSystemsTable {
	var filteredSAPSystemsTable SAPSystemsTable
	if len(sid) == 0 {
		return t
	}

	for _, r := range t {
		if internal.Contains(sid, r.SID) {
			filteredSAPSystemsTable = append(filteredSAPSystemsTable, r)
		}
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

func NewSAPSystemListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sapSystemsList sapsystem.SAPSystemsList
		query := c.Request.URL.Query()
		sidFilter := query["sid"]

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

		sapSystemsTable := NewSAPSystemsTable(sapSystemsList, hostList)
		sapSystemsTable = sapSystemsTable.filter(sidFilter)

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
