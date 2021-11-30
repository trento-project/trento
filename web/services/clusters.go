package services

import (
	"encoding/json"
	"errors"

	"github.com/lib/pq"
	"github.com/trento-project/trento/internal"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
)

//go:generate mockery --name=ClustersService --inpackage --filename=clusters_mock.go

type ClustersService interface {
	GetAll(*ClustersFilter, *Page) (models.ClusterList, error)
	GetByID(string) (*models.Cluster, error)
	GetCount() (int, error)
	GetAllClusterTypes() ([]string, error)
	GetAllSIDs() ([]string, error)
	GetAllTags() ([]string, error)
}

type ClustersFilter struct {
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
	db := s.db.Preload("Tags").Scopes(Paginate(page))

	if filter != nil {
		if len(filter.Name) > 0 {
			db = s.db.Where("name IN (?)", filter.Name)
		}

		if len(filter.ClusterType) > 0 {
			db = s.db.Where("cluster_type IN (?)", filter.ClusterType)
		}

		if len(filter.SIDs) > 0 {
			db = s.db.Where("sids && ?", pq.Array(filter.SIDs))
		}

		if len(filter.Tags) > 0 {
			db = db.Where("id IN (?)", s.db.Model(&models.Tag{}).
				Select("resource_id").
				Where("resource_type = ?", models.TagClusterResourceType).
				Where("value IN ?", filter.Tags),
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

	if filter != nil && len(filter.Health) > 0 {
		clusterList = filterByHealth(clusterList, filter.Health)
	}

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
		Where("sids IS NOT NULL").
		Distinct().
		Pluck("unnest(sids)", &sids).
		Error

	if err != nil {
		return nil, err
	}

	return []string(sids), nil
}

func (s *clustersService) GetAllTags() ([]string, error) {
	var tags []string

	err := s.db.
		Model(&models.Tag{ResourceType: models.TagClusterResourceType}).
		Distinct().
		Pluck("value", &tags).
		Error

	if err != nil {
		return nil, err

	}

	return tags, nil
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
	health, _ := s.checksService.GetAggregatedChecksResultByCluster(cluster.ID)

	cluster.Health = health.String()
	cluster.PassingCount = health.PassingCount
	cluster.WarningCount = health.WarningCount
	cluster.CriticalCount = health.CriticalCount

}

func (s *clustersService) enrichClusterNodes(nodes []*models.HANAClusterNode, clusterID string, hosts []*entities.Host) {
	for _, node := range nodes {
		for _, host := range hosts {
			if node.Name == host.Name {
				node.IPAddresses = append(node.IPAddresses, host.IPAddresses...)
				break
			}
		}
		c, _ := s.checksService.GetAggregatedChecksResultByHost(clusterID)
		if _, ok := c[node.Name]; ok {
			node.Health = c[node.Name].String()
		}
	}
}

func filterByHealth(clusterList models.ClusterList, health []string) models.ClusterList {
	var filteredClusterList models.ClusterList

	for _, c := range clusterList {
		if internal.Contains(health, c.Health) {
			filteredClusterList = append(filteredClusterList, c)
		}
	}

	return filteredClusterList
}
