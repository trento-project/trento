package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApiPingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
