package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"

	"github.com/trento-project/trento/web/services"
)

type TargetsList []*Targets

type Targets struct {
	Targets []string          `json:"targets,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

// ApiGetPrometheusHttpSdTargets godoc
// @Summary Get prometheus HTTP SD targets
// @Produce json
// @Success 200 {object} TargetsList
// @Error 500
// @Router /prometheus/targets [get]
func ApiGetPrometheusHttpSdTargets(s services.PrometheusService) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetsList, err := s.GetHttpSDTargets()
		if err != nil {
			c.Error(err)
			return
		}

		var targetsListJson TargetsList

		mapstructure.Decode(targetsList, &targetsListJson)

		c.JSON(http.StatusOK, targetsListJson)
	}
}
