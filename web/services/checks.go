package services

import (
	"encoding/json"
	"errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/trento-project/trento/web/datapipeline"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
)

const (
	partialChecksHealth = "config_checks"
)

//go:generate mockery --name=ChecksService --inpackage --filename=checks_mock.go

type ChecksService interface {
	// Check catalog services
	GetChecksCatalog() (models.ChecksCatalog, error)
	GetChecksCatalogByGroup() (models.GroupedCheckList, error)
	CreateChecksCatalogEntry(check *models.Check) error // seems to be never used
	CreateChecksCatalog(checkList models.ChecksCatalog) error
	// Check result services
	CreateChecksResult(checksResult *models.ChecksResult) error
	GetLastExecutionByGroup() ([]*models.ChecksResult, error)
	GetChecksResultByCluster(clusterId string) (*models.ChecksResult, error)
	GetChecksResultAndMetadataByCluster(clusterId string) (*models.ChecksResultAsList, error)
	GetAggregatedChecksResultByHost(clusterId string) (map[string]*models.AggregatedCheckData, error)
	GetAggregatedChecksResultByCluster(clusterId string) (*models.AggregatedCheckData, error)
	// Selected checks services
	GetSelectedChecksById(id string) (models.SelectedChecks, error)
	CreateSelectedChecks(id string, selectedChecksList []string) error
	// Connection data services
	GetConnectionSettingsById(id string) (map[string]models.ConnectionSettings, error)
	GetConnectionSettingsByNode(node string) (models.ConnectionSettings, error)
	CreateConnectionSettings(node, cluster, user string) error
}

type checksService struct {
	db                      *gorm.DB
	premiumDetectionService PremiumDetectionService
}

func NewChecksService(db *gorm.DB, premiumDetectionService PremiumDetectionService) *checksService {
	return &checksService{
		db:                      db,
		premiumDetectionService: premiumDetectionService,
	}
}

/*
Checks catalog services
*/

func (c *checksService) GetChecksCatalog() (models.ChecksCatalog, error) {
	var checksEntity entities.CheckList
	isPremiumActive, _ := c.premiumDetectionService.IsPremiumActive()

	var result *gorm.DB
	qb := c.db.Order("payload->>'name'")

	if isPremiumActive {
		result = qb.Find(&checksEntity)
	} else {
		result = qb.Find(&checksEntity, datatypes.JSONQuery("payload").Equals(false, "premium"))
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return checksEntity.ToModel()
}

func (c *checksService) GetChecksCatalogByGroup() (models.GroupedCheckList, error) {
	groupedCheckMap := make(map[string]models.ChecksCatalog)

	checkList, err := c.GetChecksCatalog()
	if err != nil {
		return nil, err
	}

	for _, c := range checkList {
		if _, ok := groupedCheckMap[c.Group]; !ok {
			groupedCheckMap[c.Group] = models.ChecksCatalog{}
		}
		groupedCheckMap[c.Group] = append(groupedCheckMap[c.Group], c)
	}

	groupedCheckList := make(models.GroupedCheckList, 0)

	for group, checks := range groupedCheckMap {
		g := &models.GroupedChecks{Group: group, Checks: checks}
		groupedCheckList = append(groupedCheckList, g)
	}

	return groupedCheckList, nil
}

func (c *checksService) CreateChecksCatalogEntry(check *models.Check) error {
	checkJson, err := json.Marshal(&check)
	if err != nil {
		return err
	}

	checkEntity := entities.Check{ID: check.ID, Payload: checkJson}
	result := c.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&checkEntity)

	return result.Error
}

func (c *checksService) CreateChecksCatalog(checkList models.ChecksCatalog) error {

	var checkEntityList entities.CheckList
	for _, check := range checkList {
		checkJson, err := json.Marshal(&check)
		if err != nil {
			return err
		}
		checkEntityList = append(checkEntityList, &entities.Check{ID: check.ID, Payload: checkJson})
	}

	result := c.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&checkEntityList)

	if result.Error != nil {
		return result.Error
	}

	// Remove old not updated checks
	return c.db.Not(&checkEntityList).Delete(entities.CheckList{}).Error
}

/*
Checks result services
*/

func (c *checksService) CreateChecksResult(checksResult *models.ChecksResult) error {
	jsonData, err := json.Marshal(&checksResult)
	if err != nil {
		return err
	}

	event := entities.ChecksResult{GroupID: checksResult.ID, Payload: jsonData}
	err = c.db.Create(&event).Error
	if err != nil {
		return err
	}

	// Project the current health state
	aggregatedHealth, err := c.GetAggregatedChecksResultByCluster(checksResult.ID)
	if err != nil {
		return err
	}

	err = datapipeline.ProjectHealth(
		c.db, checksResult.ID, partialChecksHealth, aggregatedHealth.String())
	if err != nil {
		return err
	}

	return nil
}

