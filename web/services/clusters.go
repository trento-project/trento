package services

import (
	"encoding/json"
	"errors"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/internal/cloud"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
)

//go:generate mockery --name=ClustersService --inpackage --filename=clusters_mock.go

type ClustersService interface {
	GetAll(*ClustersFilter, *Page) (models.ClusterList, error)
	GetByID(string) (*models.Cluster, error)
	GetCount() (int, error)
	GetAllClusterNames() ([]string, error)
	GetAllClusterTypes() ([]string, error)
	GetAllSIDs() ([]string, error)
	GetAllTags() ([]string, error)
	GetAllClustersSettings() (models.ClustersSettings, error)
	GetClusterSettingsByID(id string) (*models.ClusterSettings, error)
}

type ClustersFilter struct {
	ID          []string
	Name        []string
	ClusterType []string
	SIDs        []string
	Tags        []string
	Health      []string
}

type clustersService struct {
	db            *gorm.DB
	checksService ChecksService
}

func NewClustersService(db *gorm.DB, checksService ChecksService) *clustersService {
	return &clustersService{
		db:            db,
		checksService: checksService,
	}
}

func (s *clustersService) GetAll(filter *ClustersFilter, page *Page) (models.ClusterList, error) {
	var clusters []entities.Cluster

	db := s.db.Preload("Health").Preload("Tags").Scopes(Paginate(page))

	if filter != nil {
		if len(filter.ID) > 0 {
			db = db.Where("id IN (?)", filter.ID)
		}

		if len(filter.Name) > 0 {
			db = db.Where("name IN (?)", filter.Name)
		}

		if len(filter.ClusterType) > 0 {
			db = db.Where("cluster_type IN (?)", filter.ClusterType)
		}

		if len(filter.SIDs) > 0 {
			db = db.Where("sid IN (?)", filter.SIDs)
		}

		if len(filter.Tags) > 0 {
			db = db.Where("id IN (?)", s.db.Model(&models.Tag{}).
				Select("resource_id").
				Where("resource_type = ?", models.TagClusterResourceType).
				Where("value IN ?", filter.Tags),
			)
		}

		if len(filter.Health) > 0 {
			db = db.Where("id IN (?)", s.db.Model(&entities.HealthState{}).
				Select("id").
				Where("health IN ?", filter.Health),
			)
		}
	}

	err := db.Order("name").Order("id").Find(&clusters).Error
	if err != nil {
		return nil, err
	}

	var clusterList models.ClusterList
	for _, cluster := range clusters {
		clusterList = append(clusterList, cluster.ToModel())
	}

	s.enrichClusterList(clusterList)

	return clusterList, nil
}

func (s *clustersService) GetByID(clusterID string) (*models.Cluster, error) {
	var cluster entities.Cluster

	err := s.db.
		Preload("Hosts").
		Where("id = ?", clusterID).First(&cluster).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	clusterModel := cluster.ToModel()

	switch cluster.ClusterType {
	case models.ClusterTypeHANAScaleUp, models.ClusterTypeHANAScaleOut:
		var clusterDetailHANA entities.HANAClusterDetails

		err := json.Unmarshal(cluster.Details, &clusterDetailHANA)
		if err != nil {
			return nil, err
		}

		detail := clusterDetailHANA.ToModel()
		s.enrichClusterNodes(detail.Nodes, cluster.ID, cluster.Hosts)
		s.enrichCluster(clusterModel)
		clusterModel.Details = detail
	default:
		clusterModel.Details = nil
	}

	return clusterModel, nil
}

func (s *clustersService) GetCount() (int, error) {
	var count int64
	err := s.db.Model(&entities.Cluster{}).Count(&count).Error

	return int(count), err
}

func (s *clustersService) GetAllClusterNames() ([]string, error) {
	var clusterNames []string

	err := s.db.Model(&entities.Cluster{}).
		Distinct().
		Order("name").
		Pluck("name", &clusterNames).
		Error

	if err != nil {
		return nil, err
	}

	return clusterNames, nil
}

func (s *clustersService) GetAllClusterTypes() ([]string, error) {
	var clusterTypes []string

	err := s.db.Model(&entities.Cluster{}).
		Distinct().
		Pluck("cluster_type", &clusterTypes).
		Error

	if err != nil {
		return nil, err
	}

	return clusterTypes, nil
}

