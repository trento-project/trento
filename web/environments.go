package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/environments"
)

func NewEnvironmentListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		environments, err := environments.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "environments.html.tmpl", gin.H{
			"Environments": environments,
		})
	}
}

func NewEnvironmentHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		environments, err := environments.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		envName := c.Param("env")
		_, ok := environments[envName]
		if !ok {
			_ = c.Error(NotFoundError("could not find environment"))
			return
		}

		c.HTML(http.StatusOK, "environment.html.tmpl", gin.H{
			"Environments": environments,
			"EnvName":      envName,
		})
	}
}

func NewLandscapeListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var envName string

		query := c.Request.URL.Query()
		if len(query["environment"]) > 0 {
			envName = query["environment"][0]
		}

		environments, err := environments.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "landscapes.html.tmpl", gin.H{
			"Environments": environments,
			"EnvName":      envName,
		})
	}
}

func NewLandscapeHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var envName string
		landName := c.Param("land")

		query := c.Request.URL.Query()
		if len(query["environment"]) > 0 {
			envName = query["environment"][0]
		}

		environments, err := environments.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		environment, ok := environments[envName]
		if !ok {
			_ = c.Error(NotFoundError("could not find environment"))
			return
		}

		_, ok = environment.Landscapes[landName]
		if !ok {
			_ = c.Error(NotFoundError("could not find landscape"))
			return
		}

		c.HTML(http.StatusOK, "landscape.html.tmpl", gin.H{
			"Environments": environments,
			"EnvName":      envName,
			"LandName":     landName,
		})
	}
}

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
		var envName, landName string

		query := c.Request.URL.Query()
		if len(query["environment"]) > 0 {
			envName = query["environment"][0]
		}

		if len(query["landscape"]) > 0 {
			landName = query["landscape"][0]
		}

		environments, err := environments.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "sapsystems.html.tmpl", gin.H{
			"Environments": environments,
			"EnvName":      envName,
			"LandName":     landName,
		})
	}
}

func NewSAPSystemHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var envName, landName string

		query := c.Request.URL.Query()
		if len(query["environment"]) > 0 {
			envName = query["environment"][0]
		}

		if len(query["landscape"]) > 0 {
			landName = query["landscape"][0]
		}

		envs, err := environments.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		environment, ok := envs[envName]
		if !ok {
			_ = c.Error(NotFoundError("could not find environment"))
			return
		}

		landscape, ok := environment.Landscapes[landName]
		if !ok {
			_ = c.Error(NotFoundError("could not find landscape"))
			return
		}

		system, ok := landscape.SAPSystems[c.Param("sys")]
		if !ok {
			_ = c.Error(NotFoundError("could not find system"))
			return
		}

		c.HTML(http.StatusOK, "sapsystem.html.tmpl", gin.H{
			"Environment": environment,
			"Landscape":   landscape,
			"SAPSystem":   system,
		})
	}
}
