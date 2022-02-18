package services

import (
	"github.com/trento-project/trento/web/models"
)

//go:generate mockery --name=HealthSummaryService --inpackage --filename=health_summary_service_mock.go
type HealthSummaryService interface {
	GetHealthSummary() (models.HealthSummary, error)
}

type healthSummaryService struct {
	sapSystemsService SAPSystemsService
	hostsService      HostsService
	clustersService   ClustersService
}

func NewHealthSummaryService(sapSystemsService SAPSystemsService,
	clustersService ClustersService,
	hostsService HostsService) HealthSummaryService {
	return &healthSummaryService{
		sapSystemsService: sapSystemsService,
		clustersService:   clustersService,
		hostsService:      hostsService,
	}
}

func (s *healthSummaryService) GetHealthSummary() (models.HealthSummary, error) {
	var healthSummary models.HealthSummary

	sapSystems, err := s.sapSystemsService.GetAllApplications(nil, nil)
	if err != nil {
		return nil, err
	}

	for _, sapSystem := range sapSystems {
		var clusterIDs []string
		var hostIDs []string

		for _, instance := range sapSystem.GetAllInstances() {
			clusterIDs = append(clusterIDs, instance.ClusterID)
			hostIDs = append(hostIDs, instance.HostID)
		}

		clusters, err := s.clustersService.GetAll(&ClustersFilter{ID: clusterIDs}, nil)
		if err != nil {
			return nil, err
		}

		hosts, err := s.hostsService.GetAll(&HostsFilter{ID: hostIDs}, nil)
		if err != nil {
			return nil, err
		}

		healthSummary = append(healthSummary, models.SAPSystemHealthSummary{
			ID:             sapSystem.ID,
			SID:            sapSystem.SID,
			DatabaseHealth: computeDatabaseHealth(sapSystem.AttachedDatabase),
			ClustersHealth: computeAggregatedClustersHealth(clusters),
			HostsHealth:    computeAggregatedHostsHealth(hosts),
		})
	}

	return healthSummary, nil
}

func computeDatabaseHealth(database *models.SAPSystem) string {
	switch database.Health {
	case models.SAPSystemHealthPassing:
		return models.HealthSummaryHealthPassing
	case models.SAPSystemHealthWarning:
		return models.HealthSummaryHealthWarning
	case models.SAPSystemHealthCritical:
		return models.HealthSummaryHealthCritical
	default:
		return models.HealthSummaryHealthUnknown
	}
}

func computeAggregatedClustersHealth(clusters []*models.Cluster) string {
	var hasWarningCluster, hasUnknownCluster bool

	for _, c := range clusters {
		switch c.Health {
		case models.CheckCritical:
			return models.HealthSummaryHealthCritical
		case models.CheckWarning:
			hasWarningCluster = true
		case models.HealthSummaryHealthUnknown:
			hasUnknownCluster = true
		}
	}

	if hasWarningCluster {
		return models.HealthSummaryHealthWarning
	}

	if hasUnknownCluster {
		return models.HealthSummaryHealthUnknown
	}

	return models.HealthSummaryHealthPassing
}

func computeAggregatedHostsHealth(hosts []*models.Host) string {
	var hasWarningHost, hasUnknownHost bool

	for _, h := range hosts {
		switch h.Health {
		case models.HostHealthCritical:
			return models.HealthSummaryHealthCritical
		case models.HostHealthWarning:
			hasWarningHost = true
		case models.HostHealthUnknown:
			hasUnknownHost = true
		}
	}

	if hasWarningHost {
		return models.HealthSummaryHealthWarning
	}

	if hasUnknownHost {
		return models.HealthSummaryHealthUnknown
	}

	return models.HealthSummaryHealthPassing
}
