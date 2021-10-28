package datapipeline

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	"gorm.io/gorm"
)

type ProjectorTestSuite struct {
	suite.Suite
	db *gorm.DB
	tx *gorm.DB
}

func TestProjectorTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectorTestSuite))
}

func (suite *ProjectorTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())

	suite.db.AutoMigrate(&Subscription{})
}

func (suite *ProjectorTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(Subscription{})
}

func (suite *ProjectorTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
}

func (suite *ProjectorTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

// TestProjector_Project tests that a projector updates the subscription correctly
func (suite *ProjectorTestSuite) TestProjector_Project() {
	projector := NewProjector("dummy_projector", suite.tx)
	handler := func(dataCollectedEvent *DataCollectedEvent, db *gorm.DB) error {
		return nil
	}

	projector.AddHandler("dummy_discovery_type", handler)

	projector.Project(&DataCollectedEvent{ID: 123, DiscoveryType: "dummy_discovery_type", AgentID: "345"})
	var subscription Subscription

	suite.tx.First(&subscription)
	suite.Equal(int64(123), subscription.LastProjectedEventID)
	suite.Equal("dummy_projector", subscription.ProjectorID)
	suite.Equal("345", subscription.AgentID)
	suite.NotEmpty(subscription.UpdatedAt)
}

// TestProkector_Project tests that a projector does not update the subscription in case of error
func (suite *ProjectorTestSuite) TestProjector_Project_Error() {
	projector := NewProjector("dummy_projector", suite.tx)
	handler := func(dataCollectedEvent *DataCollectedEvent, db *gorm.DB) error {
		return fmt.Errorf("kaboom")
	}

	suite.tx.Create(&Subscription{
		LastProjectedEventID: 123,
		ProjectorID:          "dummy_projector",
		AgentID:              "345",
	})

	projector.AddHandler("dummy_discovery_type", handler)
	projector.Project(&DataCollectedEvent{ID: 666, DiscoveryType: "dummy_discovery_type", AgentID: "345"})

	var subscription Subscription
	suite.tx.First(&subscription)

	suite.Equal(int64(123), subscription.LastProjectedEventID)
	suite.Equal("dummy_projector", subscription.ProjectorID)
	suite.Equal("345", subscription.AgentID)
	suite.NotEmpty(subscription.UpdatedAt)
}

// TestProjector_Concurrency tests that a projector routine waits for the previous one to finish
// if the same projector tries to project data regarding the same Agent
func (suite *ProjectorTestSuite) TestProjector_Concurrency() {
	projector := NewProjector("dummy_projector", suite.tx)
	ch := make(chan struct{})
	ch2 := make(chan struct{})

	handler := func(dataCollectedEvent *DataCollectedEvent, _ *gorm.DB) error {
		if dataCollectedEvent.ID == 1 {
			ch <- struct{}{}
			time.Sleep(500 * time.Millisecond)
		} else {
			ch2 <- struct{}{}
		}
		return nil
	}

	projector.AddHandler("dummy_discovery_type", handler)
	projector.AddHandler("dummy_discovery_type_2", handler)

	agentId := "agent_id"
	go projector.Project(&DataCollectedEvent{ID: 1, DiscoveryType: "dummy_discovery_type", AgentID: agentId})
	<-ch

	go projector.Project(&DataCollectedEvent{ID: 2, DiscoveryType: "dummy_discovery_type_2", AgentID: agentId})
	<-ch2

	var subscription Subscription
	suite.tx.First(&subscription)

	suite.Equal(int64(2), subscription.LastProjectedEventID)
}

// TestProjector_SkipPastEvent tests that a projector does not project and update the subscription
// if the event is older than the last projected one
func (suite *ProjectorTestSuite) TestProjector_SkipPastEvent() {
	projector := NewProjector("dummy_projector", suite.tx)
	projected := false

	projector.AddHandler("dummy_discovery_type", func(dataCollectedEvent *DataCollectedEvent, _ *gorm.DB) error {
		projected = true
		return nil
	})

	suite.tx.Create(&Subscription{
		LastProjectedEventID: 123,
		ProjectorID:          "dummy_projector",
		AgentID:              "345",
	})

	projector.Project(&DataCollectedEvent{ID: 120, DiscoveryType: "dummy_discovery_type", AgentID: "345"})

	var subscription Subscription
	suite.tx.First(&subscription)

	suite.False(projected)
	suite.Equal(int64(123), subscription.LastProjectedEventID)
}
