package services

import (
	"time"

	"github.com/lib/pq"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const HeartbeatTreshold = 5 * 2 * time.Second

var timeSince = time.Since

//go:generate mockery --name=HostsNextService --inpackage --filename=hosts_next_mock.go
type HostsNextService interface {
	GetAll(filters map[string][]string) (models.HostList, error)
	GetAllSIDs() ([]string, error)
	GetAllTags() ([]string, error)
	Heartbeat(agentID string) error
}

type hostsNextService struct {
	db *gorm.DB
}

func NewHostsNextService(db *gorm.DB) *hostsNextService {
	return &hostsNextService{db}
}

func (s *hostsNextService) GetAll(filters map[string][]string) (models.HostList, error) {
	var hosts []entities.Host
	db := s.db.Preload("Tags").Preload("Heartbeat")

	if tags, ok := filters["tags"]; ok {
		db = db.Where("agent_id IN (?)", s.db.Model(&models.Tag{}).
			Select("resource_id").
			Where("resource_type = ?", models.TagHostResourceType).
			Where("value IN ?", tags),
		)
	}

	if sids, ok := filters["sids"]; ok {
		if len(sids) > 0 {
			db = db.Where("sids && ?", pq.Array(sids))
		}
	}

	err := db.Find(&hosts).Error
	if err != nil {
		return nil, err
	}

	var hostList models.HostList
	for _, h := range hosts {
		host := h.ToModel()
		host.Health = computeHealth(&h)

		hostList = append(hostList, host)
	}

	return hostList, nil
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
