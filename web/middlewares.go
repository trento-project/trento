package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/web/services"
)

func EulaMiddleware(premiumDetection services.PremiumDetectionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		requiresEulaAcceptance, err := premiumDetection.RequiresEulaAcceptance()

		if err != nil {
			log.Error(err)
			c.HTML(http.StatusInternalServerError, "error.html.tmpl", gin.H{"Error": "Error checking whether the instance is premium or not"})
			c.Abort()
			return
		}

		if c.Request.URL.Path != "/accept-eula" && c.Request.URL.Path != "/eula" && requiresEulaAcceptance {
			c.Redirect(http.StatusFound, "/eula")
			c.HTML(http.StatusOK, "eula.html.tmpl", gin.H{"Title": "License agreement"})
			c.Abort()
			return
		}
		c.Next()
	}
}
