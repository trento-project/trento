package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeData struct {
	Title     string
	Paragraph string
}

func HomeHandler(c *gin.Context) {
	data := HomeData{
		Title: defaultLayoutData.Title,
		Paragraph: `An open cloud-native web console improving on the
				life of SAP Applications administrators`,
	}
	c.HTML(http.StatusOK, "home.html.tmpl", data)
}