func (c *checksService) GetLastExecutionByGroup() ([]*models.ChecksResult, error) {
	var checksResults []entities.ChecksResult

	err := c.db.Where("(group_id, created_at) IN (?)", c.db.Model(&entities.ChecksResult{}).
		Select("group_id, max(created_at)").
		Group("group_id")).Order("id").Find(&checksResults).Error
	if err != nil {
		return nil, err
	}

	var checksResultModels []*models.ChecksResult
	for _, checksResult := range checksResults {
		cModel, _ := checksResult.ToModel()
		checksResultModels = append(checksResultModels, cModel)
	}

	return checksResultModels, nil
}

func (c *checksService) GetChecksResultByCluster(clusterId string) (*models.ChecksResult, error) {
	var checksResult entities.ChecksResult
	result := c.db.Where("group_id", clusterId).Last(&checksResult)

	if result.Error != nil {
		return nil, result.Error
	}

	return checksResult.ToModel()
}

func (c *checksService) GetChecksResultAndMetadataByCluster(clusterId string) (*models.ChecksResultAsList, error) {
	cResultByCluster, err := c.GetChecksResultByCluster(clusterId)
	if err != nil {
		return nil, err
	}

	checkList, err := c.GetChecksCatalog()
	if err != nil {
		return nil, err
	}

	resultSet := &models.ChecksResultAsList{}
	resultSet.Hosts = cResultByCluster.Hosts
	resultSet.Checks = []*models.ChecksByHost{}

	for _, checkMeta := range checkList {
		for checkId, checkByHost := range cResultByCluster.Checks {
			if checkId == checkMeta.ID {
				current := &models.ChecksByHost{
					Group:       checkMeta.Group,
					Description: checkMeta.Description,
					Hosts:       checkByHost.Hosts,
					ID:          checkId,
				}
				resultSet.Checks = append(resultSet.Checks, current)
				continue
			}
		}
	}

	return resultSet, nil
}

func (c *checksService) GetAggregatedChecksResultByHost(clusterId string) (map[string]*models.AggregatedCheckData, error) {
	cResultByCluster, err := c.GetChecksResultByCluster(clusterId)
	if err != nil {
		return nil, err
	}

	return cResultByCluster.GetAggregatedChecksResultByHost(), nil
}

func (c *checksService) GetAggregatedChecksResultByCluster(clusterId string) (*models.AggregatedCheckData, error) {
	cResultByCluster, err := c.GetChecksResultByCluster(clusterId)
	if err != nil {
		return nil, err
	}

	return cResultByCluster.GetAggregatedChecksResultByCluster(), nil
}

/*
Selected checks services
*/

func (c *checksService) GetSelectedChecksById(id string) (models.SelectedChecks, error) {
	selectedChecks := models.SelectedChecks{
		ID:             "",
		SelectedChecks: []string{},
	}

	catalog, err := c.GetChecksCatalog()
	if err != nil {
		return selectedChecks, err
	}

	err = c.db.Where("id", id).First(&selectedChecks).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return selectedChecks, nil
		}
		return selectedChecks, err
	}

	set := make(map[string]struct{})
	filteredChecks := []string{}

	for _, availableCheck := range catalog {
		set[availableCheck.ID] = struct{}{}
	}

	for _, s := range selectedChecks.SelectedChecks {
		if _, ok := set[s]; ok {
			filteredChecks = append(filteredChecks, s)
		}
	}

	selectedChecks.SelectedChecks = filteredChecks

	return selectedChecks, err
}

func (c *checksService) CreateSelectedChecks(id string, selectedChecksList []string) error {
	selectedChecks := models.SelectedChecks{
		ID:             id,
		SelectedChecks: selectedChecksList,
	}

	result := c.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&selectedChecks)

	return result.Error
}

/*
Checks connection user services
*/

func (c *checksService) GetConnectionSettingsByNode(node string) (models.ConnectionSettings, error) {
	var connUser models.ConnectionSettings

	result := c.db.Where("node", node).First(&connUser)

	return connUser, result.Error
}

func (c *checksService) GetConnectionSettingsById(id string) (map[string]models.ConnectionSettings, error) {
	var connUsersList []models.ConnectionSettings

	result := c.db.Where("id", id).Find(&connUsersList)

	if result.Error != nil {
		return nil, result.Error
	}

	connUsersMap := make(map[string]models.ConnectionSettings)
	for _, user := range connUsersList {
		connUsersMap[user.Node] = user
	}

	return connUsersMap, nil
}

func (c *checksService) CreateConnectionSettings(id, node, user string) error {
	connUser := models.ConnectionSettings{
		ID:   id,
		Node: node,
		User: user,
	}

	result := c.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&connUser)

	return result.Error
}
