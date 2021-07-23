package projectors

import (
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/trento-project/trento/web/models"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"

	"github.com/cloudquery/sqlite"
	"gorm.io/gorm"

	"github.com/trento-project/trento/web/projectors/mocks"
)

func setupProjectorTests() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(Subscription{}, models.Host{})
	if err != nil {
		panic(err)
	}

	return db
}

func teardownProjectorTests(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	if err := sqlDB.Close(); err != nil {
		panic(err)
	}
}

// TestProjector_Run tests that the lastSeenIndex is updated if the handler Project functions is successful
func TestProjector_Run(t *testing.T) {
	var persistedSubscription *Subscription
	var wg sync.WaitGroup

	db := setupProjectorTests()
	defer teardownProjectorTests(db)
	wg.Add(1)

	handler := new(mocks.ProjectorHandler)
	handler.On("GetName").Return("dummy")
	handler.On("Query", mock.Anything).Return(struct{}{}, uint64(1337), nil)
	handler.On("Project", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		wg.Done()
	}).Return(nil)

	projector := NewProjector(handler, db)
	projector.Run()

	if waitTimeout(&wg, 1*time.Second) {
		t.Fatal("timeout")
	}

	assert.Equal(t, uint64(1337), projector.subscription.LastSeenIndex)
	db.Where(&Subscription{Projector: "dummy"}).First(&persistedSubscription)
	assert.Equal(t, projector.subscription.LastSeenIndex, persistedSubscription.LastSeenIndex)
}

// TestProjector_Run_Error tests that the lastSeenIndex is not updated if the handler Project functions returns an error
func TestProjector_Run_Error(t *testing.T) {
	var persistedSubscription *Subscription
	var wg sync.WaitGroup

	db := setupProjectorTests()
	defer teardownProjectorTests(db)
	wg.Add(1)

	handler := new(mocks.ProjectorHandler)

	handler.On("GetName").Return("dummy")
	handler.On("Query", mock.Anything).Return(struct{}{}, uint64(1337), nil)
	handler.On("Project", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		wg.Done()
	}).Return(errors.New("error"))

	projector := NewProjector(handler, db)
	projector.Run()

	if waitTimeout(&wg, 1*time.Second) {
		t.Fatal("timeout")
	}
	assert.Equal(t, uint64(0), projector.subscription.LastSeenIndex)
	db.Where(&Subscription{Projector: "dummy"}).First(&persistedSubscription)
	assert.Equal(t, projector.subscription.LastSeenIndex, persistedSubscription.LastSeenIndex)
}

// TestHostsHandler_Project tests that consul nodes are projected correctly
func TestHostsHandler_Project(t *testing.T) {
	var hosts []models.Host

	db := setupProjectorTests()
	defer teardownProjectorTests(db)

	consulNodes := []*consulApi.Node{
		{
			Node:    "node1",
			Address: "192.168.1.1",
			Meta: map[string]string{
				"trento-cloud-provider":  "azure",
				"trento-ha-cluster":      "cluster",
				"trento-sap-system":      "sys1",
				"trento-sap-landscape":   "land1",
				"trento-sap-environment": "env1",
			},
		},
		{
			Node:    "node2",
			Address: "192.168.1.2",
			Meta: map[string]string{
				"trento-ha-cluster":      "cluster",
				"trento-sap-system":      "sys2",
				"trento-sap-landscape":   "land2",
				"trento-sap-environment": "env2",
				"trento-cloud-provider":  "gcp",
			},
		},
	}

	handler := NewHostsHandler("hosts", time.Second, nil)
	if err := handler.Project(db, consulNodes); err != nil {
		t.Fatal(err)
	}

	expected := []models.Host{
		{
			Name:          "node1",
			Address:       "192.168.1.1",
			Cluster:       "cluster",
			Environment:   "env1",
			Landscape:     "land1",
			SAPSystem:     "sys1",
			CloudProvider: "azure",
		},
		{
			Name:          "node2",
			Address:       "192.168.1.2",
			Cluster:       "cluster",
			Environment:   "env2",
			Landscape:     "land2",
			SAPSystem:     "sys2",
			CloudProvider: "gcp",
		},
	}

	db.Find(&hosts, []string{"node1", "node2"})
	assert.Equal(t, 2, len(hosts))
	assert.Equal(t, expected, hosts)
}

// TestHostsHandler_Project tests that consul nodes health statuses are projected correctly
func TestHostsHealthHandler_Project(t *testing.T) {
	var hosts []models.Host

	db := setupProjectorTests()
	defer teardownProjectorTests(db)

	healthChecks := consulApi.HealthChecks{
		&consulApi.HealthCheck{
			Node:   "node1",
			Status: consulApi.HealthPassing,
		},
		&consulApi.HealthCheck{
			Node:   "node1",
			Status: consulApi.HealthWarning,
		},
		&consulApi.HealthCheck{
			Node:   "node2",
			Status: consulApi.HealthPassing,
		},
		&consulApi.HealthCheck{
			Node:   "node2",
			Status: consulApi.HealthCritical,
		},
	}

	handler := NewHostsHealthHandler("hosts_health", time.Second, nil)
	if err := handler.Project(db, healthChecks); err != nil {
		t.Fatal(err)
	}

	db.Find(&hosts, []string{"node1", "node2"})
	assert.Equal(t, 2, len(hosts))
	assert.Equal(t, consulApi.HealthWarning, hosts[0].Health)
	assert.Equal(t, consulApi.HealthCritical, hosts[1].Health)
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
