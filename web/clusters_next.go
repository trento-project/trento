package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

func NewClusterListHealthContainer(clusterList models.ClusterList) *HealthContainer {
	h := &HealthContainer{}
	for _, c := range clusterList {
		switch c.Health {
		case models.CheckPassing:
			h.PassingCount += 1
		case models.CheckWarning:
			h.WarningCount += 1
		case models.CheckCritical:
			h.CriticalCount += 1
		}
	}
	return h
}

func NewClusterListNextHandler(clustersService services.ClustersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()

		clusterList, err := clustersService.GetAll(query)

		if err != nil {
			_ = c.Error(err)
			return
		}

		healthContainer := NewClusterListHealthContainer(clusterList)
		healthContainer.Layout = "horizontal"

		page := c.DefaultQuery("page", "1")
		perPage := c.DefaultQuery("per_page", "10")
		pagination := NewPaginationWithStrings(len(clusterList), page, perPage) // not nice here
		firstElem, lastElem := pagination.GetSliceNumbers()

		c.HTML(http.StatusOK, "clusters_next.html.tmpl", gin.H{
			"ClustersTable":   clusterList[firstElem:lastElem],
			"AppliedFilters":  query,
			"Pagination":      pagination,
			"HealthContainer": healthContainer,
		})
	}
}
