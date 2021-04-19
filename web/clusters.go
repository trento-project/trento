package web

import (
	"fmt"
	//"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/trento-project/trento/internal/consul"
)

type Cluster struct {
	Name string
}

type ClusterList map[string]*Cluster

func NewClustersListHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusters, err := loadClusters(client, "")
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
		cluster_name := c.Param("name")

		cluster, err := loadClusters(client, cluster_name)
		if err != nil {
			_ = c.Error(err)
			return
		}

		filter_query := fmt.Sprintf("Meta[\"trento-ha-cluster\"] == \"%s\"", cluster_name)
		hosts, err := loadHosts(client, filter_query, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "cluster.html.tmpl", gin.H{
			"Cluster": cluster[cluster_name],
			"Hosts":   hosts,
		})
	}
}

func loadClusters(client consul.Client, cluster_name string) (ClusterList, error) {
	var clusters = ClusterList{}

	kv_path := fmt.Sprintf("%s/%s", consul.KvClustersPath, cluster_name)

	entries, _, err := client.KV().List(kv_path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not query Consul for Cluster KV values")
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Key, "clusters/") {
			continue
		}

		key_values := strings.Split(entry.Key, "/")
		// 2 is used as Split creates a last empty entry
		cluster_id := key_values[len(key_values)-2]

		if strings.HasSuffix(entry.Key, "/") {
			clusters[cluster_id] = &Cluster{Name: cluster_id}
			continue
		}

		//value := key_values[len(key_values)-1]
		// This could be done with a more automatic way in the future when we define the
		// Cluster and KV structure
		//switch value {
		//case "name":
		//	clusters[cluster_id].Name = string(entry.Value)
		//}

	}
	return clusters, nil
}
