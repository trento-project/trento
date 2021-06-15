package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/internal/ruleset"
)

type Cluster struct {
	Name string
}

type ClusterList map[string]*Cluster

func NewClusterListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusters, err := cluster.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "clusters.html.tmpl", gin.H{
			"Clusters": clusters,
		})
	}
}

type rulesetsForm struct {
	Rulesets []string `form:"rulesets[]"`
}

func getRulesets(c *gin.Context, client consul.Client, h hosts.HostList) (ruleset.RuleSets, error) {
	var cForm rulesetsForm

	r, err := ruleset.NewRuleSets([]string{})
	if err != nil {
		return nil, errors.Wrap(err, "could not load embedded rulesets")
	}

	// POST action
	if c.Request.Method == http.MethodPost {
		c.ShouldBind(&cForm)
		err := r.Enable(cForm.Rulesets)
		if err != nil {
			return nil, errors.Wrap(err, "could not enable the ruleset")
		}

		for _, loadedHost := range h {
			err = r.Store(client, loadedHost.Name())
			if err != nil {
				return nil, errors.Wrap(err, "could not load the ruleset")
			}
		}
	} else { // GET action, load rulesets from 1st host

		r, err = ruleset.Load(client, h[0].Name())
		if err != nil {
			return nil, errors.Wrap(err, "could not load the ruleset")
		}

	}

	return r, nil
}

func NewClusterHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterName := c.Param("name")

		cluster, err := cluster.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		filter_query := fmt.Sprintf("Meta[\"trento-ha-cluster\"] == \"%s\"", clusterName)
		loadedHosts, err := hosts.Load(client, filter_query, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		rules, _ := getRulesets(c, client, loadedHosts)

		c.HTML(http.StatusOK, "cluster.html.tmpl", gin.H{
			"Cluster":  cluster[clusterName],
			"Hosts":    loadedHosts,
			"Rulesets": rules,
		})
	}
}
