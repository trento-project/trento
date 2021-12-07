package services

import (
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/trento-project/trento/internal"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const HeartbeatTreshold = internal.HeartbeatInterval * 2

var timeSince = time.Since

//go:generate mockery --name=HostsService --inpackage --filename=hosts_mock.go
type HostsService interface {
	GetAll(*HostsFilter, *Page) (models.HostList, error)
	GetByID(string) (*models.Host, error)
	GetCount() (int, error)
	GetAllSIDs() ([]string, error)
	GetAllTags() ([]string, error)
	Heartbeat(agentID string) error
}

type HostsFilter struct {
	SIDs   []string
	Tags   []string
	Health []string
}

type hostsService struct {
	db *gorm.DB
}

func NewHostsService(db *gorm.DB) *hostsService {
	return &hostsService{db}
}

func (s *hostsService) GetAll(filter *HostsFilter, page *Page) (models.HostList, error) {
	var hosts []entities.Host
	db := s.db.
		Model(&entities.Host{}).
		Scopes(Paginate(page)).
		Preload("Tags").
		Preload("Heartbeat").
		Preload("SAPSystemInstances").
		Preload("SAPSystemInstances.Host")

	if filter != nil {
		if len(filter.SIDs) > 0 {
			db = db.Where("agent_id IN (?)", s.db.Model(&entities.SAPSystemInstance{}).
				Select("agent_id").
				Where("sid IN ?", filter.SIDs),
			)
		}

		if len(filter.Tags) > 0 {
			db = db.Where("agent_id IN (?)", s.db.Model(&models.Tag{}).
				Select("resource_id").
				Where("resource_type = ?", models.TagHostResourceType).
				Where("value IN ?", filter.Tags),
			)
		}
	}

	err := db.Order("name").Find(&hosts).Error
	if err != nil {
		return nil, err
	}

	var hostList models.HostList
	for _, h := range hosts {
		host := h.ToModel()
		host.Health = computeHealth(&h)

		if filter != nil && len(filter.Health) > 0 {
			if !internal.Contains(filter.Health, host.Health) {
				continue
			}
		}
		hostList = append(hostList, host)
	}

	return hostList, nil
}

func (s *hostsService) GetByID(id string) (*models.Host, error) {
	var host entities.Host
	err := s.db.
		Where("agent_id = ?", id).
		First(&host).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return host.ToModel(), nil
}

func (s *hostsService) GetCount() (int, error) {
	var count int64
	err := s.db.Model(&entities.Host{}).Count(&count).Error

	return int(count), err
}

func (s *hostsService) GetAllSIDs() ([]string, error) {
	var sids pq.StringArray

	err := s.db.
		Model(&entities.Host{}).
		Joins("JOIN sap_system_instances ON sap_system_instances.agent_id = hosts.agent_id AND sid IS NOT NULL").
		Distinct().
		Pluck("sap_system_instances.sid", &sids).
		Error

	if err != nil {
		return nil, err
	}

	return []string(sids), nil
}

func (s *hostsService) GetAllTags() ([]string, error) {
	var tags []string

	err := s.db.
		Model(&models.Tag{}).
		Where("resource_type = ?", models.TagHostResourceType).
		Distinct().
		Pluck("value", &tags).
		Error

	if err != nil {
		return nil, err

	}

	return tags, nil
}

func (s *hostsService) Heartbeat(agentID string) error {
	heartbeat := &entities.HostHeartbeat{
		AgentID: agentID,
	}

	return s.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "agent_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
	}).Create(heartbeat).Error
}

func computeHealth(host *entities.Host) string {
	if host.Heartbeat == nil {
		return models.HostHealthUnknown
	}

	if timeSince(host.Heartbeat.UpdatedAt) > HeartbeatTreshold {
		return models.HostHealthCritical
	}

	return models.HostHealthPassing
}
