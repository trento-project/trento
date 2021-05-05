package web

import (
	"net/http"
	//"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"
)

type EnvironmentHealth struct {
	Health    string
	HealthMap map[string]string
}

func (e *EnvironmentHealth) updateHealth(n string, h string) EnvironmentHealth {
	e.HealthMap[n] = h

	if h == "critical" {
		e.Health = h
	} else if h == "warning" && e.Health != "critical" {
		e.Health = h
	}

	return *e
}

type SAPSystem struct {
	Name  string
	Hosts HostList
}

type SAPSystemList map[string]*SAPSystem

func (s *SAPSystem) Health() EnvironmentHealth {
	var health = EnvironmentHealth{
		Health:    "passing",
		HealthMap: make(map[string]string),
	}

	for _, host := range s.Hosts {
		h := host.Health()
		health = health.updateHealth(host.Name(), h)
	}

	return health
}

type Landscape struct {
	Name       string
	SAPSystems SAPSystemList
}

type LandscapeList map[string]*Landscape

func (l *Landscape) Health() EnvironmentHealth {
	var health = EnvironmentHealth{
		Health:    "passing",
		HealthMap: make(map[string]string),
	}

	for _, system := range l.SAPSystems {
		h := system.Health().Health
		health = health.updateHealth(system.Name, h)
	}

	return health
}

func (l *Landscape) Ungrouped() bool {
	return l.Name == consul.KvUngrouped
}

type Environment struct {
	Name       string
	Landscapes LandscapeList
}

type EnvironmentList map[string]*Environment

func (e *Environment) Health() EnvironmentHealth {
	var health = EnvironmentHealth{
		Health:    "passing",
		HealthMap: make(map[string]string),
	}

	for _, land := range e.Landscapes {
		h := land.Health().Health
		health = health.updateHealth(land.Name, h)
	}

	return health
}

func (e *Environment) Ungrouped() bool {
	return e.Name == consul.KvUngrouped
}

func NewEnvironmentsListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		environments, err := loadEnvironments(client)
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
		environments, err := loadEnvironments(client)
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

		environments, err := loadEnvironments(client)
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

		environments, err := loadEnvironments(client)
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

		environments, err := loadEnvironments(client)
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

		environments, err := loadEnvironments(client)
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

func loadEnvironments(client consul.Client) (EnvironmentList, error) {
	var environments = EnvironmentList{}

	envs, err := client.KV().ListMap(consul.KvEnvironmentsPath, consul.KvEnvironmentsPath)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting the environments data")
	}

	for env, envValue := range envs {
		landscapes, err := loadLandscapes(client, env, envValue)
		if err != nil {
			return nil, err
		}
		environments[env] = &Environment{Name: env, Landscapes: landscapes}
	}

	return environments, nil
}

func loadLandscapes(client consul.Client, env string, envValue interface{}) (LandscapeList, error) {
	var landscapes = LandscapeList{}

	lands := envValue.(map[string]interface{})["landscapes"]

	for land, landValue := range lands.(map[string]interface{}) {
		sapsystems, err := loadSAPSystems(client, env, land, landValue)
		if err != nil {
			return nil, err
		}
		landscapes[land] = &Landscape{Name: land, SAPSystems: sapsystems}
	}

	return landscapes, nil
}

func loadSAPSystems(client consul.Client, env string, land string, landValue interface{}) (SAPSystemList, error) {
	var sapsystems = SAPSystemList{}

	syss := landValue.(map[string]interface{})["sapsystems"]

	for sys, _ := range syss.(map[string]interface{}) {
		query := CreateFilterMetaQuery(map[string][]string{
			"trento-sap-environment": []string{env},
			"trento-sap-landscape":   []string{land},
			"trento-sap-system":      []string{sys},
		})
		hosts, err := loadHosts(client, query, []string{})
		if err != nil {
			return nil, errors.Wrap(err, "could not query Consul for hosts")
		}

		sapsystems[sys] = &SAPSystem{Name: sys, Hosts: hosts}
	}

	return sapsystems, nil
}
