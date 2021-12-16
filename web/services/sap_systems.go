package services

import (
	"errors"
	"net"

	"github.com/lib/pq"
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
	applications, err := s.getAllByType(models.SAPSystemTypeApplication, models.TagSAPSystemResourceType, filter, page)

	for _, a := range applications {
		a.AttachedDatabase, err = s.getAttachedDatabase(a.DBName, a.DBHost)
		if err != nil {
			return nil, err
		}
	}

	return applications, err
}

func (s *sapSystemsService) GetAllDatabases(filter *SAPSystemFilter, page *Page) (models.SAPSystemList, error) {
	return s.getAllByType(models.SAPSystemTypeDatabase, models.TagDatabaseResourceType, filter, page)
}
func (s *sapSystemsService) GetByID(ID string) (*models.SAPSystem, error) {
	var instances entities.SAPSystemInstances

	err := s.db.
		Where("id = ?", ID).
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

	db := s.db.
		Preload("Host").
		Scopes(Paginate(page)).
		Preload("Tags", "resource_type = (?)", tagResourceType).
		Where("type = ?", sapSystemType)

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
	s.ernichSAPSystemList(sapSystemList)

	return sapSystemList, nil
}

func (s *sapSystemsService) getAttachedDatabase(dbName string, dbHost string) (*models.SAPSystem, error) {
	var primaryInstance entities.SAPSystemInstance

	db := s.db.
		Model(&entities.SAPSystemInstance{}).
		Joins("JOIN hosts ON sap_system_instances.agent_id = hosts.agent_id")

	ip := net.ParseIP(dbHost)
	if ip.To4() == nil {
		db = db.Where("hosts.name = ?", dbHost)
	} else {
		db = db.Where("hosts.ip_addresses && ?", pq.Array([]string{dbHost}))
	}

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
		Find(&instances).
		Error

	if err != nil {
		return nil, err
	}

	return instances.ToModel()[0], nil
}

func (s *sapSystemsService) ernichSAPSystemList(sapSystemList models.SAPSystemList) {
	sids := make(map[string]int)
	for _, sapSystem := range sapSystemList {
		sids[sapSystem.SID] += 1
	}
	for _, sapSystem := range sapSystemList {
		if sids[sapSystem.SID] > 1 {
			sapSystem.HasDuplicatedSID = true
		}
	}
}
