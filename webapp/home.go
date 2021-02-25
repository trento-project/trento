package webapp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func homeHandler(c *gin.Context) {
	viewModel := gin.H{
		"title": "SUSE Console for SAP Applications",
	}
	c.HTML(http.StatusOK, "home.html.tmpl", viewModel)
}
