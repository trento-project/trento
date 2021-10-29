package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"

	"github.com/trento-project/trento/internal/consul"

	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

type JSONChecksSettings struct {
	SelectedChecks     []string          `json:"selected_checks" binding:"required"`
	ConnectionSettings map[string]string `json:"connection_settings" binding:"required"`
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

type JSONChecksResults struct {
	Hosts  map[string]*JSONHosts       `json:"hosts,omitempty" binding:"required"`
	Checks map[string]*JSONCheckResult `json:"checks,omitempty" binding:"required"`
}

type JSONHosts struct {
	Result    string `json:"result,omitempty"`
	Reachable bool   `json:"reachable,omitempty"`
	Msg       string `json:"msg,omitempty"`
}

type JSONCheckResult struct {
	ID          string                `json:"id,omitempty"`
	Hosts       map[string]*JSONHosts `json:"hosts,omitempty"`
	Group       string                `json:"group,omitempty"`
	Description string                `json:"description,omitempty"`
}

// ApiCheckResultsHandler godoc
// @Summary Get a specific cluster's check results
// @Produce json
// @Param cluster_id path string true "Cluster Id"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
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

// ApiCreateChecksResultstaHandler godoc
// @Summary Create a checks results entry
// @Produce json
// @Param id path string true "Resource Id"
// @Param Body body JSONChecksResults true "Checks results"
// @Success 201 {object} JSONChecksResults
// @Error 500
// @Router /api/checks/{id}/results [post]
func ApiCreateChecksResultstaHandler(s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var r JSONChecksResults

		id := c.Param("id")

		err := c.BindJSON(&r)
		if err != nil {
			_ = c.Error(BadRequestError("unable to parse JSON body"))
			return
		}

		var results models.Results
		// This is the easier way to decode the json format in the internal models
		mapstructure.Decode(r, &results)

		err = s.CreateChecksResultsById(id, &results)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, &r)
	}
}

// ApiCreateChecksCatalogHandler godoc
// @Summary Create/Updates the checks catalog
// @Produce json
// @Param Body body JSONChecksCatalog true "Checks catalog"
// @Success 200 {object} JSONChecksCatalog
// @Failure 500 {object} map[string]string
// @Router /api/checks/catalog [put]
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
// @Router /api/checks/{id}/settings [get]
func ApiCheckGetSettingsByIdHandler(s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		selectedChecks, err := s.GetSelectedChecksById(id)
		if err != nil {
			_ = c.Error(NotFoundError("could not find check selection"))
			return
		}

		connSettings, err := s.GetConnectionSettingsById(id)
		if err != nil {
			_ = c.Error(NotFoundError("could not find connection settings"))
			return
		}

		var jsonCheckSetting JSONChecksSettings
		jsonCheckSetting.ConnectionSettings = make(map[string]string)

		jsonCheckSetting.SelectedChecks = selectedChecks.SelectedChecks
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
// @Router /api/checks/{id}/settings [post]
func ApiCheckCreateSettingsByIdHandler(s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var r JSONChecksSettings

		err := c.BindJSON(&r)
		if err != nil {
			_ = c.Error(BadRequestError("unable to parse JSON body"))
			return
		}

		err = s.CreateSelectedChecks(id, r.SelectedChecks)
		if err != nil {
			_ = c.Error(err)
			return
		}

		for node, user := range r.ConnectionSettings {
			err = s.CreateConnectionSettings(id, node, user)
			if err != nil {
				_ = c.Error(err)
				return
			}
		}

		c.JSON(http.StatusCreated, &r)
	}
}
