package web

import (
	//"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"
)

const envIndex int = 3
const landIndex int = 5

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

type Environment struct {
	Name       string
	Landscapes LandscapeList
}

type EnvironmentList map[string]*Environment

func (l *Environment) Health() EnvironmentHealth {
	var health = EnvironmentHealth{
		Health:    "passing",
		HealthMap: make(map[string]string),
	}

	for _, land := range l.Landscapes {
		h := land.Health().Health
		health = health.updateHealth(land.Name, h)
	}

	return health
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

/* loadEnvironments needs a fixed kv structure to work. Here an example
trento/v0/environments/
trento/v0/environments/env1/
trento/v0/environments/env1/landscapes/
trento/v0/environments/env1/landscapes/land1/
trento/v0/environments/env1/landscapes/land1/sapsystems/
trento/v0/environments/env1/landscapes/land1/sapsystems/sys1/
trento/v0/environments/env1/landscapes/land1/sapsystems/sys2/
trento/v0/environments/env1/landscapes/land2/
trento/v0/environments/env1/landscapes/land2/sapsystems/
trento/v0/environments/env1/landscapes/land2/sapsystems/sys3/
trento/v0/environments/env1/landscapes/land2/sapsystems/sys4/
trento/v0/environments/env2/
trento/v0/environments/env2/landscapes/
trento/v0/environments/env2/landscapes/land3/
trento/v0/environments/env2/landscapes/land3/sapsystems/
trento/v0/environments/env2/landscapes/land3/sapsystems/sys5/
*/
func loadEnvironments(client consul.Client) (EnvironmentList, error) {
	var (
		environments = EnvironmentList{}
		reserveKeys  = []string{"environments", "landscapes", "sapsystems"}
	)

	entries, _, err := client.KV().Keys(consul.KvEnvironmentsPath, "", nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for Environments KV values")
	}

	for _, entry := range entries {
		// Remove individual values, even though there is not any defined by now.
		if !strings.HasSuffix(entry, "/") {
			continue
		}

		keyValues := strings.Split(strings.TrimSuffix(entry, "/"), "/")
		lastKey := keyValues[len(keyValues)-1]
		lastKeyParent := keyValues[len(keyValues)-2]

		if contains(reserveKeys, lastKey) {
			continue
		}

		_, envFound := environments[lastKeyParent]
		if lastKeyParent == "environments" && !envFound {
			env := &Environment{Name: lastKey, Landscapes: make(LandscapeList)}
			environments[lastKey] = env
		}

		environments, err = loadLandscapes(client, environments, keyValues)
		if err != nil {
			return nil, errors.Wrap(err, "could not get the SAP landscapes")
		}
	}

	return environments, nil
}

func loadLandscapes(client consul.Client, environments EnvironmentList, values []string) (EnvironmentList, error) {
	lastKey := values[len(values)-1]
	lastKeyParent := values[len(values)-2]

	_, landFound := environments[lastKeyParent]
	if lastKeyParent == "landscapes" && !landFound {
		land := &Landscape{Name: lastKey, SAPSystems: make(SAPSystemList)}
		envName := values[envIndex]
		environments[envName].Landscapes[lastKey] = land
	}

	environments, err := loadSAPSystems(client, environments, values)
	if err != nil {
		return nil, errors.Wrap(err, "could not get the SAP systems")
	}

	return environments, nil
}

func loadSAPSystems(client consul.Client, environments EnvironmentList, values []string) (EnvironmentList, error) {
	lastKey := values[len(values)-1]
	lastKeyParent := values[len(values)-2]

	_, sysFound := environments[lastKeyParent]
	if lastKeyParent == "sapsystems" && !sysFound {
		envName := values[envIndex]
		landName := values[landIndex]
		// Get the nodes with these meta-data entries
		query := CreateFilterMetaQuery(map[string][]string{
			"trento-sap-environment": []string{envName},
			"trento-sap-landscape":   []string{landName},
			"trento-sap-system":      []string{lastKey},
		})
		hosts, err := loadHosts(client, query, []string{})
		if err != nil {
			return nil, errors.Wrap(err, "could not query Consul for hosts")
		}
		sapsystem := &SAPSystem{Name: lastKey, Hosts: hosts}

		environments[envName].Landscapes[landName].SAPSystems[lastKey] = sapsystem
	}

	return environments, nil
}
