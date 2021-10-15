package services

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/trento-project/trento/internal"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
)

//go:generate mockery --name=ClustersService --inpackage --filename=clusters_mock.go

type ClustersService interface {
	GetAll(filters map[string][]string) (models.ClusterList, error)
}

type clustersService struct {
	db            *gorm.DB
	checksService ChecksService
	tagsService   TagsService
}

func NewClustersService(db *gorm.DB, checksService ChecksService, tagsService TagsService) *clustersService {
	return &clustersService{
		db:            db,
		checksService: checksService,
		tagsService:   tagsService,
	}
}

func (s *clustersService) GetAll(filters map[string][]string) (models.ClusterList, error) {
	var clusterList models.ClusterList
	db := s.db

	if sids, ok := filters["sid"]; ok {
		if len(sids) > 0 {
			db = s.db.Where("sids && ?", pq.Array(sids))
		}
	}

	for _, f := range []string{"name", "cluster_type"} {
		if v, ok := filters[f]; ok {
			if len(v) > 0 {
				q := fmt.Sprintf("%s IN (?)", f)
				db = s.db.Where(q, v)
			}
		}
	}

	err := db.Find(&clusterList).Error
	if err != nil {
		return nil, err
	}

	err = s.enrichClusterData(clusterList)
	if err != nil {
		return nil, err
	}

	if tagsFilter, ok := filters["tags"]; ok {
		clusterList = filterByTags(clusterList, tagsFilter)
	}

	if healthFilter, ok := filters["health"]; ok {
		clusterList = filterByHealth(clusterList, healthFilter)
	}

	//db.Model(&User{}).Where("name = ?", "jinzhu").Count(&count) > 0
	return clusterList, nil
}

func (s *clustersService) enrichClusterData(clusterList models.ClusterList) error {
	names := make(map[string]int)
	for _, c := range clusterList {
		names[c.Name] += 1
	}

	for _, c := range clusterList {
		if names[c.Name] > 1 {
			c.HasDuplicatedName = true
		}
		health, _ := s.checksService.GetAggregatedChecksResultByCluster(c.ID)
		c.Health = health.String()

		tags, err := s.tagsService.GetAllByResource(models.TagClusterResourceType, c.ID)
		if err != nil {
			return err
		}
		c.Tags = tags
	}

	return nil
}

func filterByTags(clusterList models.ClusterList, tags []string) models.ClusterList {
	var filteredClusterList models.ClusterList

	for _, c := range clusterList {
		for _, t := range tags {
			if internal.Contains(c.Tags, t) {
				filteredClusterList = append(filteredClusterList, c)
				break
			}
		}
	}

	return filteredClusterList
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
