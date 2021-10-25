package web

import (
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func setupTestDependencies() Dependencies {
	return Dependencies{
		webEngine:       gin.Default(),
		collectorEngine: gin.Default(),
		store:           cookie.NewStore([]byte("secret")),
	}
}

func setupTestConfig() *Config {
	return &Config{
		Host: "",
		Port: 80,
	}
}
