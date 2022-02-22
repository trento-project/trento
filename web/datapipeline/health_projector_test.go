package datapipeline

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	_ "github.com/trento-project/trento/test"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/entities"
	"gorm.io/gorm"
)

type HealthProjectorTestSuite struct {
	suite.Suite
	db *gorm.DB
	tx *gorm.DB
}

func TestHealthProjectorTestSuite(t *testing.T) {
	suite.Run(t, new(HealthProjectorTestSuite))
}

func (suite *HealthProjectorTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())

	suite.db.AutoMigrate(entities.HealthState{})
}

func (suite *HealthProjectorTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(entities.HealthState{})
}

func (suite *HealthProjectorTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
}

func (suite *HealthProjectorTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func (suite *HealthProjectorTestSuite) Test_ProjectHealth() {
	err := ProjectHealth(suite.tx, "1", "my_health_value", "passing")
	suite.NoError(err)

	var health entities.HealthState
	suite.tx.First(&health)

	var partialHealth map[string]string
	json.Unmarshal(health.PartialHealths, &partialHealth)

	suite.Equal("1", health.ID)
	suite.Equal("passing", health.Health)
	suite.Equal(map[string]string{"my_health_value": "passing"}, partialHealth)
}

func (suite *HealthProjectorTestSuite) Test_ProjectHealth_Update() {
	partialHealths1, _ := json.Marshal(map[string]string{"my_health_value": "passing"})
	suite.tx.Create(&entities.HealthState{
		ID:             "1",
		Health:         "passing",
		PartialHealths: partialHealths1,
	})

	err := ProjectHealth(suite.tx, "1", "my_health_value", "critical")
	suite.NoError(err)

	var health entities.HealthState
	suite.tx.First(&health)

	var partialHealth map[string]string
	json.Unmarshal(health.PartialHealths, &partialHealth)

	suite.Equal("1", health.ID)
	suite.Equal("critical", health.Health)
	suite.Equal(map[string]string{"my_health_value": "critical"}, partialHealth)
}

func (suite *HealthProjectorTestSuite) Test_ProjectHealth_New() {
	partialHealths1, _ := json.Marshal(map[string]string{"my_health_value": "passing"})
	suite.tx.Create(&entities.HealthState{
		ID:             "1",
		Health:         "passing",
		PartialHealths: partialHealths1,
	})

	err := ProjectHealth(suite.tx, "1", "my_new_health", "warning")
	suite.NoError(err)

	var health entities.HealthState
	suite.tx.First(&health)

	var partialHealth map[string]string
	json.Unmarshal(health.PartialHealths, &partialHealth)

	suite.Equal("1", health.ID)
	suite.Equal("warning", health.Health)
	suite.Equal(
		map[string]string{
			"my_health_value": "passing",
			"my_new_health":   "warning",
		},
		partialHealth,
	)
}

func (suite *HealthProjectorTestSuite) Test_ComputeOverallHealth_Passing() {
	health := computeOverallHealth(
		map[string]string{
			"health":       "passing",
			"other_health": "passing",
		},
	)
	suite.Equal("passing", health)
}

func (suite *HealthProjectorTestSuite) Test_ComputeOverallHealth_Unknown() {
	health := computeOverallHealth(
		map[string]string{
			"health":         "passing",
			"unknonw_health": "unknown",
			"other_health":   "passing",
		},
	)
	suite.Equal("unknown", health)
}

func (suite *HealthProjectorTestSuite) Test_ComputeOverallHealth_Warning() {
	health := computeOverallHealth(
		map[string]string{
			"health":         "passing",
			"warning_health": "warning",
			"unknonw_health": "unknown",
			"other_health":   "passing",
		},
	)
	suite.Equal("warning", health)
}

func (suite *HealthProjectorTestSuite) Test_ComputeOverallHealth_Critical() {
	health := computeOverallHealth(
		map[string]string{
			"health":          "passing",
			"critical_health": "critical",
			"warning_health":  "warning",
			"unknonw_health":  "unknown",
			"other_health":    "passing",
		},
	)
	suite.Equal("critical", health)
}
