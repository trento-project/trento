package web

import (
	//"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"
)

const KVEnvironmentsPath string = "trento/environments"
const envIndex int = 2
const landIndex int = 4

type SAPSystem struct {
	Name  string
	Hosts HostList
}

type SAPSystemList map[string]*SAPSystem

type Landscape struct {
	Name       string
	SAPSystems SAPSystemList
}

type LandscapeList map[string]*Landscape

type Environment struct {
	Name       string
	Landscapes LandscapeList
}

type EnvironmentList map[string]*Environment

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

func NewLandscapesListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		environments, err := loadEnvironments(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "landscapes.html.tmpl", gin.H{
			"Environments": environments,
			"EnvName":      c.Param("env"),
		})
	}
}

func NewSAPSystemsListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		environments, err := loadEnvironments(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "sapsystems.html.tmpl", gin.H{
			"Environments": environments,
			"EnvName":      c.Param("env"),
			"LandName":     c.Param("land"),
		})
	}
}

func NewSAPSystemHostsListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		environments, err := loadEnvironments(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "sapsystem.html.tmpl", gin.H{
			"Environments": environments,
			"EnvName":      c.Param("env"),
			"LandName":     c.Param("land"),
			"SAPSysName":   c.Param("sys"),
		})
	}
}

/* loadEnvironments needs a fixed kv structure to work. Here an example
trento/environments/
trento/environments/env1/
trento/environments/env1/landscapes/
trento/environments/env1/landscapes/land1/
trento/environments/env1/landscapes/land1/sapsystems/
trento/environments/env1/landscapes/land1/sapsystems/sys1/
trento/environments/env1/landscapes/land1/sapsystems/sys2/
trento/environments/env1/landscapes/land2/
trento/environments/env1/landscapes/land2/sapsystems/
trento/environments/env1/landscapes/land2/sapsystems/sys3/
trento/environments/env1/landscapes/land2/sapsystems/sys4/
trento/environments/env2/
trento/environments/env2/landscapes/
trento/environments/env2/landscapes/land3/
trento/environments/env2/landscapes/land3/sapsystems/
trento/environments/env2/landscapes/land3/sapsystems/sys5/
*/
func loadEnvironments(client consul.Client) (EnvironmentList, error) {
	var (
		environments = EnvironmentList{}
		reserveKeys  = []string{"environments", "landscapes", "sapsystems"}
	)

	entries, _, err := client.KV().Keys(KVEnvironmentsPath, "", nil)
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
