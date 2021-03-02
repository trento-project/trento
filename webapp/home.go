package webapp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func homeHandler(c *gin.Context) {
	viewModel := gin.H{
		"title":  "SUSE Console for SAP Applications",
		"footer": "© 2019-2020 SUSE, all rights reserved.",
	}
	c.HTML(http.StatusOK, "home", viewModel)
}
