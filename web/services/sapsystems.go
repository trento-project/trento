package services

import (
	"fmt"
	"net"

	"github.com/trento-project/trento/internal"
	"github.com/trento-project/trento/internal/consul"
	"github.com/trento-project/trento/internal/hosts"
	"github.com/trento-project/trento/internal/sapsystem"
)

//go:generate mockery --name=SAPSystemsService --inpackage --filename=sapsystems_mock.go

type SAPSystemsService interface {
	GetSAPSystems() (sapsystem.SAPSystemsList, error)
	GetSAPSystemsById(id string) (sapsystem.SAPSystemsList, error)
	GetSAPSystemsByType(systemType int) (sapsystem.SAPSystemsList, error)
	GetAttachedDatabasesById(id string) (sapsystem.SAPSystemsList, error)
}

type sapSystemsService struct {
	consul consul.Client
}

func NewSAPSystemsService(client consul.Client) SAPSystemsService {
	return &sapSystemsService{consul: client}
}

func (s *sapSystemsService) GetSAPSystems() (sapsystem.SAPSystemsList, error) {
	var sapSystemsList sapsystem.SAPSystemsList

	hostList, err := hosts.Load(s.consul, "", nil)
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

func (s *sapSystemsService) GetSAPSystemsById(id string) (sapsystem.SAPSystemsList, error) {
	var sapSystemsListBySid sapsystem.SAPSystemsList

	sapSystemsList, err := s.GetSAPSystems()
	if err != nil {
		return nil, err
	}

	for _, s := range sapSystemsList {
		if s.Id == id {
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

func (s *sapSystemsService) GetAttachedDatabasesById(id string) (sapsystem.SAPSystemsList, error) {
	var sapDatabases sapsystem.SAPSystemsList

	// Find current SAP system
	sapSystemsList, err := s.GetSAPSystemsById(id)
	if err != nil {
		return sapDatabases, err
	}

	if len(sapSystemsList) == 0 {
		return sapDatabases, fmt.Errorf("system with %s not found", id)
	}

	dbs, dbsfound := sapSystemsList[0].Profile["dbs"]
	hdb, hdbfound := dbs.(map[string]interface{})["hdb"]
	dbname, dbnamefound := hdb.(map[string]interface{})["dbname"]
	sapdbhost, sapdbhostfound := sapSystemsList[0].Profile["SAPDBHOST"].(string)

	if !dbsfound || !hdbfound || !dbnamefound || !sapdbhostfound {
		return sapDatabases, fmt.Errorf("database data not available in the system %s", id)
	}

	// Get IP address if the sapdbhost is configured with a name
	addresses, err := net.LookupIP(sapdbhost)
	if err != nil {
		return sapDatabases, err
	}

	// Find Nodes with the database address
	query := fmt.Sprintf("Meta[\"trento-host-ip-addresses\"] contains \"%s\"", addresses[0])
	hostList, err := hosts.Load(s.consul, query, nil)
	if err != nil {
		return sapDatabases, err
	}

	if len(hostList) == 0 {
		return sapDatabases, nil
	}

	hostListNames := make([]string, 0)
	for _, h := range hostList {
		hostListNames = append(hostListNames, h.Name())
	}

	// Get all databases
	databaseList, err := s.GetSAPSystemsByType(sapsystem.Database)
	if err != nil {
		return sapDatabases, err
	}

	var databaseId string

	// Find the database that matches the database name and has the database address
	for _, database := range databaseList {
		for _, tenants := range database.Databases {
			if tenants.Database == dbname && internal.Contains(hostListNames, tenants.Host) {
				databaseId = database.Id
				break
			}
		}
	}

	if databaseId == "" {
		return sapDatabases, nil
	}

	// Get all the databases instances with the found id
	attachedDatabaseList, err := s.GetSAPSystemsById(databaseId)
	if err != nil {
		return sapDatabases, err
	}

	return attachedDatabaseList, nil
}
