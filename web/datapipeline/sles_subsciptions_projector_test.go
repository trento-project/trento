package datapipeline

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	_ "github.com/trento-project/trento/test"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/entities"
	"gorm.io/gorm"
)

type SlesSubscriptionsProjectorTestSuite struct {
	suite.Suite
	db            *gorm.DB
	tx            *gorm.DB
	subsProjector *projector
}

func TestSlesSubscriptionsProjectorTestSuite(t *testing.T) {
	suite.Run(t, new(SlesSubscriptionsProjectorTestSuite))
}

func (suite *SlesSubscriptionsProjectorTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())

	suite.db.AutoMigrate(&Subscription{}, &entities.SlesSubscription{})
}

func (suite *SlesSubscriptionsProjectorTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(Subscription{}, entities.SlesSubscription{})
}

func (suite *SlesSubscriptionsProjectorTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
	suite.subsProjector = NewSlesSubscriptionsProjector(suite.tx)

	suite.tx.Create(&entities.SlesSubscription{
		AgentID:            "779cdd70-e9e2-58ca-b18a-bf3eb3f71244",
		ID:                 "SLES_SAP",
		Version:            "15.1",
		Arch:               "x86_64",
		Status:             "Registered",
		StartsAt:           "2017-03-20 09:55:32 UTC",
		ExpiresAt:          "2021-03-20 09:55:32 UTC",
		SubscriptionStatus: "ACTIVE",
		Type:               "internal",
	})
}

func (suite *SlesSubscriptionsProjectorTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func (suite *SlesSubscriptionsProjectorTestSuite) Test_SlesSubscriptionsProjector() {

	jsonFile, err := os.Open("./test/fixtures/discovery/subscriptions/expected_published_subscriptions_discovery.json")
	if err != nil {
		panic(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var dataCollectedEvent *DataCollectedEvent
	json.Unmarshal(byteValue, &dataCollectedEvent)

	subsProjector_SubscriptionDiscoveryHandler(dataCollectedEvent, suite.tx)

	var projectedSub entities.SlesSubscription
	suite.tx.Last(&projectedSub)

	expectedSub := entities.SlesSubscription{
		AgentID:            "779cdd70-e9e2-58ca-b18a-bf3eb3f71244",
		ID:                 "SLES_SAP",
		Version:            "15.2",
		Arch:               "x86_64",
		Status:             "Registered",
		StartsAt:           "2021-09-17 13:41:34 UTC",
		ExpiresAt:          "2026-09-17 13:41:34 UTC",
		SubscriptionStatus: "ACTIVE",
		Type:               "internal",
	}

	suite.Equal(expectedSub, projectedSub)
}
