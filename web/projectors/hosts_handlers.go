package projectors

import (
	"fmt"
	"time"

	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm/clause"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/trento-project/trento/internal/consul"
	"gorm.io/gorm"
)

type HostsHandler struct {
	name     string
	waitTime time.Duration
	client   consul.Client
}

func NewHostsHandler(name string, waitTime time.Duration, client consul.Client) *HostsHandler {
	return &HostsHandler{
		name,
		waitTime,
		client,
	}
}

func (h *HostsHandler) GetName() string {
	return h.name
}

// Query long polls consul catalog and returns the nodes and last seen index
func (h *HostsHandler) Query(lastIndex uint64) (interface{}, uint64, error) {
	query := &consulApi.QueryOptions{WaitIndex: lastIndex, WaitTime: h.waitTime}
	consulNodes, meta, err := h.client.Catalog().Nodes(query)
	if err != nil {
		return nil, 0, err
	}

	return consulNodes, meta.LastIndex, err
}

// Project processes the consul nodes and stores them in the database
func (h *HostsHandler) Project(db *gorm.DB, data interface{}) error {
	consulNodes, ok := data.([]*consulApi.Node)

	if !ok {
		return fmt.Errorf("projector handler %s: invalid type", h.name)
	}

	var hosts []*models.Host
	for _, n := range consulNodes {
		host := &models.Host{
			Name:          n.Node,
			Address:       n.Address,
			Cluster:       n.Meta["trento-ha-cluster"],
			CloudProvider: n.Meta["trento-cloud-provider"],
			SAPSystem:     n.Meta["trento-sap-system"],
			Landscape:     n.Meta["trento-sap-landscape"],
			Environment:   n.Meta["trento-sap-environment"],
		}
		hosts = append(hosts, host)
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if hosts != nil {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "name"}},
				DoUpdates: clause.AssignmentColumns([]string{"name", "address", "cluster", "cloud_provider", "sap_system", "landscape", "environment"}),
			}).Create(hosts).Error; err != nil {
				return err
			}

			if err := tx.Not(hosts).Delete(&models.Host{}).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
