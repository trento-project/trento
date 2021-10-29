package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:generate swag init -g api.go -o ../docs/api
// @title Trento API
// @version 1.0
// @description Trento API

// @contact.name Trento Project
// @contact.url https://www.trento-project.io
// @contact.email  trento-project@suse.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api
// @schemes http

func ApiPingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
