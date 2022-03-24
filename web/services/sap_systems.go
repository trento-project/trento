package services

import (
	"errors"
	"fmt"
	"net"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
)

//go:generate mockery --name=SAPSystemsService --inpackage --filename=sap_systems_service_mock.go
type SAPSystemsService interface {
	GetAllApplications(filter *SAPSystemFilter, page *Page) (models.SAPSystemList, error)
	GetAllDatabases(filter *SAPSystemFilter, page *Page) (models.SAPSystemList, error)
	GetByID(ID string) (*models.SAPSystem, error)
	GetApplicationsCount() (int, error)
	GetDatabasesCount() (int, error)
	GetAllApplicationsSIDs() ([]string, error)
	GetAllDatabasesSIDs() ([]string, error)
	GetAllApplicationsTags() ([]string, error)
	GetAllDatabasesTags() ([]string, error)
}

type SAPSystemFilter struct {
	Tags []string
	SIDs []string
}

type sapSystemsService struct {
	db *gorm.DB
}

func NewSAPSystemsService(db *gorm.DB) *sapSystemsService {
	return &sapSystemsService{db}
}

func (s *sapSystemsService) GetAllApplications(filter *SAPSystemFilter, page *Page) (models.SAPSystemList, error) {
	return s.getAllByType(models.SAPSystemTypeApplication, models.TagSAPSystemResourceType, filter, page)
}

func (s *sapSystemsService) GetAllDatabases(filter *SAPSystemFilter, page *Page) (models.SAPSystemList, error) {
	return s.getAllByType(models.SAPSystemTypeDatabase, models.TagDatabaseResourceType, filter, page)
}
func (s *sapSystemsService) GetByID(ID string) (*models.SAPSystem, error) {
	var instances entities.SAPSystemInstances

	err := s.db.
		Where("id = ?", ID).
		Order("sid, instance_number, system_replication, id").
		Find(&instances).
		Error

	if err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, nil
	}

	return instances.ToModel()[0], nil
}

func (s *sapSystemsService) GetApplicationsCount() (int, error) {
	var count int64

	err := s.db.
		Model(&entities.SAPSystemInstance{}).
		Distinct("id").
		Group("type").
		Having("type = ?", models.SAPSystemTypeApplication).
		Count(&count).
		Error

	return int(count), err
}

func (s *sapSystemsService) GetDatabasesCount() (int, error) {
	var count int64

	err := s.db.
		Model(&entities.SAPSystemInstance{}).
		Distinct("id").
		Group("type").
		Having("type = ?", models.SAPSystemTypeDatabase).
		Count(&count).
		Error

	return int(count), err
}

func (s *sapSystemsService) GetAllApplicationsSIDs() ([]string, error) {
	var sids pq.StringArray

	err := s.db.
		Model(&entities.SAPSystemInstance{}).
		Statement.Where("type = ?", models.SAPSystemTypeApplication).
		Distinct().
		Pluck("sid", &sids).
		Error

	if err != nil {
		return nil, err
	}

	return []string(sids), nil
}

func (s *sapSystemsService) GetAllDatabasesSIDs() ([]string, error) {
	var sids pq.StringArray

	err := s.db.
		Model(&entities.SAPSystemInstance{}).
		Statement.Where("type = ?", models.SAPSystemTypeDatabase).
		Distinct().
		Pluck("sid", &sids).
		Error

	if err != nil {
		return nil, err
	}

	return []string(sids), nil
}

func (s *sapSystemsService) GetAllApplicationsTags() ([]string, error) {
	var tags []string

	err := s.db.
		Model(&models.Tag{}).
		Where("resource_type = ?", models.TagSAPSystemResourceType).
		Distinct().
		Pluck("value", &tags).
		Error

	if err != nil {
		return nil, err

	}

	return tags, nil
}

func (s *sapSystemsService) GetAllDatabasesTags() ([]string, error) {
	var tags []string

	err := s.db.
		Model(&models.Tag{}).
		Where("resource_type = ?", models.TagDatabaseResourceType).
		Distinct().
		Pluck("value", &tags).
		Error

	if err != nil {
		return nil, err

	}

	return tags, nil
}

