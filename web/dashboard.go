package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html.tmpl", struct{}{})
}
