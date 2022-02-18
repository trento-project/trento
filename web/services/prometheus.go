package services

import (
	"fmt"

	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
)

const nodeExporterPort = 9100

//go:generate mockery --name=PrometheusService --inpackage --filename=prometheus_mock.go
type PrometheusService interface {
	GetHttpSDTargets() (models.PrometheusTargetsList, error)
}

type prometheusService struct {
	db *gorm.DB
}

func NewPrometheusService(db *gorm.DB) *prometheusService {
	return &prometheusService{db}
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
				"agentID":  host.AgentID,
				"hostname": host.Name,
			},
		}
		targetsList = append(targetsList, targets)
	}

	return targetsList, nil
}
