package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

func NewHostsNextHealthContainer(hostList models.HostList) *HealthContainer {
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

func NewHostListNextHandler(hostsService services.HostsService) gin.HandlerFunc {
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

		hContainer := NewHostsNextHealthContainer(hostList)
		hContainer.Layout = "horizontal"

		c.HTML(http.StatusOK, "hosts_next.html.tmpl", gin.H{
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
