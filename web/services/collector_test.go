package services

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/datapipeline"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
)

type CollectorServiceTestSuite struct {
	suite.Suite
	db               *gorm.DB
	tx               *gorm.DB
	ch               chan *datapipeline.DataCollectedEvent
	collectorService *collectorService
}

func TestCollectorServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CollectorServiceTestSuite))
}

func (suite *CollectorServiceTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())

	suite.db.AutoMigrate(&datapipeline.DataCollectedEvent{})
}

func (suite *CollectorServiceTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(models.Tag{})
}

func (suite *CollectorServiceTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()

	ch := make(chan *datapipeline.DataCollectedEvent, 1)
	suite.ch = ch
	suite.collectorService = NewCollectorService(suite.tx, ch)
}

func (suite *CollectorServiceTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func (suite *CollectorServiceTestSuite) TestCollectorService_StoreEvent() {
	suite.collectorService.StoreEvent(&datapipeline.DataCollectedEvent{
		AgentID:       "agent_id",
		DiscoveryType: "test_discovery_type",
		Payload:       []byte("{}"),
	})

	eventFromChannel := <-suite.ch
	var eventFromDB datapipeline.DataCollectedEvent
	suite.tx.First(&eventFromDB)

	suite.EqualValues(eventFromChannel.ID, eventFromDB.ID)
	suite.EqualValues(eventFromChannel.AgentID, eventFromDB.AgentID)
	suite.EqualValues(eventFromChannel.DiscoveryType, eventFromDB.DiscoveryType)
	suite.EqualValues(eventFromChannel.Payload, eventFromDB.Payload)
}

func (suite *CollectorServiceTestSuite) TestCollectorService_StoreEvent_Skipping() {
	firstEvent := &datapipeline.DataCollectedEvent{
		AgentID:       "agent_id",
		DiscoveryType: "test_discovery_type",
		Payload:       []byte("{}"),
	}
	err := suite.collectorService.StoreEvent(firstEvent)
	suite.NoError(err)
	<-suite.ch

	var eventFromDBFirst datapipeline.DataCollectedEvent
	suite.tx.Last(&eventFromDBFirst)

	secondEvent := &datapipeline.DataCollectedEvent{
		AgentID:       "agent_id",
		DiscoveryType: "test_discovery_type",
		Payload:       []byte("{}"),
	}
	err = suite.collectorService.StoreEvent(secondEvent)
	suite.NoError(err)
	<-suite.ch

	var eventFromDBSecond datapipeline.DataCollectedEvent
	err = suite.tx.Last(&eventFromDBSecond).Error
	suite.NoError(err)

	var allEvents []*datapipeline.DataCollectedEvent
	err = suite.tx.Find(&allEvents).Error
	suite.NoError(err)

	suite.EqualValues(eventFromDBFirst.ID, eventFromDBSecond.ID)
	suite.EqualValues(eventFromDBFirst.CreatedAt, eventFromDBSecond.CreatedAt)
	suite.Equal(1, len(allEvents))
}
