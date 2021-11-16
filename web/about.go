package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/trento-project/trento/web/services"
)

func NewAboutHandler(s services.SubscriptionsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		isPremium, subscribedCount, err := s.IsPremium()
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "about.html.tmpl", gin.H{
			"Title":           defaultLayoutData.Title,
			"Version":         defaultLayoutData.Version,
			"IsPremium":       isPremium,
			"SubscribedCount": subscribedCount,
		})
	}
}
