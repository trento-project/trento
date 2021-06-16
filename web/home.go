package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeData struct {
	Title string
}

func HomeHandler(c *gin.Context) {
	data := HomeData{
		Title: defaultLayoutData.Title,
	}
	c.HTML(http.StatusOK, "home.html.tmpl", data)
}
