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
		var hostIDs []string
		var clusterIDs []string

		for _, instance := range sapSystem.GetAllInstances() {
			hostIDs = append(hostIDs, instance.HostID)
			clusterIDs = append(clusterIDs, instance.ClusterID)
		}

		hosts, err := s.hostsService.GetAll(&HostsFilter{ID: hostIDs}, nil)
		if err != nil {
			return nil, err
		}

		clusters, err := s.clustersService.GetAll(&ClustersFilter{
			ID:          clusterIDs,
			ClusterType: []string{models.ClusterTypeHANAScaleUp},
		}, nil)
		if err != nil {
			return nil, err
		}

		healthSummary = append(healthSummary, models.SAPSystemHealthSummary{
			ID:              sapSystem.ID,
			SID:             sapSystem.SID,
			SAPSystemHealth: computeSAPSystemHealth(sapSystem),
			DatabaseHealth:  computeSAPSystemHealth(sapSystem.AttachedDatabase),
			ClustersHealth:  computeAggregatedClustersHealth(clusters),
			HostsHealth:     computeAggregatedHostsHealth(hosts),
		})
	}

	return healthSummary, nil
}

func computeSAPSystemHealth(sapsystem *models.SAPSystem) string {
	if sapsystem == nil {
		return models.HealthSummaryHealthUnknown
	}

	switch sapsystem.Health {
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
	if len(clusters) == 0 {
		return models.HealthSummaryHealthUnknown
	}

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
