package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

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
// @Router /tags [get]
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
// @Param id path string true "Host id"
// @Param Body body JSONTag true "The tag to create"
// @Success 201 {object} JSONTag
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hosts/{id}/tags [post]
func ApiHostCreateTagHandler(hostsService services.HostsService, tagsService services.TagsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		host, err := hostsService.GetByID(id)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if host == nil {
			_ = c.Error(NotFoundError("could not find host"))
			return
		}

		var r JSONTag

		err = c.BindJSON(&r)
		if err != nil {
			_ = c.Error(BadRequestError("unable to parse JSON body"))
			return
		}

		err = tagsService.Create(r.Tag, models.TagHostResourceType, id)
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
// @Param id path string true "Host id"
// @Param tag path string true "Tag"
// @Success 204 {object} map[string]interface{}
// @Router /hosts/{id}/tags/{tag} [delete]
func ApiHostDeleteTagHandler(hostsService services.HostsService, tagsService services.TagsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		tag := c.Param("tag")

		host, err := hostsService.GetByID(id)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if host == nil {
			_ = c.Error(NotFoundError("could not find host"))
			return
		}

		err = tagsService.Delete(tag, models.TagHostResourceType, id)
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
// @Router /clusters/{id}/tags [post]
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
// @Router /clusters/{id}/tags/{tag} [delete]
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
// @Param id path string true "SAPSystem id"
// @Param Body body JSONTag true "The tag to create"
// @Success 201 {object} JSONTag
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sapsystems/{id}/tags [post]
func ApiSAPSystemCreateTagHandler(sapSystemsService services.SAPSystemsService, tagsService services.TagsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		systemList, err := sapSystemsService.GetSAPSystemsById(id)
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

		err = tagsService.Create(r.Tag, models.TagSAPSystemResourceType, id)
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
// @Param id path string true "SAPSystem id"
// @Param tag path string true "Tag"
// @Success 204 {object} map[string]interface{}
// @Router /sapsystems/{id}/tags/{tag} [delete]
func ApiSAPSystemDeleteTagHandler(sapSystemsService services.SAPSystemsService, tagsService services.TagsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		tag := c.Param("tag")

		systemList, err := sapSystemsService.GetSAPSystemsById(id)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if len(systemList) == 0 {
			_ = c.Error(NotFoundError("could not find system"))
			return
		}

		err = tagsService.Delete(tag, models.TagSAPSystemResourceType, id)
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
// @Param id path string true "Database id"
// @Param Body body JSONTag true "The tag to create"
// @Success 201 {object} JSONTag
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /databases/{id}/tags [post]
func ApiDatabaseCreateTagHandler(sapSystemsService services.SAPSystemsService, tagsService services.TagsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		systemList, err := sapSystemsService.GetSAPSystemsById(id)
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

		err = tagsService.Create(r.Tag, models.TagDatabaseResourceType, id)
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
// @Param id path string true "Database id"
// @Param tag path string true "Tag"
// @Success 204 {object} map[string]interface{}
// @Router /databases/{id}/tags/{tag} [delete]
func ApiDatabaseDeleteTagHandler(sapSystemsService services.SAPSystemsService, tagsService services.TagsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		tag := c.Param("tag")

		systemList, err := sapSystemsService.GetSAPSystemsById(id)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if len(systemList) == 0 {
			_ = c.Error(NotFoundError("could not find database"))
			return
		}

		err = tagsService.Delete(tag, models.TagDatabaseResourceType, id)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