func (s *clustersService) GetAllSIDs() ([]string, error) {
	var sids pq.StringArray

	err := s.db.Model(&entities.Cluster{}).
		Distinct().
		Where("sid IS NOT NULL AND sid <> ''").
		Order("sid").
		Pluck("sid", &sids).
		Error

	if err != nil {
		return nil, err
	}

	return []string(sids), nil
}

func (s *clustersService) GetAllTags() ([]string, error) {
	var tags []string

	err := s.db.
		Model(&models.Tag{}).
		Where("resource_type = ?", models.TagClusterResourceType).
		Distinct().
		Pluck("value", &tags).
		Error

	if err != nil {
		return nil, err

	}

	return tags, nil
}

func (s *clustersService) GetAllClustersSettings() (models.ClustersSettings, error) {
	var clusters []*entities.Cluster

	err := s.db.
		Preload("Hosts").
		Find(&clusters).
		Error

	if err != nil {
		return nil, err
	}

	clustersSettings := models.ClustersSettings{}
	for _, cluster := range clusters {
		clusterSettings, err := s.loadSettings(cluster)
		if err != nil {
			log.Error(err)
			return clustersSettings, err
		}

		clustersSettings = append(clustersSettings, clusterSettings)
	}

	return clustersSettings, nil
}

func getDefaultUserName(host *entities.Host) (string, error) {
	switch host.CloudProvider {
	case cloud.Azure:
		var metadata entities.AzureCloudData
		err := json.Unmarshal(host.CloudData, &metadata)
		if err != nil {
			return "", err
		}
		return metadata.AdminUsername, nil
	default:
		return "root", nil
	}
}
func (s *clustersService) GetClusterSettingsByID(id string) (*models.ClusterSettings, error) {
	var cluster entities.Cluster

	err := s.db.
		Preload("Hosts").
		Where("id = ?", id).
		First(&cluster).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return s.loadSettings(&cluster)
}

func (s *clustersService) loadSettings(cluster *entities.Cluster) (*models.ClusterSettings, error) {
	var hosts []*models.HostConnection

	selectedChecks, err := s.checksService.GetSelectedChecksById(cluster.ID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	connectionSettings, err := s.checksService.GetConnectionSettingsById(cluster.ID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, host := range cluster.Hosts {
		var username string

		hostConnectionSettings, found := connectionSettings[host.Name]
		if found {
			username = hostConnectionSettings.User
		}

		if username == "" {
			username, err = getDefaultUserName(host)
			if err != nil {
				return nil, err
			}
		}

		hosts = append(hosts, &models.HostConnection{
			Name:    host.Name,
			Address: host.SSHAddress,
			User:    username,
		})
	}

	return &models.ClusterSettings{
		ID:             cluster.ID,
		SelectedChecks: selectedChecks.SelectedChecks,
		Hosts:          hosts,
	}, nil
}

func (s *clustersService) enrichClusterList(clusterList models.ClusterList) {
	names := make(map[string]int)
	for _, c := range clusterList {
		names[c.Name] += 1
	}

	for _, c := range clusterList {
		if names[c.Name] > 1 {
			c.HasDuplicatedName = true
		}
		s.enrichCluster(c)
	}
}

func (s *clustersService) enrichCluster(cluster *models.Cluster) {
	health, err := s.checksService.GetAggregatedChecksResultByCluster(cluster.ID)

	if err == nil {
		cluster.PassingCount = health.PassingCount
		cluster.WarningCount = health.WarningCount
		cluster.CriticalCount = health.CriticalCount
	}
}

func (s *clustersService) enrichClusterNodes(nodes []*models.HANAClusterNode, clusterID string, hosts []*entities.Host) {
	checkResults, checkResultsErr := s.checksService.GetAggregatedChecksResultByHost(clusterID)

	for _, node := range nodes {
		for _, host := range hosts {
			if node.Name == host.Name {
				node.HostID = host.AgentID
				node.IPAddresses = append(node.IPAddresses, host.IPAddresses...)
				node.Health = models.CheckUndefined
				break
			}
		}

		if checkResultsErr != nil {
			continue
		}

		if _, ok := checkResults[node.Name]; ok {
			node.Health = checkResults[node.Name].String()
		}
	}
}
