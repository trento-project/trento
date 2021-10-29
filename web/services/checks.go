package services

import (
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/trento-project/trento/web/models"
)

//go:generate mockery --name=ChecksService --inpackage --filename=checks_mock.go
type AggregatedCheckData struct {
	PassingCount  int
	WarningCount  int
	CriticalCount int
}

func (a *AggregatedCheckData) String() string {
	if a.CriticalCount > 0 {
		return models.CheckCritical
	} else if a.WarningCount > 0 {
		return models.CheckWarning
	} else if a.PassingCount > 0 {
		return models.CheckPassing
	}

	return models.CheckUndefined
}

type ChecksService interface {
	// Check catalog services
	GetChecksCatalog() (models.ChecksCatalog, error)
	GetChecksCatalogByGroup() (models.GroupedCheckList, error)
	CreateChecksCatalogEntry(check *models.Check) error
	CreateChecksCatalog(checkList models.ChecksCatalog) error
	// Check result services
	GetChecksResultById(id string) (*models.Results, error)
	CreateChecksResultsById(id string, checkResults *models.Results) error
	GetChecksResultAndMetadataById(id string) (*models.ResultsAsList, error)
	GetAggregatedChecksResultByHost(id string) (map[string]*AggregatedCheckData, error)
	GetAggregatedChecksResultById(id string) (*AggregatedCheckData, error)
	// Selected checks services
	GetSelectedChecksById(id string) (models.SelectedChecks, error)
	CreateSelectedChecks(id string, selectedChecksList []string) error
	// Connection data services
	GetConnectionSettingsById(id string) (map[string]models.ConnectionSettings, error)
	GetConnectionSettingsByNode(node string) (models.ConnectionSettings, error)
	CreateConnectionSettings(node, cluster, user string) error
}

type checksService struct {
	db *gorm.DB
}

func NewChecksService(db *gorm.DB) *checksService {
	return &checksService{db: db}
}

/*
Checks catalog services
*/

func (c *checksService) GetChecksCatalog() (models.ChecksCatalog, error) {
	var checksRaw []*models.CheckRaw
	result := c.db.Order("payload->>'name'").Find(&checksRaw)
	if result.Error != nil {
		return nil, result.Error
	}

	var checksCatalog models.ChecksCatalog

	for _, checkRaw := range checksRaw {
		var check models.Check
		err := json.Unmarshal(checkRaw.Payload, &check)
		if err != nil {
			return nil, err
		}
		checksCatalog = append(checksCatalog, &check)
	}

	return checksCatalog, nil
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

	checkRaw := models.CheckRaw{ID: check.ID, Payload: checkJson}
	result := c.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&checkRaw)

	return result.Error
}

func (c *checksService) CreateChecksCatalog(checkList models.ChecksCatalog) error {
	for _, check := range checkList {
		result := c.CreateChecksCatalogEntry(check)
		if result != nil {
			return result
		}
	}

	return nil
}

/*
Checks result services
*/

func (c *checksService) GetChecksResultById(id string) (*models.Results, error) {
	var resultRaw models.CheckResultsRaw
	result := c.db.Where("group_id", id).Last(&resultRaw)

	if result.Error != nil {
		return nil, result.Error
	}

	var checkResult models.Results
	err := json.Unmarshal(resultRaw.Payload, &checkResult)
	return &checkResult, err
}

func (c *checksService) CreateChecksResultsById(id string, checkResults *models.Results) error {
	jsonData, err := json.Marshal(&checkResults)
	if err != nil {
		return err
	}

	event := models.CheckResultsRaw{GroupID: id, Payload: jsonData}
	result := c.db.Create(&event)

	return result.Error
}

func (c *checksService) GetChecksResultAndMetadataById(id string) (*models.ResultsAsList, error) {
	cResultById, err := c.GetChecksResultById(id)
	if err != nil {
		return nil, err
	}

	checkList, err := c.GetChecksCatalog()
	if err != nil {
		return nil, err
	}

	resultSet := &models.ResultsAsList{}
	resultSet.Hosts = cResultById.Hosts
	resultSet.Checks = []*models.ChecksByHost{}

	for _, checkMeta := range checkList {
		for checkId, checkByHost := range cResultById.Checks {
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

func (c *checksService) GetAggregatedChecksResultByHost(id string) (map[string]*AggregatedCheckData, error) {
	cResultById, err := c.GetChecksResultById(id)
	if err != nil {
		return nil, err
	}

	aCheckDataByHost := make(map[string]*AggregatedCheckData)

	for _, check := range cResultById.Checks {
		for hostName, host := range check.Hosts {
			if _, ok := aCheckDataByHost[hostName]; !ok {
				aCheckDataByHost[hostName] = &AggregatedCheckData{}
			}
			switch host.Result {
			case models.CheckCritical:
				aCheckDataByHost[hostName].CriticalCount += 1
			case models.CheckWarning:
				aCheckDataByHost[hostName].WarningCount += 1
			case models.CheckPassing:
				aCheckDataByHost[hostName].PassingCount += 1
			}
		}
	}

	return aCheckDataByHost, nil
}

func (c *checksService) GetAggregatedChecksResultById(id string) (*AggregatedCheckData, error) {
	aCheckData := &AggregatedCheckData{}

	aCheckDataByHost, err := c.GetAggregatedChecksResultByHost(id)
	if err != nil {
		return aCheckData, err
	}

	for _, aData := range aCheckDataByHost {
		aCheckData.CriticalCount += aData.CriticalCount
		aCheckData.WarningCount += aData.WarningCount
		aCheckData.PassingCount += aData.PassingCount
	}

	return aCheckData, nil
}

/*
Selected checks services
*/

func (c *checksService) GetSelectedChecksById(id string) (models.SelectedChecks, error) {
	var selectedChecks models.SelectedChecks

	result := c.db.Where("id", id).First(&selectedChecks)

	return selectedChecks, result.Error
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
