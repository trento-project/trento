package datapipeline

import (
	"encoding/json"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/agent/discovery/mocks"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
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

func (s *SAPSystemsProjectorTestSuite) Test_SAPSystemDiscoveryHandler_Database() {
	discoveredSAPSystemMock := mocks.NewDiscoveredSAPSystemDatabaseMock()

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

	s.Equal(models.SAPSystemTypeDatabase, projectedSAPSystemInstance.Type)
	s.Equal("e06e328f8d6b0f46c1e66ffcd44d0dd7", projectedSAPSystemInstance.ID)
	s.Equal("00", projectedSAPSystemInstance.InstanceNumber)
	s.Equal("HDB|HDB_WORKER", projectedSAPSystemInstance.Features)
	s.Equal("Primary", projectedSAPSystemInstance.SystemReplication)
	s.Equal("SFAIL", projectedSAPSystemInstance.SystemReplicationStatus)
	s.Equal(pq.StringArray{"PRD"}, projectedSAPSystemInstance.Tenants)
	s.Equal("vmhana01", projectedSAPSystemInstance.SAPHostname)
	s.Equal("0.3", projectedSAPSystemInstance.StartPriority)
	s.Equal(50013, projectedSAPSystemInstance.HttpPort)
	s.Equal(50014, projectedSAPSystemInstance.HttpsPort)
}

func (s *SAPSystemsProjectorTestSuite) Test_SAPSystemDiscoveryHandler_Database_Obsolete() {
	err := s.tx.Create(&entities.SAPSystemInstance{
		SID:                     "PRD",
		Type:                    models.SAPSystemTypeDatabase,
		ID:                      "b6fa9c04ee8280357a35baf9ee73539d",
		AgentID:                 "agent_id",
		InstanceNumber:          "00",
		Features:                "HDB|HDB_WORKER",
		SystemReplication:       "Primary",
		SystemReplicationStatus: "SFAIL",
		Tenants:                 pq.StringArray{"PRD"},
		SAPHostname:             "vmhana01",
		StartPriority:           "0.3",
		HttpPort:                50013,
		HttpsPort:               50014,
	}).Error

	s.NoError(err)

	discoveredSAPSystemMock := mocks.NewDiscoveredSAPSystemDatabaseMock()

	requestBody, _ := json.Marshal(discoveredSAPSystemMock)
	SAPSystemsProjector_SAPSystemsDiscoveryHandler(&DataCollectedEvent{
		ID:            1,
		AgentID:       "agent_id",
		DiscoveryType: SAPsystemDiscovery,
		Payload:       requestBody,
	}, s.tx)

	var projectedSAPSystemInstance []entities.SAPSystemInstance
	s.tx.Find(&projectedSAPSystemInstance)

	s.Equal(1, len(projectedSAPSystemInstance))
	s.Equal("e06e328f8d6b0f46c1e66ffcd44d0dd7", projectedSAPSystemInstance[0].ID)
}

func (s *SAPSystemsProjectorTestSuite) Test_SAPSystemDiscoveryHandler_Application() {
	discoveredSAPSystemMock := mocks.NewDiscoveredSAPSystemApplicationMock()

	requestBody, _ := json.Marshal(discoveredSAPSystemMock)
	SAPSystemsProjector_SAPSystemsDiscoveryHandler(&DataCollectedEvent{
		ID:            1,
		AgentID:       "agent_id",
		DiscoveryType: SAPsystemDiscovery,
		Payload:       requestBody,
	}, s.tx)

	var projectedSAPSystemInstance entities.SAPSystemInstance
	s.tx.First(&projectedSAPSystemInstance)

	s.EqualValues("HA1", projectedSAPSystemInstance.SID)

	s.Equal(models.SAPSystemTypeApplication, projectedSAPSystemInstance.Type)
	s.Equal("7b65dc281f9fae2c8e68e6cab669993e", projectedSAPSystemInstance.ID)
	s.Equal("02", projectedSAPSystemInstance.InstanceNumber)
	s.Equal("ABAP|GATEWAY|ICMAN|IGS", projectedSAPSystemInstance.Features)
	s.Equal("PRD", projectedSAPSystemInstance.DBName)
	s.Equal("10.74.1.12", projectedSAPSystemInstance.DBHost)
	s.Equal("sapha1aas1", projectedSAPSystemInstance.SAPHostname)
	s.Equal("3", projectedSAPSystemInstance.StartPriority)
	s.Equal(50213, projectedSAPSystemInstance.HttpPort)
	s.Equal(50214, projectedSAPSystemInstance.HttpsPort)
}
