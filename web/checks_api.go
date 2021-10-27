package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

type JSONSelectedChecks struct {
	SelectedChecks []string `json:"selected_checks" binding:"required"`
}

// ApiCheckResultsHandler godoc
// @Summary Get a specific cluster's check results
// @Produce json
// @Param cluster_id path string true "Cluster Id"
// @Success 200 {object} map[string]interface{}
// @Error 500
// @Router /api/clusters/{cluster_id}/results [get]
func ApiClusterCheckResultsHandler(client consul.Client, s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterId := c.Param("cluster_id")

		checkResults, err := s.GetChecksResultAndMetadataByCluster(clusterId)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, checkResults)
	}
}

// ApiCheckGetSelectedHandler godoc
// @Summary Get selected checks from resource
// @Accept json
// @Produce json
// @Param id path string true "Resource id"
// @Success 200 {object} JSONSelectedChecks
// @Failure 404 {object} map[string]string
// @Router /api/checks/{id}/selected [get]
func ApiCheckGetSelectedHandler(s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		selectedChecks, err := s.GetSelectedChecksById(id)
		if err != nil {
			_ = c.Error(NotFoundError("could not find check selection"))
			return
		}

		var jsonSelectedChecks JSONSelectedChecks
		jsonSelectedChecks.SelectedChecks = selectedChecks.SelectedChecks

		c.JSON(http.StatusOK, jsonSelectedChecks)
	}
}

// ApiCheckCreateSelectedHandler godoc
// @Summary Create check selection for the resource
// @Accept json
// @Produce json
// @Param id path string true "Resource id"
// @Param Body body JSONSelectedChecks true "Selected checks"
// @Success 201 {object} JSONSelectedChecks
// @Failure 500 {object} map[string]string
// @Router /api/checks/{id}/selected [post]
func ApiCheckCreateSelectedHandler(s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var r JSONSelectedChecks

		err := c.BindJSON(&r)
		if err != nil {
			_ = c.Error(BadRequestError("unable to parse JSON body"))
			return
		}

		err = s.CreateSelectedChecks(id, r.SelectedChecks)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, &r)
	}
}

// ApiCheckGetConnectionDataHandler godoc
// @Summary Get users connection data
// @Accept json
// @Produce json
// @Param id path string true "Checks group id"
// @Success 200 {object} map[string]models.ConnectionData
// @Failure 404 {object} map[string]string
// @Router /api/checks/{id}/connection_data [get]
func ApiCheckGetConnectionDataHandler(s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		connectionData, err := s.GetConnectionDataById(id)
		if err != nil {
			_ = c.Error(NotFoundError("could not find connection data"))
			return
		}

		c.JSON(http.StatusOK, connectionData)
	}
}

// ApiCheckCreateConnectionDataHandler godoc
// @Summary Create connection data for the node
// @Accept json
// @Produce json
// @Param id path string true "Checks group id"
// @Param Body body models.ConnectionData true "Checks connection user data"
// @Success 201 {object} models.ConnectionData
// @Failure 500 {object} map[string]string
// @Router /api/checks/{id}/connection_data [post]
func ApiCheckCreateConnectionDataHandler(s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var r models.ConnectionData

		err := c.BindJSON(&r)
		if err != nil {
			_ = c.Error(BadRequestError("unable to parse JSON body"))
			return
		}

		err = s.CreateConnectionData(id, r.Node, r.User)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, &r)
	}
}
