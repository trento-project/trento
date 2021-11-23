package datapipeline

import (
	"encoding/json"
	"testing"

	"github.com/trento-project/trento/agent/discovery/mocks"

	"github.com/stretchr/testify/suite"
	_ "github.com/trento-project/trento/test"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/entities"
	"gorm.io/gorm"
)

type HostTelemetryProjectorTestSuite struct {
	suite.Suite
	db *gorm.DB
	tx *gorm.DB
}

func TestHostTelemetryProjectorTestSuite(t *testing.T) {
	suite.Run(t, new(HostTelemetryProjectorTestSuite))
}

func (suite *HostTelemetryProjectorTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())

	suite.db.AutoMigrate(&Subscription{}, &entities.HostTelemetry{})
}

func (suite *HostTelemetryProjectorTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(Subscription{}, entities.HostTelemetry{})
}

func (suite *HostTelemetryProjectorTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
}

func (suite *HostTelemetryProjectorTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

// Test_HostDiscoveryHandler tests the HostDiscoveryHandler function execution on a HostDiscovery published by an agent
func (s *HostTelemetryProjectorTestSuite) Test_HostDiscoveryHandler() {
	discoveredHostMock := mocks.NewDiscoveredHostMock()

	requestBody, _ := json.Marshal(discoveredHostMock)

	hostTelemetryProjector_HostDiscoveryHandler(&DataCollectedEvent{
		ID:            1,
		AgentID:       "agent_id",
		DiscoveryType: HostDiscovery,
		Payload:       requestBody,
	}, s.tx)

	var projectedTelemetry entities.HostTelemetry
	s.tx.First(&projectedTelemetry)

	s.Equal(discoveredHostMock.HostName, projectedTelemetry.HostName)
	s.Equal(discoveredHostMock.CPUCount, projectedTelemetry.CPUCount)
	s.Equal(discoveredHostMock.SocketCount, projectedTelemetry.SocketCount)
	s.Equal(discoveredHostMock.TotalMemoryMB, projectedTelemetry.TotalMemoryMB)
	s.Equal(discoveredHostMock.OSVersion, projectedTelemetry.SLESVersion)
	s.Equal("", projectedTelemetry.CloudProvider)
}

// Test_CloudDiscoveryHandler tests the loudDiscoveryHandler function execution on a CloudDiscovery published by an agent
func (s *HostTelemetryProjectorTestSuite) Test_CloudDiscoveryHandler() {
	discoveredCloudMock := mocks.NewDiscoveredCloudMock()

	requestBody, _ := json.Marshal(discoveredCloudMock)

	hostTelemetryProjector_CloudDiscoveryHandler(&DataCollectedEvent{
		ID:            1,
		AgentID:       "agent_id",
		DiscoveryType: CloudDiscovery,
		Payload:       requestBody,
	}, s.tx)

	var projectedTelemetry entities.HostTelemetry
	s.tx.First(&projectedTelemetry)

	s.Equal("", projectedTelemetry.HostName)
	s.Equal(0, projectedTelemetry.CPUCount)
	s.Equal(0, projectedTelemetry.SocketCount)
	s.Equal(0, projectedTelemetry.TotalMemoryMB)
	s.Equal("", projectedTelemetry.SLESVersion)

	expectedCloudProvider := discoveredCloudMock.Provider
	s.Equal(expectedCloudProvider, projectedTelemetry.CloudProvider)
}

// Test_TelemetryProjector tests the TelemetryProjector projects all of the discoveries it is interested in, resulting in a single telemetry readmodel
func (s *HostTelemetryProjectorTestSuite) Test_TelemetryProjector() {
	telemetryProjector := NewHostTelemetryProjector(s.tx)

	discoveredCloudMock := mocks.NewDiscoveredCloudMock()
	discoveredHostMock := mocks.NewDiscoveredHostMock()

	agentDiscoveries := make(map[string]interface{})
	agentDiscoveries[CloudDiscovery] = discoveredCloudMock
	agentDiscoveries[HostDiscovery] = discoveredHostMock

	evtID := int64(1)

	for discoveryType, discoveredData := range agentDiscoveries {
		requestBody, _ := json.Marshal(discoveredData)

		telemetryProjector.Project(&DataCollectedEvent{
			ID:            evtID,
			AgentID:       "agent_id",
			DiscoveryType: discoveryType,
			Payload:       requestBody,
		})
		evtID++
	}

	var projectedTelemetry entities.HostTelemetry
	s.tx.First(&projectedTelemetry)

	s.Equal(discoveredHostMock.HostName, projectedTelemetry.HostName)
	s.Equal(discoveredHostMock.CPUCount, projectedTelemetry.CPUCount)
	s.Equal(discoveredHostMock.SocketCount, projectedTelemetry.SocketCount)
	s.Equal(discoveredHostMock.TotalMemoryMB, projectedTelemetry.TotalMemoryMB)

	s.Equal(discoveredHostMock.OSVersion, projectedTelemetry.SLESVersion)

	expectedCloudProvider := discoveredCloudMock.Provider
	s.Equal(expectedCloudProvider, projectedTelemetry.CloudProvider)
}
