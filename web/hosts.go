package web

import (
	"net/http"
	"strconv"

	"github.com/trento-project/trento/web/models"

	"github.com/trento-project/trento/web/service"

	"github.com/gin-gonic/gin"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/cloud"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/internal/sapsystem"
)

func NewHostsHealthContainer(hosts []models.Host) *HealthContainer {
	h := &HealthContainer{}
	for _, host := range hosts {
		switch host.Health {
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

func NewHostListHandler(h service.IHostsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageNr, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
		page := &service.Page{
			PageNr:   pageNr,
			PageSize: pageSize,
		}
		pagination := NewPagination(h.GetHostsCount(), pageNr, pageSize)

		query := c.Request.URL.Query()
		filters := map[string][]string{
			"health":      query["health"],
			"sap_system":  query["sap_system"],
			"landscape":   query["landscape"],
			"environment": query["environment"],
		}
		hs := h.GetHosts(page, filters)

		hc := NewHostsHealthContainer(hs)
		hc.Layout = "horizontal"

		filtersOptions := map[string][]string{
			"sap_system":  h.GetHostsSAPSystems(),
			"landscape":   h.GetHostsLandscapes(),
			"environment": h.GetHostsEnvironments(),
		}

		c.HTML(http.StatusOK, "hosts.html.tmpl", gin.H{
			"Hosts":           hs,
			"Pagination":      pagination,
			"HealthContainer": hc,
			"Filters":         filtersOptions,
			"AppliedFilters":  query,
		})
	}
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

func NewHAChecksHandler(client consul.Client) gin.HandlerFunc {
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

		host := hosts.NewHost(*catalogNode.Node, client)
		c.HTML(http.StatusOK, "ha_checks.html.tmpl", gin.H{
			"Hostname": host.Name(),
			"HAChecks": host.HAChecks(),
		})
	}
}
