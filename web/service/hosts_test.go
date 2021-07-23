package service

import (
	"testing"

	consulApi "github.com/hashicorp/consul/api"

	"github.com/stretchr/testify/assert"

	"github.com/cloudquery/sqlite"
	"github.com/trento-project/trento/web/models"

	"gorm.io/gorm"
)

var hostsFixtures = []models.Host{
	{
		Name:          "foo",
		Address:       "192.168.1.1",
		Health:        consulApi.HealthPassing,
		Environment:   "env1",
		Landscape:     "land1",
		SAPSystem:     "sys1",
		CloudProvider: "azure",
	},
	{
		Name:          "bar",
		Address:       "192.168.1.2",
		Health:        consulApi.HealthCritical,
		Environment:   "env2",
		Landscape:     "land2",
		SAPSystem:     "sys2",
		CloudProvider: "aws",
	},
	{
		Name:          "buzz",
		Address:       "192.168.1.3",
		Health:        consulApi.HealthWarning,
		Environment:   "env3",
		Landscape:     "land3",
		SAPSystem:     "sys3",
		CloudProvider: "gcp",
	},
}

func setupHostsTest() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(models.Host{})
	if err != nil {
		panic(err)
	}

	db.Create(&hostsFixtures)
	return db
}

func teardownHostsTest(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	if err := sqlDB.Close(); err != nil {
		panic(err)
	}
}

func TestHostsService_GetHosts(t *testing.T) {
	db := setupHostsTest()
	hostsService := NewHostsService(db)
	defer teardownHostsTest(db)

	assert.Equal(t, hostsFixtures, hostsService.GetHosts(nil, map[string][]string{}))
}

func TestHostsService_GetHosts_Pagination(t *testing.T) {
	db := setupHostsTest()
	hostsService := NewHostsService(db)
	defer teardownHostsTest(db)

	expected := []models.Host{
		hostsFixtures[2],
	}

	assert.Equal(t, expected, hostsService.GetHosts(&Page{
		PageNr:   2,
		PageSize: 2,
	}, map[string][]string{}))
}

func TestHostsService_GetHosts_Filter(t *testing.T) {
	db := setupHostsTest()
	hostsService := NewHostsService(db)
	defer teardownHostsTest(db)

	filter := map[string][]string{
		"health": {"passing", "warning"},
	}
	expected := []models.Host{
		hostsFixtures[0],
		hostsFixtures[2],
	}
	assert.Equal(t, expected, hostsService.GetHosts(nil, filter))

	filter = map[string][]string{
		"environment": {"env1", "env3"},
	}
	expected = []models.Host{
		hostsFixtures[0],
		hostsFixtures[2],
	}
	assert.Equal(t, expected, hostsService.GetHosts(nil, filter))

	filter = map[string][]string{
		"landscape": {"land1", "land3"},
	}
	expected = []models.Host{
		hostsFixtures[0],
		hostsFixtures[2],
	}
	assert.Equal(t, expected, hostsService.GetHosts(nil, filter))

	filter = map[string][]string{
		"sap_system": {"sys1", "sys3"},
	}
	expected = []models.Host{
		hostsFixtures[0],
		hostsFixtures[2],
	}
	assert.Equal(t, expected, hostsService.GetHosts(nil, filter))
}

func TestHostsService_GetHostsCount(t *testing.T) {
	db := setupHostsTest()
	hostsService := NewHostsService(db)
	defer teardownHostsTest(db)

	assert.Equal(t, 3, hostsService.GetHostsCount())
}

func TestHostsService_GetHostsSAPSystems(t *testing.T) {
	db := setupHostsTest()
	hostsService := NewHostsService(db)
	defer teardownHostsTest(db)

	assert.Equal(t, []string{"sys1", "sys2", "sys3"}, hostsService.GetHostsSAPSystems())
}
