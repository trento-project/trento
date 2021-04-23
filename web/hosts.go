package web

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/aquasecurity/bench-common/check"
	"github.com/gin-gonic/gin"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"
)

const TrentoPrefix string = "trento-"

type HostList []*Host

type Host struct {
	consulApi.Node
	client consul.Client
}

func (n *Host) Health() string {
	checks, _, _ := n.client.Health().Node(n.Name(), nil)
	return checks.AggregatedStatus()
}

func (n *Host) Name() string {
	return n.Node.Node
}

func (n *Host) TrentoMeta() map[string]string {
	filtered_meta := make(map[string]string)

	for key, value := range n.Node.Meta {
		if strings.HasPrefix(key, TrentoPrefix) {
			filtered_meta[key] = value
		}
	}
	return filtered_meta
}

// todo: this method was rushed, needs to be completely rewritten to have the checker webservice decoupled in a dedicated HTTP client
func (n *Host) Checks() *check.Controls {
	checks := &check.Controls{}

	var err error
	resp, err := http.Get(fmt.Sprintf("http://%s:%d", n.Address, 8700)) // todo: use a Consul service instead of directly addressing the node IP and port
	if err != nil {
		log.Print(err)
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return nil
	}

	err = json.Unmarshal(body, checks)
	if err != nil {
		log.Print(err)
		return nil
	}
	return checks
}

// Use github.com/hashicorp/go-bexpr to create the filter
// https://github.com/hashicorp/consul/blob/master/agent/consul/catalog_endpoint.go#L298
func CreateFilterMetaQuery(query map[string][]string) string {
	var filters []string

	if len(query) != 0 {
		var filter string
		for key, values := range query {
			if strings.HasPrefix(key, TrentoPrefix) {
				filter = ""
				for _, value := range values {
					filter = fmt.Sprintf("%sMeta[\"%s\"] == \"%s\"", filter, key, value)
					if values[len(values)-1] != value {
						filter = fmt.Sprintf("%s or ", filter)
					}
				}
				filters = append(filters, filter)
			}
		}
	}
	return strings.Join(filters, " and ")
}

func NewHostsListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()
		query_filter := CreateFilterMetaQuery(query)
		health_filter := query["health"]

		hosts, err := loadHosts(client, query_filter, health_filter)
		if err != nil {
			_ = c.Error(err)
			return
		}

		filters, err := loadFilters(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "hosts.html.tmpl", gin.H{
			"Hosts":          hosts,
			"Filters":        filters,
			"AppliedFilters": query,
		})
	}
}

func loadHosts(client consul.Client, query_filter string, health_filter []string) (HostList, error) {
	var hosts = HostList{}

	query := &consulApi.QueryOptions{Filter: query_filter}
	consul_nodes, _, err := client.Catalog().Nodes(query)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for nodes")
	}
	for _, node := range consul_nodes {
		populated_host := &Host{*node, client}
		// This check could be done in the frontend maybe
		if len(health_filter) == 0 || contains(health_filter, populated_host.Health()) {
			hosts = append(hosts, populated_host)
		}
	}

	return hosts, nil
}

func loadFilters(client consul.Client) (map[string][]string, error) {
	filter_data := make(map[string][]string)

	environments, err := loadEnvironments(client)
	if err != nil {
		return nil, errors.Wrap(err, "could not get the filters")
	}

	for envKey, envValue := range environments {
		filter_data["environments"] = append(filter_data["environments"], envKey)
		for landKey, landValue := range envValue.Landscapes {
			filter_data["landscapes"] = append(filter_data["landscapes"], landKey)
			for sysKey, _ := range landValue.SAPSystems {
				filter_data["sapsystems"] = append(filter_data["sapsystems"], sysKey)
			}
		}
	}

	return filter_data, nil
}

func loadHealthChecks(client consul.Client, node string) ([]*consulApi.HealthCheck, error) {

	checks, _, err := client.Health().Node(node, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for health checks")
	}

	return checks, nil
}

func NewHostHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		catalogNode, _, err := client.Catalog().Node(name, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		checks, err := loadHealthChecks(client, name)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.HTML(http.StatusOK, "host.html.tmpl", gin.H{
			"Host":         &Host{*catalogNode.Node, client},
			"HealthChecks": checks,
		})
	}
}

func NewCheckHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		checkid := c.Param("checkid")
		catalogNode, _, err := client.Catalog().Node(name, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		host := &Host{*catalogNode.Node, client}
		c.HTML(http.StatusOK, "ha_checks.html.tmpl", gin.H{
			"HostName":     name,
			"CheckID":      checkid,
			"CheckContent": host.Checks(),
		})
	}
}
