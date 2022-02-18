package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trento-project/trento/web/services"
)

// ApiSAPSystemsHealthSummaryHandler godoc
// @Summary Retrieve SAP Systems Health Summary
// @Accept json
// @Produce json
// @Success 200 {object} models.HealthSummary
// @Failure 500 {object} map[string]string
// @Router /sapsystems/health [get]
func ApiSAPSystemsHealthSummaryHandler(healthSummaryService services.HealthSummaryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		healthSummary, err := healthSummaryService.GetHealthSummary()
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, healthSummary)
	}
}
