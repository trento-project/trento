package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trento-project/trento/web/services"
)

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

		c.HTML(http.StatusOK, "hosts_next.html.tmpl", gin.H{
			"Hosts":          hostList[firstElem:lastElem],
			"AppliedFilters": query,
			"FilterSIDs":     filterSIDs,
			"FilterTags":     filterTags,
			"Pagination":     pagination,
		})
	}
}
