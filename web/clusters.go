package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

func NewClustersHealthContainer(clusterList models.ClusterList) *HealthContainer {
	h := &HealthContainer{}
	for _, c := range clusterList {
		switch c.Health {
		case models.CheckPassing:
			h.PassingCount += 1
		case models.CheckWarning:
			h.WarningCount += 1
		case models.CheckCritical:
			h.CriticalCount += 1
		}
	}
	return h
}

func NewClusterListHandler(clustersService services.ClustersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()

		clustersFilter := &services.ClustersFilter{
			Name:        query["name"],
			SIDs:        query["sids"],
			ClusterType: query["cluster_type"],
			Health:      query["health"],
			Tags:        query["tags"],
		}

		pageNumber, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil {
			pageNumber = 1
		}
		pageSize, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
		if err != nil {
			pageSize = 10
		}

		page := &services.Page{
			Number: pageNumber,
			Size:   pageSize,
		}

		clusterList, err := clustersService.GetAll(clustersFilter, page)
		if err != nil {
			_ = c.Error(err)
			return
		}

		filterClusterTypes, err := clustersService.GetAllClusterTypes()
		if err != nil {
			_ = c.Error(err)
			return
		}

		filterSIDs, err := clustersService.GetAllSIDs()
		if err != nil {
			_ = c.Error(err)
			return
		}

		filterTags, err := clustersService.GetAllTags()
		if err != nil {
			_ = c.Error(err)
			return
		}

		healthContainer := NewClustersHealthContainer(clusterList)
		healthContainer.Layout = "horizontal"

		count, err := clustersService.GetCount()
		if err != nil {
			_ = c.Error(err)
			return
		}
		pagination := NewPagination(count, pageNumber, pageSize)

		c.HTML(http.StatusOK, "clusters.html.tmpl", gin.H{
			"ClustersTable":      clusterList,
			"AppliedFilters":     query,
			"FilterClusterTypes": filterClusterTypes,
			"FilterSIDs":         filterSIDs,
			"FilterTags":         filterTags,
			"Pagination":         pagination,
			"HealthContainer":    healthContainer,
		})
	}
}

func NewClusterGenericHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterId := c.Param("id")

		clusters, err := cluster.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		clusterItem, ok := clusters[clusterId]
		if !ok {
			_ = c.Error(NotFoundError("could not find cluster"))
			return
		}

		filterQuery := fmt.Sprintf("Meta[\"trento-ha-cluster-id\"] == \"%s\"", clusterId)
		hosts, err := hosts.Load(client, filterQuery, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "cluster_generic.html.tmpl", gin.H{
			"Cluster": clusterItem,
			"Hosts":   hosts,
		})
	}
}

func NewClusterHandler(clusterService services.ClustersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterID := c.Param("id")

		cluster, err := clusterService.GetByID(clusterID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if cluster == nil {
			_ = c.Error(NotFoundError("could not find cluster"))
			return
		}

		if cluster.ClusterType == models.ClusterTypeUnknown {
			c.Redirect(http.StatusFound, fmt.Sprintf("/clusters/%s/generic", clusterID))
		}

		hContainer := &HealthContainer{
			CriticalCount: cluster.CriticalCount,
			WarningCount:  cluster.WarningCount,
			PassingCount:  cluster.PassingCount,
			Layout:        "vertical",
		}

		c.HTML(http.StatusOK, "cluster_hana.html.tmpl", gin.H{
			"Cluster":         cluster,
			"HealthContainer": hContainer,
			"Alerts":          GetAlerts(c),
		})
	}
}
