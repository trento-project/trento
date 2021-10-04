package services

import (
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/internal/sapsystem"
)

//go:generate mockery --name=SAPSystemsService

type SAPSystemsService interface {
	GetSAPSystems() (sapsystem.SAPSystemsList, error)
	GetSAPSystemsBySid(sid string) (sapsystem.SAPSystemsList, error)
	GetSAPSystemsByType(systemType int) (sapsystem.SAPSystemsList, error)
}

type sapSystemsService struct {
	consul consul.Client
}

func NewSAPSystemsService(client consul.Client) SAPSystemsService {
	return &sapSystemsService{consul: client}
}

func (s *sapSystemsService) GetSAPSystems() (sapsystem.SAPSystemsList, error) {
	var sapSystemsList sapsystem.SAPSystemsList

	hostList, err := hosts.Load(s.consul, "", nil, nil)
	if err != nil {
		return nil, err
	}

	for _, h := range hostList {
		sapSystems, err := h.GetSAPSystems()
		if err != nil {
			return nil, err
		}

		for _, sapSystem := range sapSystems {
			sapSystemsList = append(sapSystemsList, sapSystem)
		}
	}

	return sapSystemsList, nil
}

func (s *sapSystemsService) GetSAPSystemsBySid(sid string) (sapsystem.SAPSystemsList, error) {
	var sapSystemsListBySid sapsystem.SAPSystemsList

	sapSystemsList, err := s.GetSAPSystems()
	if err != nil {
		return nil, err
	}

	for _, s := range sapSystemsList {
		if s.SID == sid {
			sapSystemsListBySid = append(sapSystemsListBySid, s)
		}
	}

	return sapSystemsListBySid, nil
}

func (s *sapSystemsService) GetSAPSystemsByType(systemType int) (sapsystem.SAPSystemsList, error) {
	var sapSystemsListByType sapsystem.SAPSystemsList

	sapSystemsList, err := s.GetSAPSystems()
	if err != nil {
		return nil, err
	}

	for _, sapSystem := range sapSystemsList {
		if sapSystem.Type == systemType {
			sapSystemsListByType = append(sapSystemsListByType, sapSystem)
		}
	}

	return sapSystemsListByType, nil
}
