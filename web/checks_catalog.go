package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/trento-project/trento/web/services"
)

func NewChecksCatalogHandler(s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		checkList, err := s.GetChecksCatalogByGroup()
		if err != nil {
			tipMsg := AlertCatalogNotFound().Text
			_ = c.Error(InternalServerError(err.Error()))
			_ = c.Error(InternalServerError(tipMsg))
			return
		}

		c.HTML(http.StatusOK, "checks_catalog.html.tmpl", gin.H{
			"ChecksCatalog": checkList.OrderById(),
		})
	}
}
