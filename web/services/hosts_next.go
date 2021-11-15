package services

import (
	"time"

	"github.com/lib/pq"
	"github.com/trento-project/trento/internal"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const HeartbeatTreshold = 5 * 2 * time.Second

var timeSince = time.Since

//go:generate mockery --name=HostsNextService --inpackage --filename=hosts_next_mock.go
type HostsNextService interface {
	GetAll(*HostsFilter, *Page) (models.HostList, error)
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

type hostsNextService struct {
	db *gorm.DB
}

func NewHostsNextService(db *gorm.DB) *hostsNextService {
	return &hostsNextService{db}
}

func (s *hostsNextService) GetAll(filter *HostsFilter, page *Page) (models.HostList, error) {
	var hosts []entities.Host
	db := s.db.Preload("Tags").Preload("Heartbeat")

	if filter != nil {
		if len(filter.SIDs) > 0 {
			db = db.Where("sids && ?", pq.Array(filter.SIDs))
		}

		if len(filter.Tags) > 0 {
			db = db.Where("agent_id IN (?)", s.db.Model(&models.Tag{}).
				Select("resource_id").
				Where("resource_type = ?", models.TagHostResourceType).
				Where("value IN ?", filter.Tags),
			)
		}
	}

	err := db.Find(&hosts).Order("name").Error
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

func (s *hostsNextService) GetCount() (int, error) {
	var count int64
	err := s.db.Model(&entities.Host{}).Count(&count).Error

	return int(count), err
}

func (s *hostsNextService) GetAllSIDs() ([]string, error) {
	var sids pq.StringArray

	err := s.db.Model(&entities.Host{}).
		Where("sids IS NOT NULL").
		Distinct().
		Pluck("unnest(sids)", &sids).
		Error

	if err != nil {
		return nil, err
	}

	return []string(sids), nil
}

func (s *hostsNextService) GetAllTags() ([]string, error) {
	var tags []string

	err := s.db.
		Model(&models.Tag{
			ResourceType: models.TagHostResourceType,
		}).
		Distinct().
		Pluck("value", &tags).
		Error

	if err != nil {
		return nil, err

	}

	return tags, nil
}

func (s *hostsNextService) Heartbeat(agentID string) error {
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
