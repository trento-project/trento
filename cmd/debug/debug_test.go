package debug

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/datapipeline"
	"github.com/trento-project/trento/web/entities"
	"gorm.io/gorm"
)

type DebugTestSuite struct {
	suite.Suite
	db *gorm.DB
	tx *gorm.DB
}

func TestDebugTestSuite(t *testing.T) {
	suite.Run(t, new(DebugTestSuite))
}

func (suite *DebugTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())
}

func (suite *DebugTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
}

func (suite *DebugTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func (suite *DebugTestSuite) TestPruneEvents() {
	suite.tx.AutoMigrate(&datapipeline.DataCollectedEvent{})

	events := []datapipeline.DataCollectedEvent{
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

	pruneEvents(suite.tx, 24*10*time.Hour)

	var prunedEvents []datapipeline.DataCollectedEvent
	suite.tx.Find(&prunedEvents)

	suite.Equal(1, len(prunedEvents))
	suite.Equal(int64(3), prunedEvents[0].ID)
}

func (suite *DebugTestSuite) TestPruneChecksResults() {
	suite.tx.AutoMigrate(&entities.ChecksResult{})

	checksResults := []entities.ChecksResult{
		{
			ID:        1,
			GroupID:   "group_id",
			Payload:   []byte("{}"),
			CreatedAt: time.Now().Add(-24 * 15 * time.Hour),
		},
		{
			ID:        2,
			GroupID:   "group_id",
			Payload:   []byte("{}"),
			CreatedAt: time.Now().Add(-24 * 10 * time.Hour),
		},
		{
			ID:        3,
			GroupID:   "group_id",
			Payload:   []byte("{}"),
			CreatedAt: time.Now().Add(-24 * 6 * time.Hour),
		},
	}
	suite.tx.Create(checksResults)

	pruneChecksResults(suite.tx, 24*10*time.Hour)

	var prunedChecksResults []entities.ChecksResult
	suite.tx.Find(&prunedChecksResults)

	suite.Equal(1, len(prunedChecksResults))
	suite.Equal(int64(3), prunedChecksResults[0].ID)
}
