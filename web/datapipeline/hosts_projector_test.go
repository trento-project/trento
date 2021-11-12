package datapipeline

import (
	"encoding/json"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/agent/discovery/mocks"
	_ "github.com/trento-project/trento/test"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/entities"
	"gorm.io/gorm"
)

type HostsProjectorTestSuite struct {
	suite.Suite
	db *gorm.DB
	tx *gorm.DB
}

func TestHostsProjectorTestSuite(t *testing.T) {
	suite.Run(t, new(HostsProjectorTestSuite))
}

func (suite *HostsProjectorTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())

	suite.db.AutoMigrate(&Subscription{}, &entities.Host{})
}

func (suite *HostsProjectorTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(Subscription{}, entities.Host{})
}

func (suite *HostsProjectorTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
}

func (suite *HostsProjectorTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

// Test_HostDiscoveryHandler tests the HostDiscoveryHandler function execution on a HostDiscovery published by an agent
func (s *HostsProjectorTestSuite) Test_HostDiscoveryHandler() {
	discoveredHostMock := mocks.NewDiscoveredHostMock()

	requestBody, _ := json.Marshal(discoveredHostMock)

	hostsProjector_HostDiscoveryHandler(&DataCollectedEvent{
		ID:            1,
		AgentID:       "agent_id",
		DiscoveryType: HostDiscovery,
		Payload:       requestBody,
	}, s.tx)

	var projectedHost entities.Host
	s.tx.First(&projectedHost)

	s.Equal(discoveredHostMock.HostName, projectedHost.Name)
	s.EqualValues(discoveredHostMock.HostIpAddresses, projectedHost.IPAddresses)
	s.Equal(discoveredHostMock.AgentVersion, projectedHost.AgentVersion)

	s.Equal("", projectedHost.CloudProvider)
	s.Equal("", projectedHost.ClusterID)
	s.Equal("", projectedHost.ClusterName)
	s.Equal(pq.StringArray(nil), projectedHost.SIDs)
}

// Test_CloudDiscoveryHandler tests the loudDiscoveryHandler function execution on a CloudDiscovery published by an agent
func (s *HostsProjectorTestSuite) Test_CloudDiscoveryHandler() {
	discoveredCloudMock := mocks.NewDiscoveredCloudMock()

	requestBody, _ := json.Marshal(discoveredCloudMock)

	hostsProjector_CloudDiscoveryHandler(&DataCollectedEvent{
		ID:            1,
		AgentID:       "agent_id",
		DiscoveryType: CloudDiscovery,
		Payload:       requestBody,
	}, s.tx)

	var projectedHost entities.Host
	s.tx.First(&projectedHost)

	s.Equal(discoveredCloudMock.Provider, projectedHost.CloudProvider)

	s.Equal("", projectedHost.Name)
	s.Equal("", projectedHost.ClusterID)
	s.Equal("", projectedHost.ClusterName)
	s.Equal(pq.StringArray(nil), projectedHost.SIDs)
}

// Test_ClusterDiscoveryHandler tests the ClusterDiscoveryHandler function execution on a ClusterDiscovery published by an agent
func (s *HostsProjectorTestSuite) Test_ClusterDiscoveryHandler() {
	discoveredClusterMock := mocks.NewDiscoveredClusterMock()

	requestBody, _ := json.Marshal(discoveredClusterMock)
	hostsProjector_ClusterDiscoveryHandler(&DataCollectedEvent{
		ID:            1,
		AgentID:       "agent_id",
		DiscoveryType: ClusterDiscovery,
		Payload:       requestBody,
	}, s.tx)

	var projectedHost entities.Host
	s.tx.First(&projectedHost)

	s.Equal("47d1190ffb4f781974c8356d7f863b03", projectedHost.ClusterID)
	s.Equal("hana_cluster", projectedHost.ClusterName)

	s.Equal("", projectedHost.Name)
	s.Equal("", projectedHost.CloudProvider)
	s.Equal(pq.StringArray(nil), projectedHost.SIDs)
}

// Test_SAPSystemDiscoveryHandler tests the SAPSystemDiscoveryHandler function execution on a SAPSystemDiscovery published by an agent
func (s *HostsProjectorTestSuite) Test_SAPSystemDiscoveryHandler() {
	discoveredSAPSystemMock := mocks.NewDiscoveredSAPSystemMock()

	requestBody, _ := json.Marshal(discoveredSAPSystemMock)
	hostsProjector_SAPSystemDiscoveryHandler(&DataCollectedEvent{
		ID:            1,
		AgentID:       "agent_id",
		DiscoveryType: SAPsystemDiscovery,
		Payload:       requestBody,
	}, s.tx)

	var projectedHost entities.Host
	s.tx.First(&projectedHost)

	s.EqualValues([]string{"PRD"}, projectedHost.SIDs)

	s.Equal("", projectedHost.Name)
	s.Equal("", projectedHost.ClusterID)
	s.Equal("", projectedHost.ClusterName)
	s.Equal("", projectedHost.CloudProvider)
}

// Test_HostsProjector tests the HostsProjector projects all of the discoveries it is interested in, resulting in a single host readmodel
func (s *HostsProjectorTestSuite) Test_TelemetryProjector() {
	hostsProjector := NewHostsProjector(s.tx)

	discoveredHostMock := mocks.NewDiscoveredHostMock()
	discoveredCloudMock := mocks.NewDiscoveredCloudMock()
	discoveredClusterMock := mocks.NewDiscoveredClusterMock()
	discoveredSAPSystemMock := mocks.NewDiscoveredSAPSystemMock()

	agentDiscoveries := make(map[string]interface{})
	agentDiscoveries[HostDiscovery] = discoveredHostMock
	agentDiscoveries[CloudDiscovery] = discoveredCloudMock
	agentDiscoveries[ClusterDiscovery] = discoveredClusterMock
	agentDiscoveries[SAPsystemDiscovery] = discoveredSAPSystemMock

	evtID := int64(1)

	for discoveryType, discoveredData := range agentDiscoveries {
		requestBody, _ := json.Marshal(discoveredData)

		hostsProjector.Project(&DataCollectedEvent{
			ID:            evtID,
			AgentID:       "agent_id",
			DiscoveryType: discoveryType,
			Payload:       requestBody,
		})
		evtID++
	}

	var projectedHost entities.Host
	s.tx.First(&projectedHost)

	s.Equal(discoveredHostMock.HostName, projectedHost.Name)
	s.EqualValues(discoveredHostMock.HostIpAddresses, projectedHost.IPAddresses)
	s.Equal(discoveredCloudMock.Provider, projectedHost.CloudProvider)
	s.Equal(discoveredClusterMock.Id, projectedHost.ClusterID)
	s.Equal(discoveredClusterMock.Name, projectedHost.ClusterName)

	for _, m := range discoveredSAPSystemMock {
		s.Contains(projectedHost.SIDs, m.SID)
	}
}

func (s *HostsProjectorTestSuite) Test_filterIPAddresses() {
	ipAddresses := []string{
		"127.0.0.1",
		"10.1.74.5",
		"::1",
		"fe80::6245:bdff:fe8b:5896",
		"not_valid",
	}

	s.EqualValues([]string{"10.1.74.5"}, filterIPAddresses(ipAddresses))
}
