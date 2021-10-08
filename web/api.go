package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/web/models"
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
// @Param resource_type query string false "Filter by resource type"
// @Success 200 {object} []string
// @Failure 500 {object} map[string]string
// @Router /api/tags [get]
func ApiListTag(tagsService services.TagsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()
		resourceTypeFilter := query["resource_type"]

		tags, err := tagsService.GetAll(resourceTypeFilter...)
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
func ApiHostCreateTagHandler(client consul.Client, tagsService services.TagsService) gin.HandlerFunc {
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

		err = tagsService.Create(r.Tag, models.TagHostResourceType, name)
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
func ApiHostDeleteTagHandler(client consul.Client, tagsService services.TagsService) gin.HandlerFunc {
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

		err = tagsService.Delete(tag, models.TagHostResourceType, name)
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
func ApiClusterCreateTagHandler(client consul.Client, tagsService services.TagsService) gin.HandlerFunc {
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

		err = tagsService.Create(r.Tag, models.TagClusterResourceType, id)
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
// @Param id path string true "Cluster id"
// @Param tag path string true "Tag"
// @Success 204 {object} map[string]interface{}
// @Router /api/clusters/{id}/tags/{tag} [delete]
func ApiClusterDeleteTagHandler(client consul.Client, tagsService services.TagsService) gin.HandlerFunc {
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

		err = tagsService.Delete(tag, models.TagClusterResourceType, id)
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
// @Param sid path string true "SAPSystem id"
// @Param Body body JSONTag true "The tag to create"
// @Success 201 {object} JSONTag
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/sapsystems/{sid}/tags [post]
func ApiSAPSystemCreateTagHandler(sapSystemsService services.SAPSystemsService, tagsService services.TagsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")

		systemList, err := sapSystemsService.GetSAPSystemsBySid(sid)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if len(systemList) == 0 {
			_ = c.Error(NotFoundError("could not find system"))
			return
		}

		var r JSONTag

		err = c.BindJSON(&r)
		if err != nil {
			_ = c.Error(BadRequestError("problems parsing JSON"))
			return
		}

		err = tagsService.Create(r.Tag, models.TagSAPSystemResourceType, sid)
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
// @Param sid path string true "SAPSystem id"
// @Param tag path string true "Tag"
// @Success 204 {object} map[string]interface{}
// @Router /api/sapsystems/{sid}/tags/{tag} [delete]
func ApiSAPSystemDeleteTagHandler(sapSystemsService services.SAPSystemsService, tagsService services.TagsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")
		tag := c.Param("tag")

		systemList, err := sapSystemsService.GetSAPSystemsBySid(sid)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if len(systemList) == 0 {
			_ = c.Error(NotFoundError("could not find system"))
			return
		}

		err = tagsService.Delete(tag, models.TagSAPSystemResourceType, sid)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

// ApiDatabaseCreateTagHandler godoc
// @Summary Add tag to a HANA database
// @Accept json
// @Produce json
// @Param sid path string true "Database id"
// @Param Body body JSONTag true "The tag to create"
// @Success 201 {object} JSONTag
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/databases/{sid}/tags [post]
func ApiDatabaseCreateTagHandler(sapSystemsService services.SAPSystemsService, tagsService services.TagsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")

		systemList, err := sapSystemsService.GetSAPSystemsBySid(sid)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if len(systemList) == 0 {
			_ = c.Error(NotFoundError("could not find database"))
			return
		}

		var r JSONTag

		err = c.BindJSON(&r)
		if err != nil {
			_ = c.Error(BadRequestError("problems parsing JSON"))
			return
		}

		err = tagsService.Create(r.Tag, models.TagDatabaseResourceType, sid)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, &r)
	}
}

// ApiDatabaseDeleteTagHandler godoc
// @Summary Delete a specific tag that belongs to a HANA database
// @Accept json
// @Produce json
// @Param sid path string true "Database id"
// @Param tag path string true "Tag"
// @Success 204 {object} map[string]interface{}
// @Router /api/databases/{sid}/tags/{tag} [delete]
func ApiDatabaseDeleteTagHandler(sapSystemsService services.SAPSystemsService, tagsService services.TagsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("sid")
		tag := c.Param("tag")

		systemList, err := sapSystemsService.GetSAPSystemsBySid(sid)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if len(systemList) == 0 {
			_ = c.Error(NotFoundError("could not find database"))
			return
		}

		err = tagsService.Delete(tag, models.TagDatabaseResourceType, sid)
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
// @Param cluster_id path string true "Cluster Id"
// @Success 200 {object} map[string]interface{}
// @Error 500
// @Router /api/clusters/{cluster_id}/results [get]
func ApiClusterCheckResultsHandler(client consul.Client, s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clusterId := c.Param("cluster_id")

		checkResults, err := s.GetChecksResultAndMetadataByCluster(clusterId)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, checkResults)
	}
}
