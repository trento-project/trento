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
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "checks_catalog.html.tmpl", gin.H{
			"Checks": checkList,
		})
	}
}

func NewChecksResultHandler(s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		checkResultList, err := s.GetChecksResult("hog")
		if err != nil {
			_ = c.Error(err)
			return
		}

		checksMetadata, err := s.GetChecksCatalog()
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "checks_result.html.tmpl", gin.H{
			"ChecksResult": checkResultList,
			"ChecksMetadata": checksMetadata,
			"CheckedHosts": checkResultList.GetHostNames(),
		})
	}
}
