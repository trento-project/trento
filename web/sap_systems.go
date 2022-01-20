package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

func NewSAPSystemListHandler(sapSystemsService services.SAPSystemsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()

		tagsFilter := &services.SAPSystemFilter{
			Tags: query["tags"],
			SIDs: query["sids"],
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

		paginatedSapSystems, err := sapSystemsService.GetAllApplications(tagsFilter, page)
		if err != nil {
			_ = c.Error(err)
			return
		}

		sapSystems, err := sapSystemsService.GetAllApplications(tagsFilter, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		filterSIDs, err := sapSystemsService.GetAllApplicationsSIDs()
		if err != nil {
			_ = c.Error(err)
			return
		}

		filterTags, err := sapSystemsService.GetAllApplicationsTags()
		if err != nil {
			_ = c.Error(err)
			return
		}

		pagination := NewPagination(len(sapSystems), pageNumber, pageSize)

		c.HTML(http.StatusOK, "sap_systems.html.tmpl", gin.H{
			"Type":           models.SAPSystemTypeApplication,
			"SAPSystems":     paginatedSapSystems,
			"AppliedFilters": query,
			"FilterSIDs":     filterSIDs,
			"FilterTags":     filterTags,
			"Pagination":     pagination,
		})
	}
}

func NewHANADatabaseListHandler(sapSystemsService services.SAPSystemsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()

		tagsFilter := &services.SAPSystemFilter{
			Tags: query["tags"],
			SIDs: query["sids"],
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

		paginatedDatabases, err := sapSystemsService.GetAllDatabases(tagsFilter, page)
		if err != nil {
			_ = c.Error(err)
			return
		}

		databases, err := sapSystemsService.GetAllDatabases(tagsFilter, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		filterSIDs, err := sapSystemsService.GetAllDatabasesSIDs()
		if err != nil {
			_ = c.Error(err)
			return
		}

		filterTags, err := sapSystemsService.GetAllDatabasesTags()
		if err != nil {
			_ = c.Error(err)
			return
		}

		pagination := NewPagination(len(databases), pageNumber, pageSize)

		c.HTML(http.StatusOK, "sap_systems.html.tmpl", gin.H{
			"Type":           models.SAPSystemTypeDatabase,
			"SAPSystems":     paginatedDatabases,
			"AppliedFilters": query,
			"FilterSIDs":     filterSIDs,
			"FilterTags":     filterTags,
			"Pagination":     pagination,
		})
	}
}

func NewSAPResourceHandler(hostsService services.HostsService, sapSystemsService services.SAPSystemsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		sapSystem, err := sapSystemsService.GetByID(id)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if sapSystem == nil {
			_ = c.Error(NotFoundError("could not find system"))
			return
		}

		hosts, err := hostsService.GetAllBySAPSystemID(id)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "sap_system.html.tmpl", gin.H{
			"SAPSystem":      sapSystem,
			"Hosts":          hosts,
			"HideSAPSystems": true,
			"HideTags":       true,
		})
	}
}
