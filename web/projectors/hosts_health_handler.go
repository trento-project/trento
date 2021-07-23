package projectors

import (
	"fmt"
	"time"

	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/trento-project/trento/internal/consul"
)

type HostsHealthHandler struct {
	name     string
	waitTime time.Duration
	client   consul.Client
}

func NewHostsHealthHandler(name string, waitTime time.Duration, client consul.Client) *HostsHealthHandler {
	return &HostsHealthHandler{
		name,
		waitTime,
		client,
	}
}

func (h *HostsHealthHandler) GetName() string {
	return h.name
}

// Query long polls consul health checks endpoint and returns the health checks and last seen index
func (h *HostsHealthHandler) Query(lastIndex uint64) (interface{}, uint64, error) {
	query := &consulApi.QueryOptions{WaitIndex: lastIndex, WaitTime: h.waitTime}
	healthChecks, meta, err := h.client.Health().Checks("trento-agent", query)

	if err != nil {
		return nil, 0, err
	}

	return healthChecks, meta.LastIndex, err
}

// Project processes the consul health checks and stores them in the database
func (h *HostsHealthHandler) Project(db *gorm.DB, data interface{}) error {
	healthChecks, ok := data.(consulApi.HealthChecks)

	if !ok {
		return fmt.Errorf("projector handler %s: invalid type", h.name)
	}

	var err error
	var hosts []*models.Host

	healthChecksByNode := make(map[string]consulApi.HealthChecks)
	for _, h := range healthChecks {
		healthChecksByNode[h.Node] = append(healthChecksByNode[h.Node], h)
	}

	for node, h := range healthChecksByNode {
		hosts = append(hosts, &models.Host{
			Name:   node,
			Health: h.AggregatedStatus(),
		})
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if hosts != nil {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "name"}},
				DoUpdates: clause.AssignmentColumns([]string{"health"}),
			}).Create(hosts).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
