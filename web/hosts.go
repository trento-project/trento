package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"

	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

func NewHostsHealthContainer(hostList models.HostList) *HealthContainer {
	h := &HealthContainer{}
	for _, host := range hostList {
		switch host.Health {
		case models.HostHealthPassing:
			h.PassingCount += 1
		case models.HostHealthWarning:
			h.WarningCount += 1
		case models.HostHealthCritical:
			h.CriticalCount += 1
		}
	}
	return h
}

func NewHostListHandler(hostsService services.HostsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()

		hostsFilter := &services.HostsFilter{
			SIDs:   query["sids"],
			Health: query["health"],
			Tags:   query["tags"],
		}

		pageNumber, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil {
			pageNumber = 1
		}
		pageSize, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
		if err != nil {
			pageSize = 10
		}

		page := &services.Page{
			Number: pageNumber,
			Size:   pageSize,
		}

		hostList, err := hostsService.GetAll(hostsFilter, page)
		if err != nil {
			_ = c.Error(err)
			return
		}

		filterSIDs, err := hostsService.GetAllSIDs()
		if err != nil {
			_ = c.Error(err)
			return
		}

		filterTags, err := hostsService.GetAllTags()
		if err != nil {
			_ = c.Error(err)
			return
		}

		count, err := hostsService.GetCount()
		if err != nil {
			_ = c.Error(err)
			return
		}
		pagination := NewPagination(count, pageNumber, pageSize)

		hContainer := NewHostsHealthContainer(hostList)
		hContainer.Layout = "horizontal"

		c.HTML(http.StatusOK, "hosts.html.tmpl", gin.H{
			"Hosts":           hostList,
			"AppliedFilters":  query,
			"FilterSIDs":      filterSIDs,
			"FilterTags":      filterTags,
			"Pagination":      pagination,
			"HealthContainer": hContainer,
		})
	}
}

func ApiHostHeartbeatHandler(hostService services.HostsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		agentID := c.Param("id")

		err := hostService.Heartbeat(agentID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusNoContent, gin.H{})
	}
}

func getTrentoAgentCheck(client consul.Client, node string) (*consulApi.HealthCheck, error) {
	checks, _, err := client.Health().Node(node, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for health checks")
	}

	var trentoAgentCheck *consulApi.HealthCheck

	for _, check := range checks {
		if check.CheckID == "trentoAgent" {
			trentoAgentCheck = check
			break
		}
	}

	return trentoAgentCheck, nil
}

func NewHostHandler(hostsService services.HostsService, subsService services.SubscriptionsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		host, err := hostsService.GetByName(name)
		if err != nil {
			_ = c.Error(NotFoundError("could not find host"))
			return
		}
		if host == nil {
			_ = c.Error(NotFoundError("could not find host"))
			return
		}

		subs, err := subsService.GetHostSubscriptions(name)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "host.html.tmpl", gin.H{
			"Host":          &host,
			"Subscriptions": subs,
		})
	}
}
