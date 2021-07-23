package service

import (
	"fmt"

	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
)

//go:generate mockery --all

type IHostsService interface {
	GetHosts(page *Page, filters map[string][]string) []models.Host
	GetHostsSAPSystems() []string
	GetHostsLandscapes() []string
	GetHostsEnvironments() []string
	GetHostsCount() int
}

type HostsService struct {
	db *gorm.DB
}

func NewHostsService(db *gorm.DB) *HostsService {
	return &HostsService{db: db}
}

func (h *HostsService) GetHosts(page *Page, filters map[string][]string) []models.Host {
	var hosts []models.Host

	db := h.db.Scopes(Paginate(page))

	for _, f := range []string{"health", "sap_system", "landscape", "environment"} {
		if v, ok := filters[f]; ok {
			if len(v) > 0 {
				q := fmt.Sprintf("%s IN ?", f)
				db.Where(q, v)
			}
		}
	}

	db.Find(&hosts)
	return hosts
}

func (h *HostsService) GetHostsSAPSystems() []string {
	var sapsystems []string

	h.db.Model(&models.Host{}).Not("sap_system = ?", "").Distinct().Pluck("sap_system", &sapsystems)

	return sapsystems
}

func (h *HostsService) GetHostsLandscapes() []string {
	var landscapes []string

	h.db.Model(&models.Host{}).Not("landscape = ?", "").Distinct().Pluck("landscape", &landscapes)

	return landscapes
}

func (h *HostsService) GetHostsEnvironments() []string {
	var environments []string

	h.db.Model(&models.Host{}).Not("environment = ?", "").Distinct().Pluck("environment", &environments)

	return environments
}

func (h *HostsService) GetHostsCount() int {
	var hosts []*models.Host
	var count int64

	h.db.Find(&hosts).Count(&count)

	return int(count)
}
