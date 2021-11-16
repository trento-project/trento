package datapipeline

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/agent/discovery/mocks"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/entities"
	"gorm.io/gorm"
)

type SAPSystemsProjectorTestSuite struct {
	suite.Suite
	db *gorm.DB
	tx *gorm.DB
}

func TestSAPSystemsProjectorTestSuite(t *testing.T) {
	suite.Run(t, new(SAPSystemsProjectorTestSuite))
}

func (suite *SAPSystemsProjectorTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())

	suite.db.AutoMigrate(&Subscription{}, &entities.SAPSystemInstance{})
}

func (suite *SAPSystemsProjectorTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(Subscription{}, entities.SAPSystemInstance{})
}

func (suite *SAPSystemsProjectorTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
}

func (suite *SAPSystemsProjectorTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func (s *SAPSystemsProjectorTestSuite) Test_SAPSystemDiscoveryHandler() {
	discoveredSAPSystemMock := mocks.NewDiscoveredSAPSystemMock()

	requestBody, _ := json.Marshal(discoveredSAPSystemMock)
	SAPSystemsProjector_SAPSystemsDiscoveryHandler(&DataCollectedEvent{
		ID:            1,
		AgentID:       "agent_id",
		DiscoveryType: SAPsystemDiscovery,
		Payload:       requestBody,
	}, s.tx)

	var projectedSAPSystemInstance entities.SAPSystemInstance
	s.tx.First(&projectedSAPSystemInstance)

	s.EqualValues("PRD", projectedSAPSystemInstance.SID)

	s.Equal("database", projectedSAPSystemInstance.Type)
	s.Equal("e06e328f8d6b0f46c1e66ffcd44d0dd7", projectedSAPSystemInstance.SystemID)
	s.Equal("00", projectedSAPSystemInstance.InstanceNumber)
	s.Equal("HDB|HDB_WORKER", projectedSAPSystemInstance.Features)
}
