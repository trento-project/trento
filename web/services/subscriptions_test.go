package services

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/entities"
	"github.com/trento-project/trento/web/models"
	"gorm.io/gorm"
)

type SubscriptionServiceTestSuite struct {
	suite.Suite
	db          *gorm.DB
	tx          *gorm.DB
	subsService *subscriptionsService
}

func TestSubscriptionsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionServiceTestSuite))
}

func (suite *SubscriptionServiceTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())

	suite.db.AutoMigrate(entities.SlesSubscription{}, entities.Host{})
	loadSubsFixtures(suite.db)
}

func (suite *SubscriptionServiceTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(entities.SlesSubscription{}, entities.Host{})
}

func (suite *SubscriptionServiceTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
	suite.subsService = NewSubscriptionsService(suite.tx)
}

func (suite *SubscriptionServiceTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func loadSubsFixtures(db *gorm.DB) {
	db.Create(&entities.SlesSubscription{
		AgentID:            "1",
		ID:                 "SLES_SAP",
		Version:            "15.2",
		Arch:               "x86_64",
		Status:             "Registered",
		StartsAt:           "2019-03-20 09:55:32 UTC",
		ExpiresAt:          "2024-03-20 09:55:32 UTC",
		SubscriptionStatus: "ACTIVE",
		Type:               "internal",
	})

	db.Create(&entities.SlesSubscription{
		AgentID: "1",
		ID:      "sle-module-public-cloud",
		Version: "15.2",
		Arch:    "x86_64",
		Status:  "Registered",
	})

	db.Create(&entities.SlesSubscription{
		AgentID:            "2",
		ID:                 "SLES_SAP",
		Version:            "15.2",
		Arch:               "x86_64",
		Status:             "Registered",
		StartsAt:           "2019-03-20 09:55:32 UTC",
		ExpiresAt:          "2024-03-20 09:55:32 UTC",
		SubscriptionStatus: "ACTIVE",
		Type:               "internal",
	})

	db.Create(&entities.Host{
		AgentID: "1",
		Name:    "host1",
	})
}

func (suite *SubscriptionServiceTestSuite) TestSubscriptionService_IsTrentoPremium() {
	premium, count, err := suite.subsService.IsTrentoPremium()

	suite.Equal(premium, true)
	suite.Equal(count, int64(2))
	suite.NoError(err)

	suite.tx.Where("id", "SLES_SAP").Delete(&entities.SlesSubscription{})

	premium, count, err = suite.subsService.IsTrentoPremium()

	suite.Equal(premium, false)
	suite.Equal(count, int64(0))
	suite.NoError(err)
}

func (suite *SubscriptionServiceTestSuite) TestSubscriptionService_GetHostSubscriptions() {
	subs, err := suite.subsService.GetHostSubscriptions("host1")
	expectedSubs := []*models.SlesSubscription{
		&models.SlesSubscription{
			ID:                 "SLES_SAP",
			Version:            "15.2",
			Arch:               "x86_64",
			Status:             "Registered",
			StartsAt:           "2019-03-20 09:55:32 UTC",
			ExpiresAt:          "2024-03-20 09:55:32 UTC",
			SubscriptionStatus: "ACTIVE",
			Type:               "internal",
		},
		&models.SlesSubscription{
			ID:      "sle-module-public-cloud",
			Version: "15.2",
			Arch:    "x86_64",
			Status:  "Registered",
		},
	}
	suite.ElementsMatch(expectedSubs, subs)
	suite.NoError(err)
}
