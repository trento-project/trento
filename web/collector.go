package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trento-project/trento/web/datapipeline"
	"github.com/trento-project/trento/web/services"
)

// ApiCollectDataHandler handles the request to collect agent data from the API
func ApiCollectDataHandler(collectorService services.CollectorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var e datapipeline.DataCollectedEvent

		err := c.BindJSON(&e)
		if err != nil {
			_ = c.Error(err)
			return
		}

		err = collectorService.StoreEvent(&e)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Writer.WriteHeader(http.StatusAccepted)
	}
}
