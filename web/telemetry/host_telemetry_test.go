package telemetry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
)

type HostTelemetryTestSuite struct {
	suite.Suite
	db *gorm.DB
	tx *gorm.DB
}

func TestHostTelemetryTestSuite(t *testing.T) {
	suite.Run(t, new(HostTelemetryTestSuite))
}

func (suite *HostTelemetryTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())

	suite.db.AutoMigrate(&models.HostTelemetry{})
}

func (suite *HostTelemetryTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(models.HostTelemetry{})
}

func (suite *HostTelemetryTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
}

func (suite *HostTelemetryTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

// Test_ExtractsExpectedHostTelemetry tests that given some projected host telemetry, it extracts correctly the expected HostTelemetries
func (suite *HostTelemetryTestSuite) Test_ExtractsExpectedHostTelemetry() {
	fixtures := hostTelemetryFixtures()
	suite.tx.Create(fixtures)

	hostTelemetryExtractor := NewHostTelemetryExtractor(suite.tx)

	extracted, _ := hostTelemetryExtractor.Extract()
	extractedHostTelemetry, _ := extracted.(HostTelemetries)

	suite.EqualValues(len(fixtures), len(extractedHostTelemetry))
	suite.EqualValues(fixtures[0].AgentID, extractedHostTelemetry[0].AgentID)
	suite.EqualValues(fixtures[0].SLESVersion, extractedHostTelemetry[0].SLESVersion)
	suite.EqualValues(fixtures[0].CPUCount, extractedHostTelemetry[0].CPUCount)
	suite.EqualValues(fixtures[0].SocketCount, extractedHostTelemetry[0].SocketCount)
	suite.EqualValues(fixtures[0].TotalMemoryMB, extractedHostTelemetry[0].TotalMemoryMB)
	suite.EqualValues(fixtures[0].CloudProvider, extractedHostTelemetry[0].CloudProvider)
	suite.True(fixtures[0].UpdatedAt.Equal(extractedHostTelemetry[0].Time))

	suite.EqualValues(fixtures[1].AgentID, extractedHostTelemetry[1].AgentID)
	suite.EqualValues(fixtures[1].SLESVersion, extractedHostTelemetry[1].SLESVersion)
	suite.EqualValues(fixtures[1].CPUCount, extractedHostTelemetry[1].CPUCount)
	suite.EqualValues(fixtures[1].SocketCount, extractedHostTelemetry[1].SocketCount)
	suite.EqualValues(fixtures[1].TotalMemoryMB, extractedHostTelemetry[1].TotalMemoryMB)
	suite.EqualValues(fixtures[1].CloudProvider, extractedHostTelemetry[1].CloudProvider)
	suite.True(fixtures[1].UpdatedAt.Equal(extractedHostTelemetry[1].Time))
}

// Test_ExtractsEmptyHostTelemetry tests that given an empty set of projected host telemetry, it extracts correctly nothing
func (suite *HostTelemetryTestSuite) Test_ExtractsEmptyHostTelemetry() {
	hostTelemetryExtractor := NewHostTelemetryExtractor(suite.tx)

	extracted, err := hostTelemetryExtractor.Extract()

	suite.Error(err)
	suite.Nil(extracted)
}

func hostTelemetryFixtures() []models.HostTelemetry {
	t1 := time.Date(2020, 11, 12, 11, 45, 26, 0, time.UTC)
	t2 := time.Date(2020, 11, 13, 11, 45, 26, 0, time.UTC)

	return []models.HostTelemetry{
		{
			AgentID:       "some-agent-id",
			SLESVersion:   "15-sp2",
			CPUCount:      2,
			SocketCount:   8,
			TotalMemoryMB: 4096,
			CloudProvider: "azure",
			UpdatedAt:     t1,
		},
		{
			AgentID:       "another-agent-id",
			SLESVersion:   "15-sp2",
			CPUCount:      4,
			SocketCount:   8,
			TotalMemoryMB: 8192,
			CloudProvider: "azure",
			UpdatedAt:     t2,
		},
	}
}
