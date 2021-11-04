package datapipeline

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
)

type DataCollectedEventTestSuite struct {
	suite.Suite
	db *gorm.DB
	tx *gorm.DB
}

func TestDataCollectedEventTestSuite(t *testing.T) {
	suite.Run(t, new(DataCollectedEventTestSuite))
}

func (suite *DataCollectedEventTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())

	suite.db.AutoMigrate(&DataCollectedEvent{})
}

func (suite *DataCollectedEventTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(models.Tag{})
}

func (suite *DataCollectedEventTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
}

func (suite *DataCollectedEventTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func (suite *DataCollectedEventTestSuite) TestCollectorService_PruneEvents() {
	events := []DataCollectedEvent{
		{
			ID:            1,
			AgentID:       "agent_id",
			DiscoveryType: "test_discovery_type",
			Payload:       []byte("{}"),
			CreatedAt:     time.Now().Add(-24 * 15 * time.Hour),
		},
		{
			ID:            2,
			AgentID:       "agent_id",
			DiscoveryType: "test_discovery_type",
			Payload:       []byte("{}"),
			CreatedAt:     time.Now().Add(-24 * 10 * time.Hour),
		},
		{
			ID:            3,
			AgentID:       "agent_id",
			DiscoveryType: "test_discovery_type",
			Payload:       []byte("{}"),
			CreatedAt:     time.Now().Add(-24 * 6 * time.Hour),
		},
	}
	suite.tx.Create(events)

	PruneEvents(24*10*time.Hour, suite.tx)

	var prunedEvents []DataCollectedEvent
	suite.tx.Find(&prunedEvents)

	suite.Equal(1, len(prunedEvents))
	suite.Equal(int64(3), prunedEvents[0].ID)
}
