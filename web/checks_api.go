package web

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"

	"github.com/trento-project/trento/internal/cluster"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

type JSONChecksSettings struct {
	SelectedChecks     []string          `json:"selected_checks" binding:"required"`
	ConnectionSettings map[string]string `json:"connection_settings" binding:"required"`
	Hosts              Nodes             `json:"hosts"`
}

type JSONChecksCatalog []*JSONCheck

type JSONCheck struct {
	ID             string `json:"id,omitempty" binding:"required"`
	Name           string `json:"name,omitempty" binding:"required"`
	Group          string `json:"group,omitempty" binding:"required"`
	Description    string `json:"description,omitempty"`
	Remediation    string `json:"remediation,omitempty"`
	Implementation string `json:"implementation,omitempty"`
	Labels         string `json:"labels,omitempty"`
}

type JSONChecksGroup struct {
	Group  string          `json:"group"`
	Checks []*models.Check `json:"checks"`
}

type JSONChecksGroupedCatalog []*JSONChecksGroup

// ApiCheckCatalogHandler godoc
// @Summary Get the whole checks' catalog
// @Produce json
// @Success 200 {object} JSONChecksGroupedCatalog
// @Error 500
// @Router /checks/catalog [get]
func ApiChecksCatalogHandler(s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var groupedCatalog JSONChecksGroupedCatalog

		checkGroups, err := s.GetChecksCatalogByGroup()
		if err != nil {
			c.Error(err)
			return
		}

		for _, group := range checkGroups {
			g := JSONChecksGroup{Group: group.Group, Checks: group.Checks}
			groupedCatalog = append(groupedCatalog, &g)
		}

		sort.SliceStable(groupedCatalog, func(i, j int) bool {
			return groupedCatalog[i].Group < groupedCatalog[j].Group
		})

		c.JSON(http.StatusOK, groupedCatalog)
	}
}

// ApiCheckResultsHandler godoc
// @Summary Get a specific cluster's check results
// @Produce json
// @Param cluster_id path string true "Cluster Id"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /clusters/{cluster_id}/results [get]
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

// ApiCreateChecksCatalogHandler godoc
// @Summary Create/Updates the checks catalog
// @Produce json
// @Param Body body JSONChecksCatalog true "Checks catalog"
// @Success 200 {object} JSONChecksCatalog
// @Failure 500 {object} map[string]string
// @Router /checks/catalog [put]
func ApiCreateChecksCatalogHandler(s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {

		var r JSONChecksCatalog

		err := c.BindJSON(&r)
		if err != nil {
			_ = c.Error(BadRequestError("unable to parse JSON body"))
			return
		}

		var catalog models.ChecksCatalog

		for _, checkData := range r {
			newCheck := &models.Check{
				ID:             checkData.ID,
				Name:           checkData.Name,
				Group:          checkData.Group,
				Description:    checkData.Description,
				Remediation:    checkData.Remediation,
				Implementation: checkData.Implementation,
				Labels:         checkData.Labels,
			}
			catalog = append(catalog, newCheck)
		}

		err = s.CreateChecksCatalog(catalog)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, &r)
	}
}

// ApiCheckGetSettingsByIdHandler godoc
// @Summary Get the check settings
// @Accept json
// @Produce json
// @Param id path string true "Resource id"
// @Success 200 {object} JSONChecksSettings
// @Failure 404 {object} map[string]string
// @Router /checks/{id}/settings [get]
func ApiCheckGetSettingsByIdHandler(consul consul.Client, s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		resourceId := c.Param("id")

		// TODO: this has absolutely to be refactored once we've got the hosts service
		clusters, err := cluster.Load(consul)
		if err != nil {
			_ = c.Error(err)
			return
		}

		clusterItem, ok := clusters[resourceId]
		if !ok {
			_ = c.Error(NotFoundError("could not find cluster"))
			return
		}

		filterQuery := fmt.Sprintf("Meta[\"trento-ha-cluster-id\"] == \"%s\"", resourceId)
		hosts, err := hosts.Load(consul, filterQuery, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		nodes := NewNodes(s, clusterItem, hosts)

		selectedChecks, err := s.GetSelectedChecksById(resourceId)
		if err != nil {
			log.Error(err)
		}

		connSettings, err := s.GetConnectionSettingsById(resourceId)
		if err != nil {
			_ = c.Error(NotFoundError("could not find connection settings"))
			return
		}

		var jsonCheckSetting JSONChecksSettings
		jsonCheckSetting.ConnectionSettings = make(map[string]string)

		if len(selectedChecks.SelectedChecks) == 0 {
			jsonCheckSetting.SelectedChecks = make([]string, 0)

		} else {
			jsonCheckSetting.SelectedChecks = selectedChecks.SelectedChecks
		}

		jsonCheckSetting.Hosts = nodes
		for node, settings := range connSettings {
			jsonCheckSetting.ConnectionSettings[node] = settings.User
		}
		c.JSON(http.StatusOK, jsonCheckSetting)
	}
}

// ApiCheckCreateSettingsByIdHandler godoc
// @Summary Create the check settings
// @Accept json
// @Produce json
// @Param id path string true "Resource id"
// @Param Body body JSONChecksSettings true "Checks settings"
// @Success 201 {object} JSONChecksSettings
// @Failure 500 {object} map[string]string
// @Router /checks/{id}/settings [post]
func ApiCheckCreateSettingsByIdHandler(s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		resourceId := c.Param("id")

		var r JSONChecksSettings

		err := c.BindJSON(&r)
		if err != nil {
			_ = c.Error(BadRequestError("unable to parse JSON body"))
			return
		}

		err = s.CreateSelectedChecks(resourceId, r.SelectedChecks)
		if err != nil {
			_ = c.Error(err)
			return
		}

		for node, user := range r.ConnectionSettings {
			err = s.CreateConnectionSettings(resourceId, node, user)
			if err != nil {
				_ = c.Error(err)
				return
			}
		}

		c.JSON(http.StatusCreated, &r)
	}
}
