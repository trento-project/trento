package web

import (
	"fmt"
	//"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	//consulApi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"
)

const KV_CLUSTERS_PATH string = "trento/clusters"

type Cluster struct {
	Name string
}

type ClusterList map[string]*Cluster

func NewClustersListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusters, err := loadClusters(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "clusters.html.tmpl", gin.H{
			"Clusters": clusters,
		})
	}
}

func loadClusters(client consul.Client) (ClusterList, error) {
	var clusters = ClusterList{}

	entries, _, err := client.KV().List(KV_CLUSTERS_PATH, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for Cluster KV values")
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Key, "clusters/") {
			continue
		}

		key_values := strings.Split(entry.Key, "/")

		if strings.HasSuffix(entry.Key, "/") {
			// 2 is used as Split creates a last empty entry
			last_key := key_values[len(key_values)-2]
			clusters[last_key] = &Cluster{}
			continue
		}

		cluster_id := key_values[len(key_values)-2]
		key := key_values[len(key_values)-1]
		// This could be done with a more automatic way in the future when we define the
		// Cluster and KV structure
		switch key {
		case "name":
			clusters[cluster_id].Name = string(entry.Value)
		}

	}
	return clusters, nil
}

func NewClusterHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cluster = Cluster{}

		cluster_name := c.Param("name")

		cluster_data, _, err := client.KV().List(fmt.Sprintf("%s/%s", KV_CLUSTERS_PATH, cluster_name), nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		for _, entry := range cluster_data {
			if strings.HasSuffix(entry.Key, "/") {
				continue
			}

			key_values := strings.Split(entry.Key, "/")
			key := key_values[len(key_values)-1]
			// This could be done with a more automatic way in the future when we define the
			// Cluster and KV structure
			switch key {
			case "name":
				cluster.Name = string(entry.Value)
			}
		}

		environments, err := loadEnvironments(client, fmt.Sprintf("Meta[\"trento-ha-cluster\"] == \"%s\"", cluster_name), nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "cluster.html.tmpl", gin.H{
			"Cluster":      cluster,
			"Environments": environments,
		})
	}
}
