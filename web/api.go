package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/internal/sapsystem"
	"github.com/trento-project/trento/internal/tags"
	"github.com/trento-project/trento/web/services"
)

func ApiPingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

type JSONTag struct {
	Tag string `json:"tag" binding:"required"`
}

// ApiListTag godoc
// @Summary List all the tags in the system
// @Accept json
// @Produce json
// @Param resourceType query string false "Filter by resource type"
// @Success 200 {object} []string
// @Failure 500 {object} map[string]string
// @Router /api/tags [get]
func ApiListTag(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()
		resourceTypeFilter := query["resource_type"]

		t := tags.NewTags(client)

		tags, err := t.GetAll(resourceTypeFilter...)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if tags == nil {
			c.JSON(http.StatusOK, []string{})
			return
		}

		c.JSON(http.StatusOK, tags)
	}
}

// ApiHostCreateTagHandler godoc
// @Summary Add tag to host
// @Accept json
// @Produce json
// @Param name path string true "Host name"
// @Param Body body JSONTag true "The tag to create"
// @Success 201 {object} JSONTag
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/hosts/{name}/tags [post]
func ApiHostCreateTagHandler(client consul.Client) gin.HandlerFunc {
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

		var r JSONTag

		err = c.BindJSON(&r)
		if err != nil {
			_ = c.Error(BadRequestError("unable to parse JSON body"))
			return
		}

		t := tags.NewTags(client)
		err = t.Create(r.Tag, tags.HostResourceType, name)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, &r)
	}
}

// ApiHostDeleteTagHandler godoc
// @Summary Delete a specific tag that belongs to a host
// @Accept json
// @Produce json
// @Param name path string true "Host name"
// @Param tag path string true "Tag"
// @Success 204 {object} map[string]interface{}
// @Router /api/hosts/{name}/tags/{tag} [delete]
func ApiHostDeleteTagHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		tag := c.Param("tag")

		catalogNode, _, err := client.Catalog().Node(name, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if catalogNode == nil {
			_ = c.Error(NotFoundError("could not find host"))
			return
		}

		t := tags.NewTags(client)
		err = t.Delete(tag, tags.HostResourceType, name)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

// ApiClusterCreateTagHandler godoc
// @Summary Add tag to Cluster
// @Accept json
// @Produce json
// @Param id path string true "Cluster id"
// @Param Body body JSONTag true "The tag to create"
// @Success 201 {object} JSONTag
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/clusters/{id}/tags [post]
func ApiClusterCreateTagHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		clusters, err := cluster.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if _, ok := clusters[id]; !ok {
			_ = c.Error(NotFoundError("could not find cluster"))
			return
		}

		var r JSONTag

		err = c.BindJSON(&r)
		if err != nil {
			_ = c.Error(BadRequestError("problems parsing JSON"))
			return
		}

		t := tags.NewTags(client)
		err = t.Create(r.Tag, tags.ClusterResourceType, id)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, &r)
	}
}

// ApiClusterDeleteTagHandler godoc
// @Summary Delete a specific tag that belongs to a cluster
// @Accept json
// @Produce json
// @Param cluster path string true "Cluster id"
// @Param tag path string true "Tag"
// @Success 204 {object} map[string]interface{}
// @Router /api/clusters/{name}/tags/{tag} [delete]
func ApiClusterDeleteTagHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		tag := c.Param("tag")

		clusters, err := cluster.Load(client)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if _, ok := clusters[id]; !ok {
			_ = c.Error(NotFoundError("could not find cluster"))
			return
		}

		t := tags.NewTags(client)
		err = t.Delete(tag, tags.ClusterResourceType, id)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

// ApiSAPSystemCreateTagHandler godoc
// @Summary Add tag to SAPSystem
// @Accept json
// @Produce json
// @Param id path string true "SAPSystem id"
// @Param Body body JSONTag true "The tag to create"
// @Success 201 {object} JSONTag
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/sapsystems/{id}/tags [post]
func ApiSAPSystemCreateTagHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")

		// TODO: store sapsystem outside hosts
		hostList, err := hosts.Load(client, "", nil, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var system *sapsystem.SAPSystem
		for _, h := range hostList {
			sapSystems, err := h.GetSAPSystems()
			if err != nil {
				_ = c.Error(err)
				return
			}

			for _, s := range sapSystems {
				if s.SID == sid {
					system = s
					break
				}
			}
		}

		if system == nil {
			_ = c.Error(NotFoundError("could not find system"))
			return
		}

		var r JSONTag

		err = c.BindJSON(&r)
		if err != nil {
			_ = c.Error(BadRequestError("problems parsing JSON"))
			return
		}

		t := tags.NewTags(client)
		err = t.Create(r.Tag, tags.SAPSystemResourceType, sid)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, &r)
	}
}

// ApiSAPSystemDeleteTagHandler godoc
// @Summary Delete a specific tag that belongs to a SAPSystem
// @Accept json
// @Produce json
// @Param cluster path string true "SAPSystem id"
// @Param tag path string true "Tag"
// @Success 204 {object} map[string]interface{}
// @Router /api/sapsystems/{name}/tags/{tag} [delete]
func ApiSAPSystemDeleteTagHandler(client consul.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")
		tag := c.Param("tag")

		// TODO: store sapsystem outside hosts
		hostList, err := hosts.Load(client, "", nil, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var system *sapsystem.SAPSystem
		for _, h := range hostList {
			sapSystems, err := h.GetSAPSystems()
			if err != nil {
				_ = c.Error(err)
				return
			}

			for _, s := range sapSystems {
				if s.SID == sid {
					system = s
					break
				}
			}
		}

		if system == nil {
			_ = c.Error(NotFoundError("could not find system"))
			return
		}

		t := tags.NewTags(client)
		err = t.Delete(tag, tags.SAPSystemResourceType, sid)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

// ApiCheckResultsHandler godoc
// @Summary Get a specific cluster's check results
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Error 500
// @Router /api/clusters/{id}/results
func ApiClusterCheckResultsHandler(client consul.Client, s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterId := c.Param("cluster_id")

		checkResults, err := s.GetChecksResultByCluster(clusterId)
		if err != nil {
			c.Error(err)
			return
		}

		checksCatalog, err := s.GetChecksCatalog()
		if err != nil {
			c.Error(err)
			return
		}

		for checkId, check := range checkResults.Checks {
			check.Group = checksCatalog[checkId].Group
			check.Description = checksCatalog[checkId].Description
			check.ID = checksCatalog[checkId].ID
		}

		c.JSON(http.StatusOK, checkResults)
	}
}
