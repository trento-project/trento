package services

import (
	"encoding/json"
	"errors"
	"time"

	"fmt"
	"context"

	"github.com/lib/pq"
	"github.com/trento-project/trento/internal"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
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
	GetExportersState(hostname string) (map[string]bool, error)
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

var jobToExporter = map[string]string{
	"http_sd_clusters": "HA cluster exporter",
	"http_sd_hosts": "Host exporter",
}

func (s *hostsService) GetExportersState(hostname string) (map[string]bool, error) {
	exportersState := make(map[string]bool)

	client, err := api.NewClient(api.Config{
		Address: "http://localhost:9090",
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return exportersState, err
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := v1api.Query(ctx, fmt.Sprintf("up{hostname=\"%s\"}", hostname), time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		return exportersState, err
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	fmt.Printf("Result:\n%v\n", result)

	// Here how we could have more generic implementation
	// https://github.com/prometheus/client_golang/issues/194#issuecomment-254443951
  resultVal := result.(model.Vector)

	if len(resultVal) == 0 {
		return exportersState, nil
	}

	for _, result := range resultVal {
		exportersState[jobToExporter[string(result.Metric["job"])]] = (int(result.Value) == 1)
	}

	return exportersState, nil
}
