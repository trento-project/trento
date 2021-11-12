package web

import (
	"net/http"

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

func NewHostListNextHandler(hostsService services.HostsNextService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()

		hostList, err := hostsService.GetAll(query)
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

		page := c.DefaultQuery("page", "1")
		perPage := c.DefaultQuery("per_page", "10")
		pagination := NewPaginationWithStrings(len(hostList), page, perPage)
		firstElem, lastElem := pagination.GetSliceNumbers()

		hContainer := NewHostsNextHealthContainer(hostList)
		hContainer.Layout = "horizontal"

		c.HTML(http.StatusOK, "hosts_next.html.tmpl", gin.H{
			"Hosts":           hostList[firstElem:lastElem],
			"AppliedFilters":  query,
			"FilterSIDs":      filterSIDs,
			"FilterTags":      filterTags,
			"Pagination":      pagination,
			"HealthContainer": hContainer,
		})
	}
}

func ApiHostHeartbeatHandler(hostService services.HostsNextService) gin.HandlerFunc {
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
