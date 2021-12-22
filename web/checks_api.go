package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"

	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

type JSONChecksSettings struct {
	SelectedChecks     []string          `json:"selected_checks" binding:"required"`
	ConnectionSettings map[string]string `json:"connection_settings" binding:"required"`
	Hostnames          []string          `json:"hostnames"`
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
	Premium        bool   `json:"premium,omitempty"`
}

type JSONChecksGroup struct {
	Group  string          `json:"group"`
	Checks []*models.Check `json:"checks"`
}

type JSONChecksGroupedCatalog []*JSONChecksGroup

type JSONChecksResult struct {
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

		for _, group := range checkGroups.OrderByName() {
			g := JSONChecksGroup{Group: group.Group, Checks: group.Checks}
			groupedCatalog = append(groupedCatalog, &g)
		}

		c.JSON(http.StatusOK, groupedCatalog)
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
				Premium:        checkData.Premium,
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

// ApiCheckResultsHandler godoc
// @Summary Get a specific cluster's check results
// @Produce json
// @Param cluster_id path string true "Cluster Id"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /clusters/{cluster_id}/results [get]
func ApiClusterCheckResultsHandler(s services.ChecksService) gin.HandlerFunc {
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

// ApiCreateChecksResultHandler godoc
// @Summary Create a checks result entry
// @Produce json
// @Param id path string true "Resource Id"
// @Param Body body JSONChecksResult true "Checks result"
// @Success 201 {object} JSONChecksResult
// @Failure 500 {object} map[string]string
// @Router /checks/{id}/results [post]
func ApiCreateChecksResultHandler(s services.ChecksService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var r JSONChecksResult

		id := c.Param("id")

		err := c.BindJSON(&r)
		if err != nil {
			_ = c.Error(BadRequestError("unable to parse JSON body"))
			return
		}

		var results models.ChecksResult
		// This is the easier way to decode the json format in the internal models
		mapstructure.Decode(r, &results)

		err = s.CreateChecksResultById(id, &results)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, &r)
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
func ApiCheckGetSettingsByIdHandler(s services.ClustersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		resourceId := c.Param("id")

		clusterSettings, err := s.GetClusterSettingsByID(resourceId)
		if err != nil {
			c.Error(err)
			return
		}

		if clusterSettings == nil {
			c.Error(NotFoundError("cluster not found"))
			return
		}

		resp := &JSONChecksSettings{
			SelectedChecks:     clusterSettings.SelectedChecks,
			ConnectionSettings: make(map[string]string),
		}

		for _, host := range clusterSettings.Hosts {
			resp.ConnectionSettings[host.Name] = host.User
			resp.Hostnames = append(resp.Hostnames, host.Name)
		}

		c.JSON(http.StatusOK, resp)
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
