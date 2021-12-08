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

		sapSystems, err := sapSystemsService.GetAllApplications(tagsFilter, page)
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

		count, err := sapSystemsService.GetApplicationsCount()
		if err != nil {
			_ = c.Error(err)
			return
		}
		pagination := NewPagination(count, pageNumber, pageSize)

		c.HTML(http.StatusOK, "sap_systems.html.tmpl", gin.H{
			"Type":           models.SAPSystemTypeApplication,
			"SAPSystems":     sapSystems,
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

		databases, err := sapSystemsService.GetAllDatabases(tagsFilter, page)
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

		count, err := sapSystemsService.GetDatabasesCount()
		if err != nil {
			_ = c.Error(err)
			return
		}
		pagination := NewPagination(count, pageNumber, pageSize)

		c.HTML(http.StatusOK, "sap_systems.html.tmpl", gin.H{
			"Type":           models.SAPSystemTypeDatabase,
			"SAPSystems":     databases,
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

		// We will send the 1st entry by now, as only use the layout, which is repeated among all the
		// SAP instances within a System. It does not resolve the HANA SR scenario in any case
		c.HTML(http.StatusOK, "sap_system.html.tmpl", gin.H{
			"SAPSystem": sapSystem,
			"Hosts":     hosts,
			"HideTags":  true,
		})
	}
}
