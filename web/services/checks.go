package services

import (
	"github.com/mitchellh/mapstructure"

	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services/ara"
)

//go:generate mockery --name=ChecksService

type ChecksService interface {
	GetChecksCatalog() ([]*models.Check, error)
}

type checksService struct {
	araService ara.AraService
}

func NewChecksService(araService ara.AraService) ChecksService {
	return &checksService{araService: araService}
}

func (c *checksService) GetChecksCatalog() ([]*models.Check, error) {
	checkList := []*models.Check{}

	records, err := c.araService.GetRecordList("key=trento-metadata&order=-id")
	if err != nil {
		return checkList, err
	}

	if len(records.Results) == 0 {
		return checkList, nil
	}

	record, err := c.araService.GetRecord(records.Results[0].ID)
	if err != nil {
		return checkList, err
	}

	mapstructure.Decode(record.Value, &checkList)

	return checkList, nil
}
