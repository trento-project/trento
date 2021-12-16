package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PremiumDetectionTestSuite struct {
	suite.Suite
	subscriptions *MockSubscriptionsService
	settings      *MockSettingsService
}

func TestEngineTestSuite(t *testing.T) {
	suite.Run(t, new(PremiumDetectionTestSuite))
}

func (suite *PremiumDetectionTestSuite) SetupSuite() {
}

func (suite *PremiumDetectionTestSuite) SetupTest() {
	suite.subscriptions = new(MockSubscriptionsService)
	suite.settings = new(MockSettingsService)
}

func (suite *PremiumDetectionTestSuite) Test_DoesNotRequireEulaAcceptanceOnCommunityFlavor() {
	premiumDetection := NewPremiumDetection(
		Community,
		suite.subscriptions,
		suite.settings,
	)

	requiresEulaAcceptance, err := premiumDetection.RequiresEulaAcceptance()
	suite.NoError(err)
	suite.False(requiresEulaAcceptance)
}

func (suite *PremiumDetectionTestSuite) Test_RequiresEulaAcceptanceOnPremiumFlavor() {
	suite.settings.On("IsEulaAccepted").Return(false, nil)

	premiumDetection := NewPremiumDetection(
		Premium,
		suite.subscriptions,
		suite.settings,
	)

	requiresEulaAcceptance, err := premiumDetection.RequiresEulaAcceptance()
	suite.NoError(err)
	suite.True(requiresEulaAcceptance)
	suite.settings.AssertExpectations(suite.T())
}

func (suite *PremiumDetectionTestSuite) Test_DoesNotRequireEulaAcceptanceOnPremiumFlavor() {
	suite.settings.On("IsEulaAccepted").Return(true, nil)

	premiumDetection := NewPremiumDetection(
		Premium,
		suite.subscriptions,
		suite.settings,
	)

	requiresEulaAcceptance, err := premiumDetection.RequiresEulaAcceptance()
	suite.NoError(err)
	suite.False(requiresEulaAcceptance)
	suite.settings.AssertExpectations(suite.T())
}

func (suite *PremiumDetectionTestSuite) Test_FailsDeterminingEulaAcceptanceRequirement() {
	suite.settings.On("IsEulaAccepted").Return(false, errors.New("BOO BOO"))

	premiumDetection := NewPremiumDetection(
		Premium,
		suite.subscriptions,
		suite.settings,
	)

	requiresEulaAcceptance, err := premiumDetection.RequiresEulaAcceptance()
	suite.Error(err, "BOO BOO")
	suite.False(requiresEulaAcceptance)
	suite.settings.AssertExpectations(suite.T())
}

func (suite *PremiumDetectionTestSuite) Test_CannotPublishTelemetryOnCommunityFlavor() {
	premiumDetection := NewPremiumDetection(
		Community,
		suite.subscriptions,
		suite.settings,
	)

	canPublishTelemetry, err := premiumDetection.CanPublishTelemetry()
	suite.NoError(err)
	suite.False(canPublishTelemetry)
}

func (suite *PremiumDetectionTestSuite) Test_CanPublishTelemetryOnPremiumFlavor() {
	suite.settings.On("IsEulaAccepted").Return(true, nil)

	premiumDetection := NewPremiumDetection(
		Premium,
		suite.subscriptions,
		suite.settings,
	)

	canPublishTelemetry, err := premiumDetection.CanPublishTelemetry()
	suite.NoError(err)
	suite.True(canPublishTelemetry)
	suite.settings.AssertExpectations(suite.T())
}

func (suite *PremiumDetectionTestSuite) Test_CannotPublishTelemetryOnPremiumFlavor() {
	suite.settings.On("IsEulaAccepted").Return(false, nil)

	premiumDetection := NewPremiumDetection(
		Premium,
		suite.subscriptions,
		suite.settings,
	)

	canPublishTelemetry, err := premiumDetection.CanPublishTelemetry()
	suite.NoError(err)
	suite.False(canPublishTelemetry)
	suite.settings.AssertExpectations(suite.T())
}

func (suite *PremiumDetectionTestSuite) Test_FailsDeterminingTelemetryPublishability() {
	suite.settings.On("IsEulaAccepted").Return(false, errors.New("KABOOM"))

	premiumDetection := NewPremiumDetection(
		Premium,
		suite.subscriptions,
		suite.settings,
	)

	canPublishTelemetry, err := premiumDetection.CanPublishTelemetry()
	suite.Error(err, "KABOOM")
	suite.False(canPublishTelemetry)
	suite.settings.AssertExpectations(suite.T())
}

func (suite *PremiumDetectionTestSuite) Test_PremiumIsNotActiveOnCommunityFlavor() {
	premiumDetection := NewPremiumDetection(
		Community,
		suite.subscriptions,
		suite.settings,
	)

	isPremiumActive, err := premiumDetection.IsPremiumActive()
	suite.NoError(err)
	suite.False(isPremiumActive)
}

func (suite *PremiumDetectionTestSuite) Test_PremiumIsNotActiveOnPremiumFlavor() {
	suite.subscriptions.On("IsTrentoPremium").Return(false, nil)

	premiumDetection := NewPremiumDetection(
		Premium,
		suite.subscriptions,
		suite.settings,
	)

	isPremiumActive, err := premiumDetection.IsPremiumActive()
	suite.NoError(err)
	suite.False(isPremiumActive)
	suite.subscriptions.AssertExpectations(suite.T())
}

func (suite *PremiumDetectionTestSuite) Test_PremiumIsActiveOnPremiumFlavor() {
	suite.subscriptions.On("IsTrentoPremium").Return(true, nil)

	premiumDetection := NewPremiumDetection(
		Premium,
		suite.subscriptions,
		suite.settings,
	)

	isPremiumActive, err := premiumDetection.IsPremiumActive()
	suite.NoError(err)
	suite.True(isPremiumActive)
	suite.subscriptions.AssertExpectations(suite.T())
}

func (suite *PremiumDetectionTestSuite) Test_FailsDeterminingPremiumIsActive() {
	suite.subscriptions.On("IsTrentoPremium").Return(false, errors.New("SOME ERROR"))

	premiumDetection := NewPremiumDetection(
		Premium,
		suite.subscriptions,
		suite.settings,
	)

	isPremiumActive, err := premiumDetection.IsPremiumActive()
	suite.Error(err, "SOME ERROR")
	suite.False(isPremiumActive)
	suite.subscriptions.AssertExpectations(suite.T())
}
