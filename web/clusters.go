package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
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

func NewClusterHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterName := c.Param("name")

		clusters, err := cluster.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}
		cluster := clusters[clusterName]

		filter_query := fmt.Sprintf("Meta[\"trento-ha-cluster\"] == \"%s\"", clusterName)
		hosts, err := hosts.Load(client, filter_query, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "cluster.html.tmpl", gin.H{
			"Cluster": cluster,
			"Hosts":   hosts,
		})
	}
}