func (s *sapSystemsService) getAllByType(sapSystemType string, tagResourceType string, filter *SAPSystemFilter, page *Page) (models.SAPSystemList, error) {
	var instances entities.SAPSystemInstances

	paginationSubQuery := s.db.
		Distinct("id,sid").
		Where("type = ?", sapSystemType).
		Scopes(Paginate(page)).
		Order("sid").
		Table("sap_system_instances")

	db := s.db.
		Preload("Host").
		Preload("Tags", "resource_type = (?)", tagResourceType).
		Where("(id,sid) IN (?)", paginationSubQuery).
		Order("sid, instance_number, system_replication, id")

	if filter != nil {
		if len(filter.SIDs) > 0 {
			db = db.Where("sid IN (?)", filter.SIDs)
		}

		if len(filter.Tags) > 0 {
			db = db.Where("id IN (?)", s.db.Model(&models.Tag{}).
				Select("resource_id").
				Where("resource_type = ?", tagResourceType).
				Where("value IN ?", filter.Tags),
			)
		}
	}

	err := db.Find(&instances).Error
	if err != nil {
		return nil, err
	}

	sapSystemList := instances.ToModel()
	err = s.enrichSAPSystemList(sapSystemList)
	if err != nil {
		return nil, err
	}

	return sapSystemList, nil
}

func (s *sapSystemsService) getAttachedDatabase(dbName string, dbAddress string) (*models.SAPSystem, error) {
	var primaryInstance entities.SAPSystemInstance

	db := s.db.
		Model(&entities.SAPSystemInstance{}).
		Joins("JOIN hosts ON sap_system_instances.agent_id = hosts.agent_id")

	ip := net.ParseIP(dbAddress)
	if ip.To4() == nil {
		return nil, fmt.Errorf("received database address is not valid: %s", dbAddress)
	}

	db = db.Where("hosts.ip_addresses && ?", pq.Array([]string{dbAddress}))

	err := db.Where("tenants && ?", pq.Array([]string{dbName})).
		Select("id").
		First(&primaryInstance).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	var instances entities.SAPSystemInstances
	err = s.db.
		Where("id", primaryInstance.ID).
		Preload("Host").
		Order("sid, instance_number, system_replication, id").
		Find(&instances).
		Error

	if err != nil {
		return nil, err
	}

	sapSystem := instances.ToModel()[0]
	s.computeHealth(sapSystem)

	return sapSystem, nil
}

func (s *sapSystemsService) enrichSAPSystemList(sapSystemList models.SAPSystemList) error {
	sids := make(map[string]int)
	for _, sapSystem := range sapSystemList {
		err := s.attachDatabase(sapSystem)
		if err != nil {
			log.Warnf("could not attach database: %s", err)
		}
		s.computeHealth(sapSystem)
		// Store already found SIDs to find duplicates
		sids[sapSystem.SID] += 1
	}

	for _, sapSystem := range sapSystemList {
		if sids[sapSystem.SID] > 1 {
			sapSystem.HasDuplicatedSID = true
		}
	}

	return nil
}

func (s *sapSystemsService) attachDatabase(sapSystem *models.SAPSystem) error {
	if sapSystem.Type == models.SAPSystemTypeApplication {
		attachedDatabase, err := s.getAttachedDatabase(sapSystem.DBName, sapSystem.DBAddress)
		if err != nil {
			return err
		}
		sapSystem.AttachedDatabase = attachedDatabase
	}
	return nil
}

func (s *sapSystemsService) computeHealth(sapSystem *models.SAPSystem) {
	sapSystem.Health = models.SAPSystemHealthPassing
	for _, sapInstance := range sapSystem.GetAllInstances() {
		switch {
		case sapInstance.Health() == models.SAPSystemHealthCritical:
			sapSystem.Health = models.SAPSystemHealthCritical
		case sapSystem.Health != models.SAPSystemHealthCritical && sapInstance.Health() == models.SAPSystemHealthWarning:
			sapSystem.Health = models.SAPSystemHealthWarning
		case sapSystem.Health == models.SAPSystemHealthPassing && sapInstance.Health() == models.SAPSystemHealthUnknown:
			sapSystem.Health = models.SAPSystemHealthUnknown
		}
	}
}
