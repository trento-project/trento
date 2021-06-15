package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/environments"
)

func NewEnvironmentListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		environments, err := environments.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "environments.html.tmpl", gin.H{
			"Environments": environments,
		})
	}
}

func NewEnvironmentHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		environments, err := environments.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		envName := c.Param("env")
		_, ok := environments[envName]
		if !ok {
			_ = c.Error(NotFoundError("could not find environment"))
			return
		}

		c.HTML(http.StatusOK, "environment.html.tmpl", gin.H{
			"Environments": environments,
			"EnvName":      envName,
		})
	}
}

func NewLandscapeListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var env = ""

		query := c.Request.URL.Query()
		if len(query["environment"]) > 0 {
			env = query["environment"][0]
		}

		environments, err := environments.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "landscapes.html.tmpl", gin.H{
			"Environments": environments,
			"EnvName":      env,
		})
	}
}

func NewLandscapeHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var env = ""

		query := c.Request.URL.Query()
		if len(query["environment"]) > 0 {
			env = query["environment"][0]
		}

		environments, err := environments.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "landscape.html.tmpl", gin.H{
			"Environments": environments,
			"EnvName":      env,
			"LandName":     c.Param("land"),
		})
	}
}

func NewSAPSystemListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var env = ""
		var land = ""

		query := c.Request.URL.Query()
		if len(query["environment"]) > 0 {
			env = query["environment"][0]
		}

		if len(query["landscape"]) > 0 {
			land = query["landscape"][0]
		}

		environments, err := environments.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "sapsystems.html.tmpl", gin.H{
			"Environments": environments,
			"EnvName":      env,
			"LandName":     land,
		})
	}
}

func NewSAPSystemHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var envName, landName string

		query := c.Request.URL.Query()
		if len(query["environment"]) > 0 {
			envName = query["environment"][0]
		}

		if len(query["landscape"]) > 0 {
			landName = query["landscape"][0]
		}

		envs, err := environments.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		environment, ok := envs[envName]
		if !ok {
			_ = c.Error(NotFoundError("could not find environment"))
			return
		}

		landscape, ok := environment.Landscapes[landName]
		if !ok {
			_ = c.Error(NotFoundError("could not find landscape"))
			return
		}

		system, ok := landscape.SAPSystems[c.Param("sys")]
		if !ok {
			_ = c.Error(NotFoundError("could not find system"))
			return
		}

		c.HTML(http.StatusOK, "sapsystem.html.tmpl", gin.H{
			"Environment": environment,
			"Landscape":   landscape,
			"SAPSystem":   system,
		})
	}
}
