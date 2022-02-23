package services

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"

	prometheusModel "github.com/prometheus/common/model"
	prometheusInternal "github.com/trento-project/trento/internal/prometheus"
)

const (
	nodeExporterPort = 9100
	nodeExporterName = "Node Exporter"
)

//go:generate mockery --name=PrometheusService --inpackage --filename=prometheus_mock.go
type PrometheusService interface {
	GetHttpSDTargets() (models.PrometheusTargetsList, error)
	Query(query string, ts time.Time) (prometheusModel.Value, error)
}

type prometheusService struct {
	db            *gorm.DB
	prometheusApi prometheusInternal.PrometheusAPI
}

func NewPrometheusService(db *gorm.DB, promApi prometheusInternal.PrometheusAPI) *prometheusService {
	return &prometheusService{db, promApi}
}

func (p *prometheusService) GetHttpSDTargets() (models.PrometheusTargetsList, error) {
	var targetsList models.PrometheusTargetsList
	var hosts []entities.Host

	err := p.db.Find(&hosts).Error
	if err != nil {
		return targetsList, err
	}

	for _, host := range hosts {
		targets := &models.PrometheusTargets{
			Targets: []string{fmt.Sprintf("%s:%d", host.SSHAddress, nodeExporterPort)},
			Labels: map[string]string{
				"agentID":       host.AgentID,
				"hostname":      host.Name,
				"exporter_name": nodeExporterName,
			},
		}
		targetsList = append(targetsList, targets)
	}

	return targetsList, nil
}

func (p *prometheusService) Query(query string, ts time.Time) (prometheusModel.Value, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Debugf("Executing new query on prometheus: %s", query)
	result, warnings, err := p.prometheusApi.Query(ctx, query, ts)

	if err != nil {
		log.Errorf("Error querying prometheus: %v\n", err)
		return result, err
	}

	if len(warnings) > 0 {
		log.Warnf("Warnings querying prometheus: %v\n", warnings)
	}

	log.Debugf("Query executed successfully. Result: %v", result)
	return result, nil
}
