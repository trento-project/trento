package services

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/entities"
	"gorm.io/gorm"
)

type SettingsServiceTestSuite struct {
	suite.Suite
	db              *gorm.DB
	tx              *gorm.DB
	settingsService SettingsService
}

func TestSettingsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SettingsServiceTestSuite))
}

func (suite *SettingsServiceTestSuite) SetupSuite() {
	suite.db = helpers.SetupTestDatabase(suite.T())

	suite.db.AutoMigrate(entities.Settings{})
}

func (suite *SettingsServiceTestSuite) TearDownSuite() {
	suite.db.Migrator().DropTable(entities.Settings{})
}

func (suite *SettingsServiceTestSuite) SetupTest() {
	suite.tx = suite.db.Begin()
	suite.settingsService = NewSettingsService(suite.tx)
}

func (suite *SettingsServiceTestSuite) TearDownTest() {
	suite.tx.Rollback()
}

func (suite *SettingsServiceTestSuite) TestSettingsService_InitializesNewInstallation() {
	var settings entities.Settings
	suite.tx.Find(&settings)

	suite.Empty(settings.InstallationID)

	installationID, err := suite.settingsService.InitializeIdentifier()

	suite.tx.Find(&settings)

	suite.NoError(err)
	suite.NotEmpty(settings.InstallationID)
	suite.EqualValues(installationID.String(), settings.InstallationID)
}

func (suite *SettingsServiceTestSuite) TestSettingsService_DetectsTrentoAlreadyInstalled() {
	const dummyInstallationID = "59fd8017-b7fd-477b-9ebe-b658c558f3e9"

	settings := entities.Settings{
		InstallationID: dummyInstallationID,
	}
	suite.tx.Create(&settings)

	installationID, err := suite.settingsService.InitializeIdentifier()

	suite.NoError(err)
	suite.EqualValues(dummyInstallationID, installationID.String())
}
