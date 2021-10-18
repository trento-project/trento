package web

import (
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func setupTestDependencies() Dependencies {
	return Dependencies{
		engine: gin.Default(),
		store:  cookie.NewStore([]byte("secret")),
	}
}
