package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/internal"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	prometheusModel "github.com/prometheus/common/model"
)

const HeartbeatTreshold = internal.HeartbeatInterval * 2

var timeSince = time.Since

//go:generate mockery --name=HostsService --inpackage --filename=hosts_mock.go
type HostsService interface {
	GetAll(*HostsFilter, *Page) (models.HostList, error)
	GetByID(string) (*models.Host, error)
	GetAllBySAPSystemID(string) (models.HostList, error)
	GetCount() (int, error)
	GetAllSIDs() ([]string, error)
	GetAllTags() ([]string, error)
	Heartbeat(agentID string) error
	GetExportersState(hostname string) (map[string]string, error)
}

type HostsFilter struct {
	ID     []string
	SIDs   []string
	Tags   []string
	Health []string
}

type hostsService struct {
	db                *gorm.DB
	prometheusService PrometheusService
}

func NewHostsService(db *gorm.DB, promService PrometheusService) *hostsService {
	return &hostsService{db, promService}
}

func (s *hostsService) GetAll(filter *HostsFilter, page *Page) (models.HostList, error) {
	var hosts []entities.Host
	var healthFilteredHosts []string

	// Filter the hosts by Health
	if filter != nil && len(filter.Health) > 0 {
		var heartbeats []entities.HostHeartbeat

		err := s.db.Find(&heartbeats).Error
		if err != nil {
			return nil, err
		}

		for _, hearbeat := range heartbeats {
			hearbeatHealth := computeHearbeatHealth(&hearbeat)
			if internal.Contains(filter.Health, hearbeatHealth) {
				healthFilteredHosts = append(healthFilteredHosts, hearbeat.AgentID)
			}
		}
	}

	db := s.db.
		Model(&entities.Host{}).
		Scopes(Paginate(page)).
		Preload("Tags").
		Preload("Heartbeat").
		Preload("SAPSystemInstances").
		Preload("SAPSystemInstances.Host")

	if filter != nil {
		if len(filter.ID) > 0 {
			db = db.Where("agent_id IN (?)", filter.ID)
		}

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

		if len(filter.Health) > 0 {
			db = db.Where("agent_id IN (?)", healthFilteredHosts)
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
		hostList = append(hostList, host)
	}

	return hostList, nil
}

func (s *hostsService) GetByID(id string) (*models.Host, error) {
	var host entities.Host
	err := s.db.
		Where("agent_id = ?", id).
		Preload("Heartbeat").
		Preload("SAPSystemInstances").
		First(&host).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	hostHealth := computeHealth(&host)
	modeledHost := host.ToModel()
	modeledHost.Health = hostHealth

	if modeledHost.CloudProvider == "azure" {
		var cloudData models.AzureCloudData
		json.Unmarshal(host.CloudData, &cloudData)
		modeledHost.CloudData = cloudData
	}

	return modeledHost, nil
}

func (s *hostsService) GetAllBySAPSystemID(id string) (models.HostList, error) {
	var hosts []entities.Host

	err := s.db.
		Order("name").
		Preload("Heartbeat").
		Preload("SAPSystemInstances").
		Joins("JOIN sap_system_instances ON sap_system_instances.agent_id = hosts.agent_id").
		Where("sap_system_instances.id = ?", id).
		Find(&hosts).
		Error

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

func (s *hostsService) GetCount() (int, error) {
	var count int64
	err := s.db.Model(&entities.Host{}).Count(&count).Error

	return int(count), err
}

func (s *hostsService) GetAllSIDs() ([]string, error) {
	var sids pq.StringArray

	err := s.db.
		Model(&entities.Host{}).
		Order("sap_system_instances.sid").
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
		Order("value").
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

func initJobsStates() map[string]string {
	states := make(map[string]string)
	states[nodeExporterName] = models.HostHealthUnknown
	return states
}

func (s *hostsService) GetExportersState(hostname string) (map[string]string, error) {
	jobsState := initJobsStates()
	result, err := s.prometheusService.Query(fmt.Sprintf("up{hostname=\"%s\"}", hostname), time.Now())
	if err != nil {
		log.Warnf("error querying to prometheus: %s", err)
		return jobsState, err
	}

	resultVector := result.(prometheusModel.Vector)

	if len(resultVector) == 0 {
		return jobsState, nil
	}

	for _, r := range resultVector {
		if _, ok := r.Metric["exporter_name"]; !ok {
			continue
		}
		name := string(r.Metric["exporter_name"])
		switch int(r.Value) {
		case 0:
			jobsState[name] = models.HostHealthCritical
		case 1:
			jobsState[name] = models.HostHealthPassing
		default:
			jobsState[name] = models.HostHealthUnknown
		}
	}

	return jobsState, nil
}

func computeHealth(host *entities.Host) string {
	return computeHearbeatHealth(host.Heartbeat)
}

func computeHearbeatHealth(hearbeat *entities.HostHeartbeat) string {
	if hearbeat == nil {
		return models.HostHealthUnknown
	}

	if timeSince(hearbeat.UpdatedAt) > HeartbeatTreshold {
		return models.HostHealthCritical
	}

	return models.HostHealthPassing
}
