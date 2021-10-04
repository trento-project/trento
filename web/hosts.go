package web

import (
	"net/http"
	"strings"

	"github.com/trento-project/trento/internal"

	"github.com/gin-gonic/gin"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/cloud"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/internal/sapsystem"

	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

func NewHostsHealthContainer(hostList hosts.HostList) *HealthContainer {
	h := &HealthContainer{}
	for _, host := range hostList {
		switch host.Health() {
		case consulApi.HealthPassing:
			h.PassingCount += 1
		case consulApi.HealthWarning:
			h.WarningCount += 1
		case consulApi.HealthCritical:
			h.CriticalCount += 1
		}
	}
	return h
}

func NewHostListHandler(client consul.Client, tagsService services.TagsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()
		queryFilter := hosts.CreateFilterMetaQuery(query)
		healthFilter := query["health"]
		tagsFilter := query["tags"]

		hostList, err := hosts.Load(client, queryFilter, healthFilter)
		if err != nil {
			_ = c.Error(err)
			return
		}

		hostsTags := make(map[string][]string)
		for _, h := range hostList {
			ht, err := tagsService.GetAllByResource(models.TagHostResourceType, h.Name())
			if err != nil {
				c.Error(err)
				return
			}
			hostsTags[h.Name()] = ht
		}

		hostList = filterHostsByTags(hostList, hostsTags, tagsFilter)

		hContainer := NewHostsHealthContainer(hostList)
		hContainer.Layout = "horizontal"

		page := c.DefaultQuery("page", "1")
		perPage := c.DefaultQuery("per_page", "10")
		pagination := NewPaginationWithStrings(len(hostList), page, perPage)
		firstElem, lastElem := pagination.GetSliceNumbers()

		c.HTML(http.StatusOK, "hosts.html.tmpl", gin.H{
			"Hosts":           hostList[firstElem:lastElem],
			"SIDs":            getAllSIDs(hostList),
			"Tags":            getAllTags(hostsTags),
			"AppliedFilters":  query,
			"HealthContainer": hContainer,
			"Pagination":      pagination,
			"HostsTags":       hostsTags,
		})
	}
}

func filterHostsByTags(hostList hosts.HostList, hostsTags map[string][]string, tagsFilter []string) hosts.HostList {
	if len(tagsFilter) == 0 {
		return hostList
	}
	var filteredHostList hosts.HostList

	for _, h := range hostList {
		for _, t := range tagsFilter {
			if internal.Contains(hostsTags[h.Name()], t) {
				filteredHostList = append(filteredHostList, h)
				break
			}
		}
	}

	return filteredHostList
}

func getAllSIDs(hostList hosts.HostList) []string {
	var sids []string
	set := make(map[string]struct{})

	for _, host := range hostList {
		for _, s := range strings.Split(host.TrentoMeta()["trento-sap-systems"], ",") {
			if s == "" {
				continue
			}

			_, ok := set[s]
			if !ok {
				sids = append(sids, s)
				set[s] = struct{}{}
			}
		}
	}

	return sids
}

func getAllTags(hostsTags map[string][]string) []string {
	var tags []string
	set := make(map[string]struct{})

	for _, ht := range hostsTags {
		for _, t := range ht {
			_, ok := set[t]
			if !ok {
				tags = append(tags, t)
				set[t] = struct{}{}
			}
		}
	}

	return tags
}

func getTrentoAgentCheck(client consul.Client, node string) (*consulApi.HealthCheck, error) {
	checks, _, err := client.Health().Node(node, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for health checks")
	}

	var trentoAgentCheck *consulApi.HealthCheck

	for _, check := range checks {
		if check.CheckID == "trentoAgent" {
			trentoAgentCheck = check
			break
		}
	}

	return trentoAgentCheck, nil
}

func NewHostHandler(client consul.Client, subsService services.SubscriptionsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		catalogNode, _, err := client.Catalog().Node(name, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if catalogNode == nil {
			_ = c.Error(NotFoundError("could not find host"))
			return
		}

		trentoAgentCheck, err := getTrentoAgentCheck(client, name)
		if err != nil {
			_ = c.Error(err)
			return
		}

		systems, err := sapsystem.Load(client, name)
		if err != nil {
			_ = c.Error(err)
			return
		}

		cloudData, err := cloud.Load(client, name)
		if err != nil {
			_ = c.Error(err)
			return
		}

		subs, err := subsService.GetHostSubscriptions(name)
		if err != nil {
			_ = c.Error(err)
			return
		}

		host := hosts.NewHost(*catalogNode.Node, client)
		c.HTML(http.StatusOK, "host.html.tmpl", gin.H{
			"Host":             &host,
			"TrentoAgentCheck": trentoAgentCheck,
			"SAPSystems":       systems,
			"CloudData":        cloudData,
			"Subscriptions":    subs,
		})
	}
}
