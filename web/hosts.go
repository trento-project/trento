package web

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/cloud"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/environments"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/internal/sapsystem"
)

func NewHostsHealthContainer(hostList hosts.HostList) *HealthContainer {
	h := &HealthContainer{}
	for _, host := range hostList {
		switch host.Health() {
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

func NewHostListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()
		queryFilter := hosts.CreateFilterMetaQuery(query)
		healthFilter := query["health"]

		hostList, err := hosts.Load(client, queryFilter, healthFilter)
		if err != nil {
			_ = c.Error(err)
			return
		}

		filters, err := loadFilters(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		hContainer := NewHostsHealthContainer(hostList)
		hContainer.Layout = "horizontal"

		page := c.DefaultQuery("page", "1")
		perPage := c.DefaultQuery("per_page", "10")
		pagination := NewPaginationWithStrings(len(hostList), page, perPage)
		firstElem, lastElem := pagination.GetSliceNumbers()

		c.HTML(http.StatusOK, "hosts.html.tmpl", gin.H{
			"Hosts":           hostList[firstElem:lastElem],
			"Filters":         filters,
			"AppliedFilters":  query,
			"HealthContainer": hContainer,
			"Pagination":      pagination,
		})
	}
}

func loadFilters(client consul.Client) (map[string][]string, error) {
	filter_data := make(map[string][]string)

	envs, err := environments.Load(client)
	if err != nil {
		return nil, errors.Wrap(err, "could not get the filters")
	}

	for envKey, envValue := range envs {
		filter_data["environments"] = append(filter_data["environments"], envKey)
		for landKey, landValue := range envValue.Landscapes {
			filter_data["landscapes"] = append(filter_data["landscapes"], landKey)
			for sysKey, _ := range landValue.SAPSystems {
				filter_data["sapsystems"] = append(filter_data["sapsystems"], sysKey)
			}
		}
	}

	sort.Strings(filter_data["environments"])
	sort.Strings(filter_data["landscapes"])
	sort.Strings(filter_data["sapsystems"])

	return filter_data, nil
}

func loadHealthChecks(client consul.Client, node string) ([]*consulApi.HealthCheck, error) {

	checks, _, err := client.Health().Node(node, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for health checks")
	}

	return checks, nil
}

func NewHostHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		catalogNode, _, err := client.Catalog().Node(name, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if catalogNode == nil {
			_ = c.Error(NotFoundError("could not find host"))
			return
		}

		checks, err := loadHealthChecks(client, name)
		if err != nil {
			_ = c.Error(err)
			return
		}

		systems, err := sapsystem.Load(client, name)
		if err != nil {
			_ = c.Error(err)
			return
		}

		cloudData, err := cloud.Load(client, name)
		if err != nil {
			_ = c.Error(err)
			return
		}

		host := hosts.NewHost(*catalogNode.Node, client)
		c.HTML(http.StatusOK, "host.html.tmpl", gin.H{
			"Host":         &host,
			"HealthChecks": checks,
			"SAPSystems":   systems,
			"CloudData":    cloudData,
		})
	}
}
