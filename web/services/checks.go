package services

import (
	"fmt"
	"github.com/mitchellh/mapstructure"

	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services/ara"

	log "github.com/sirupsen/logrus"
)

//go:generate mockery --name=ChecksService

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
	GetChecksCatalog() (models.CheckList, error)
	GetChecksCatalogByGroup() (models.GroupedCheckList, error)
	GetChecksResult() (map[string]*models.Results, error)
	GetChecksResultByCluster(clusterId string) (*models.Results, error)
	GetChecksResultAndMetadataByCluster(clusterId string) (*models.ClusterCheckResults, error)
	GetAggregatedChecksResultByHost(clusterId string) (map[string]*AggregatedCheckData, error)
	GetAggregatedChecksResultByCluster(clusterId string) (*AggregatedCheckData, error)
}

type checksService struct {
	araService ara.AraService
}

func NewChecksService(araService ara.AraService) ChecksService {
	return &checksService{araService: araService}
}

func (c *checksService) GetChecksCatalog() (models.CheckList, error) {
	var checkData = models.CheckData{}

	records, err := c.araService.GetRecordList("key=trento-metadata&order=-id")
	if err != nil {
		return nil, err
	}

	if len(records.Results) == 0 {
		return nil, fmt.Errorf("Couldn't find any check catalog record. Check if the runner component is running")
	}

	record, err := c.araService.GetRecord(records.Results[0].ID)
	if err != nil {
		return nil, err
	}

	log.Debug(record.Value)

	mapstructure.Decode(record.Value, &checkData)

	return checkData.Metadata.Checks, nil
}

func (c *checksService) GetChecksCatalogByGroup() (models.GroupedCheckList, error) {
	groupedCheckMap := make(map[string]models.CheckList)

	checkList, err := c.GetChecksCatalog()
	if err != nil {
		return nil, err
	}

	for _, c := range checkList {
		if _, ok := groupedCheckMap[c.Group]; !ok {
			groupedCheckMap[c.Group] = models.CheckList{}
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

func (c *checksService) GetChecksResult() (map[string]*models.Results, error) {
	var checkData = models.CheckData{}

	records, err := c.araService.GetRecordList("key=trento-results&order=-id")
	if err != nil {
		return nil, err
	}

	if len(records.Results) == 0 {
		return nil, fmt.Errorf("Couldn't find any check result record. Check if the runner component is running")
	}

	record, err := c.araService.GetRecord(records.Results[0].ID)
	if err != nil {
		return nil, err
	}

	mapstructure.Decode(record.Value, &checkData)

	return checkData.Groups, nil
}

func (c *checksService) GetChecksResultByCluster(clusterId string) (*models.Results, error) {
	cResult, err := c.GetChecksResult()
	if err != nil {
		return nil, err
	}

	cResultByCluster, ok := cResult[clusterId]
	if !ok {
		return nil, fmt.Errorf("Cluster %s not found", clusterId)
	}

	return cResultByCluster, nil
}

func (c *checksService) GetChecksResultAndMetadataByCluster(clusterId string) (*models.ClusterCheckResults, error) {
	cResultByCluster, err := c.GetChecksResultByCluster(clusterId)
	if err != nil {
		return nil, err
	}

	checkList, err := c.GetChecksCatalog()
	if err != nil {
		return nil, err
	}

	resultSet := &models.ClusterCheckResults{}
	resultSet.Hosts = cResultByCluster.Hosts
	resultSet.Checks = []models.ClusterCheckResult{}

	for _, checkMeta := range checkList {
		for checkId, checkByHost := range cResultByCluster.Checks {
			if checkId == checkMeta.ID {
				current := models.ClusterCheckResult{
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

func (c *checksService) GetAggregatedChecksResultByHost(clusterId string) (map[string]*AggregatedCheckData, error) {
	cResultByCluster, err := c.GetChecksResultByCluster(clusterId)
	if err != nil {
		return nil, err
	}

	aCheckDataByHost := make(map[string]*AggregatedCheckData)

	for _, check := range cResultByCluster.Checks {
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

func (c *checksService) GetAggregatedChecksResultByCluster(clusterId string) (*AggregatedCheckData, error) {
	aCheckData := &AggregatedCheckData{}

	aCheckDataByHost, err := c.GetAggregatedChecksResultByHost(clusterId)
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
