package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

type ClustersSettingsResponse models.ClustersSettings

// ApiGetClustersSettingsHandler godoc
// @Summary Retrieve Settings for all the clusters. Cluster's Selected checks and Hosts connection settings
// @Accept json
// @Produce json
// @Success 200 {object} ClustersSettingsResponse
// @Failure 500 {object} map[string]string
// @Router /internal/clusters/settings [get]
func ApiGetClustersSettingsHandler(clusters services.ClustersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clustersSettings, err := clusters.GetAllClustersSettings()

		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, clustersSettings)
	}
}
