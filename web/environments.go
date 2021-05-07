package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/environments"
)


func NewEnvironmentsListHandler(client consul.Client) gin.HandlerFunc {
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

func NewEnvironmentListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		environments, err := environments.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "environment.html.tmpl", gin.H{
			"Environments": environments,
			"EnvName":      c.Param("env"),
		})
	}
}

func NewLandscapesListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var env string = ""

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

func NewLandscapeListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var env string = ""

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

func NewSAPSystemsListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var env string = ""
		var land string = ""

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

func NewSAPSystemHostsListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var env string = ""
		var land string = ""

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

		c.HTML(http.StatusOK, "sapsystem.html.tmpl", gin.H{
			"Environments": environments,
			"EnvName":      env,
			"LandName":     land,
			"SAPSysName":   c.Param("sys"),
		})
	}
}
