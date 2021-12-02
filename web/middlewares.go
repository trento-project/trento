package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/web/services"
)

func EulaMiddleware(settings services.SettingsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		eulaAccepted, err := settings.IsEulaAccepted()
		if err != nil {
			log.Error(err)
			c.HTML(http.StatusInternalServerError, "error.html.tmpl", gin.H{"Error": "EULA Error :-("})
			c.Abort()
			return
		}

		if c.Request.URL.Path != "/accept-eula" && c.Request.URL.Path != "/eula" && !eulaAccepted {
			c.Redirect(http.StatusFound, "/eula")
			c.HTML(http.StatusOK, "eula.html.tmpl", gin.H{"Title": "License agreement"})
			c.Abort()
			return
		}
		c.Next()
		return
	}
}
