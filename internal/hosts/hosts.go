package hosts

import (
	"fmt"
	"sort"
	"strings"

	"github.com/trento-project/trento/internal/tags"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/trento-project/trento/internal"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/sapsystem"
)

const TrentoPrefix string = "trento-"

type HostList []*Host

type Host struct {
	consulApi.Node
	client consul.Client
}

func (h HostList) Health() string {
	var checks consulApi.HealthChecks
	for _, n := range h {
		c, _, _ := n.client.Health().Node(n.Name(), nil)
		checks = append(checks, c...)
	}
	if checks != nil {
		return checks.AggregatedStatus()
	}

	return ""
}

func NewHost(node consulApi.Node, client consul.Client) Host {
	return Host{node, client}
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

func (n *Host) GetAgentVersionString() string {
	version, ok := n.TrentoMeta()["trento-agent-version"]

	if !ok {
		return "Not running"
	}

	return "v" + version
}

func (n *Host) GetSAPSystems() (sapsystem.SAPSystemsMap, error) {
	systems, err := sapsystem.Load(n.client, n.Name())
	if err != nil {
		return nil, err
	}
	return systems, nil
}

// Use github.com/hashicorp/go-bexpr to create the filter
// https://github.com/hashicorp/consul/blob/master/agent/consul/catalog_endpoint.go#L298
func CreateFilterMetaQuery(query map[string][]string) string {
	var filters []string
	var operator string
	// Need to sort the keys to have stable output. Mostly for unit testing
	sortedQuery := sortKeys(query)

	if len(query) != 0 {
		var filter string
		for _, key := range sortedQuery {
			switch key {
			case "trento-sap-systems":
				operator = "contains"
			default:
				operator = "=="
			}
			if strings.HasPrefix(key, TrentoPrefix) {
				filter = ""
				values := query[key]
				for _, value := range values {
					filter = fmt.Sprintf("%sMeta[\"%s\"] %s \"%s\"", filter, key, operator, value)

					if values[len(values)-1] != value {
						filter = fmt.Sprintf("%s or ", filter)
					}
				}
				filters = append(filters, fmt.Sprintf("(%s)", filter))
			}
		}
	}
	return strings.Join(filters, " and ")
}

func Load(client consul.Client, queryFilter string, healthFilter []string, tagsFilter []string) (HostList, error) {
	var hosts = HostList{}

	query := &consulApi.QueryOptions{Filter: queryFilter}
	consulNodes, _, err := client.Catalog().Nodes(query)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for nodes")
	}
	for _, node := range consulNodes {
		populatedHost := &Host{*node, client}
		// This check could be done in the frontend maybe
		if len(healthFilter) > 0 && !internal.Contains(healthFilter, populatedHost.Health()) {
			continue
		}

		if len(tagsFilter) > 0 {
			tagFound := false
			t := tags.NewTags(client)
			hostTags, err := t.GetAllByResource(tags.HostResourceType, node.Node)
			if err != nil {
				return nil, errors.Wrap(err, "could not query Tags for node")
			}

			for _, t := range tagsFilter {
				if internal.Contains(hostTags, t) {
					tagFound = true
					break
				}
			}

			if !tagFound {
				continue
			}
		}

		hosts = append(hosts, populatedHost)
	}

	return hosts, nil
}

func sortKeys(m map[string][]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
