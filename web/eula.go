package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/web/services"
)

func EulaAcceptHandler(settings services.SettingsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := settings.AcceptEula()
		if err != nil {
			log.Error(err)
			c.HTML(http.StatusInternalServerError, "error.html.tmpl", gin.H{"Error": "There was an error accepting the EULA. Please try again."})
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}
