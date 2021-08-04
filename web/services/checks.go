package services

import (
	"fmt"
	"github.com/mitchellh/mapstructure"

	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services/ara"

	log "github.com/sirupsen/logrus"
)

//go:generate mockery --name=ChecksService

type ChecksService interface {
	GetChecksCatalog() (map[string]*models.Check, error)
	GetChecksCatalogByGroup() (map[string]map[string]*models.Check, error)
	GetChecksResult() (map[string]*models.Results, error)
	GetChecksResultByCluster(clusterName string) (*models.Results, error)
}

type checksService struct {
	araService ara.AraService
}

func NewChecksService(araService ara.AraService) ChecksService {
	return &checksService{araService: araService}
}

func (c *checksService) GetChecksCatalog() (map[string]*models.Check, error) {
	var checkData = models.CheckData{}

	records, err := c.araService.GetRecordList("key=trento-metadata&order=-id")
	if err != nil {
		return nil, err
	}

	if len(records.Results) == 0 {
		return nil, nil
	}

	record, err := c.araService.GetRecord(records.Results[0].ID)
	if err != nil {
		return nil, err
	}

	log.Debug(record.Value)

	mapstructure.Decode(record.Value, &checkData)

	return checkData.Metadata.Checks, nil
}

func (c *checksService) GetChecksCatalogByGroup() (map[string]map[string]*models.Check, error) {
	groupedCheckList := make(map[string]map[string]*models.Check)

	checkList, err := c.GetChecksCatalog()
	if err != nil {
		return groupedCheckList, err
	}

	for cId, c := range checkList {
		extendedGroup := c.ExtendedGroupName()
		if _, ok := groupedCheckList[extendedGroup]; !ok {
			groupedCheckList[extendedGroup] = make(map[string]*models.Check)
		}
		groupedCheckList[extendedGroup][cId] = c
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
		return nil, nil
	}

	record, err := c.araService.GetRecord(records.Results[0].ID)
	if err != nil {
		return nil, err
	}

	mapstructure.Decode(record.Value, &checkData)

	return checkData.Groups, nil
}

func (c *checksService) GetChecksResultByCluster(clusterName string) (*models.Results, error) {
	cResult, err := c.GetChecksResult()
	if err != nil {
		return nil, err
	}

	cResultByHost, ok := cResult[clusterName]
	if !ok {
		return nil, fmt.Errorf("Cluster %s not found", clusterName)
	}

	return cResultByHost, nil
}
