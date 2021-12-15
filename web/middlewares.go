package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/web/services"
)

func EulaMiddleware(settings services.SettingsService, subscriptions services.SubscriptionsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		isPremium, err := subscriptions.IsTrentoPremium()
		if err != nil {
			log.Error(err)
			c.HTML(http.StatusInternalServerError, "error.html.tmpl", gin.H{"Error": "Error checking wether the instance is premium or not"})
			c.Abort()
			return
		}
		eulaAccepted, err := settings.IsEulaAccepted()
		if err != nil {
			log.Error(err)
			c.HTML(http.StatusInternalServerError, "error.html.tmpl", gin.H{"Error": "EULA Error :-("})
			c.Abort()
			return
		}

		if c.Request.URL.Path != "/accept-eula" && c.Request.URL.Path != "/eula" && !eulaAccepted && isPremium {
			c.Redirect(http.StatusFound, "/eula")
			c.HTML(http.StatusOK, "eula.html.tmpl", gin.H{"Title": "License agreement"})
			c.Abort()
			return
		}
		c.Next()
		return
	}
}
