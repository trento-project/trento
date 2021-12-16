package telemetry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/web/services"
)

type EngineTestSuite struct {
	suite.Suite
	dummyInstallationId    uuid.UUID
	mockedPublisher        *MockPublisher
	mockedExtractor        *MockExtractor
	mockedPremiumDetection *services.MockPremiumDetection
}

func TestEngineTestSuite(t *testing.T) {
	suite.Run(t, new(EngineTestSuite))
}

func (suite *EngineTestSuite) SetupSuite() {
	suite.dummyInstallationId = uuid.New()
	telemetryCollectionInterval = 50 * time.Millisecond
}

func (suite *EngineTestSuite) SetupTest() {
	suite.mockedPublisher = new(MockPublisher)
	suite.mockedExtractor = new(MockExtractor)
	suite.mockedPremiumDetection = new(services.MockPremiumDetection)
	suite.mockedPremiumDetection.On("CanPublishTelemetry").Return(true, nil)
}

// Test_ExtractsAndPublishesSingleTelemetry tests simple scenario of extracting and publishing single telemetry.
func (suite *EngineTestSuite) Test_ExtractsAndPublishesSingleTelemetry() {
	dummyInstallationId := suite.dummyInstallationId
	ctx, cancel := context.WithCancel(context.Background())
	callcount := 0

	suite.mockedPublisher.On("Publish", "dummy_1", dummyInstallationId, mock.Anything).Run(func(args mock.Arguments) {
		callcount++
		if callcount == 2 {
			cancel()
		}
	}).Return(nil)
	suite.mockedExtractor.On("Extract").Return(&struct{}{}, nil)

	registry := &TelemetryRegistry{
		"dummy_1": suite.mockedExtractor,
	}

	engine := NewEngine(
		dummyInstallationId,
		suite.mockedPublisher,
		registry,
		suite.mockedPremiumDetection,
	)

	engine.Start(ctx)

	suite.mockedPublisher.AssertNumberOfCalls(suite.T(), "Publish", 2)
	suite.mockedExtractor.AssertNumberOfCalls(suite.T(), "Extract", 2)
}

// Test_ExtractsAndPublishesMultipleTelemetry tests the scenarion of extracting and publishing multiple telemetry.
func (suite *EngineTestSuite) Test_ExtractsAndPublishesMultipleTelemetry() {
	dummyInstallationId := suite.dummyInstallationId
	ctx, cancel := context.WithCancel(context.Background())
	callcount := 0

	suite.mockedPublisher.On("Publish", "dummy_1", dummyInstallationId, mock.Anything).Return(nil)
	suite.mockedPublisher.On("Publish", "dummy_2", dummyInstallationId, mock.Anything).Return(nil)
	suite.mockedPublisher.On("Publish", "dummy_3", dummyInstallationId, mock.Anything).Return(nil)
	suite.mockedPublisher.On("Publish", "dummy_4", dummyInstallationId, mock.Anything).Return(nil)
	suite.mockedPublisher.On("Publish", "dummy_5", dummyInstallationId, mock.Anything).Run(func(args mock.Arguments) {
		callcount++
		if callcount == 1 {
			cancel()
		}
	}).Return(nil)
	suite.mockedExtractor.On("Extract").Return(&struct{}{}, nil)

	registry := &TelemetryRegistry{
		"dummy_1": suite.mockedExtractor,
		"dummy_2": suite.mockedExtractor,
		"dummy_3": suite.mockedExtractor,
		"dummy_4": suite.mockedExtractor,
		"dummy_5": suite.mockedExtractor,
	}

	engine := NewEngine(
		dummyInstallationId,
		suite.mockedPublisher,
		registry,
		suite.mockedPremiumDetection,
	)

	engine.Start(ctx)

	suite.mockedPublisher.AssertNumberOfCalls(suite.T(), "Publish", 5)
	suite.mockedExtractor.AssertNumberOfCalls(suite.T(), "Extract", 5)
}

// Test_ExtractsAndPublishesAlsoIdentifiedTelemetries tests the scenario of supporting identified extractors.
func (suite *EngineTestSuite) Test_ExtractsAndPublishesAlsoIdentifiedExtractors() {
	dummyInstallationId := suite.dummyInstallationId
	ctx, cancel := context.WithCancel(context.Background())
	callcount := 0

	mockedIdentifiedExtractor := new(MockIdentifiedExtractor)

	suite.mockedPublisher.On("Publish", "dummy_1", dummyInstallationId, mock.Anything).Return(nil)
	suite.mockedPublisher.On("Publish", "dummy_2", dummyInstallationId, mock.Anything).Return(nil)
	suite.mockedPublisher.On("Publish", "dummy_3", dummyInstallationId, mock.Anything).Return(nil)
	suite.mockedPublisher.On("Publish", "dummy_4", dummyInstallationId, mock.Anything).Return(nil)
	suite.mockedPublisher.On("Publish", "dummy_5", dummyInstallationId, mock.Anything).Run(func(args mock.Arguments) {
		callcount++
		if callcount == 1 {
			cancel()
		}
	}).Return(nil)

	suite.mockedExtractor.On("Extract").Return(&struct{}{}, nil)
	mockedIdentifiedExtractor.On("Extract").Return(&struct{}{}, nil)
	mockedIdentifiedExtractor.On("WithInstallationID", dummyInstallationId)

	registry := &TelemetryRegistry{
		"dummy_1": suite.mockedExtractor,
		"dummy_2": suite.mockedExtractor,
		"dummy_3": mockedIdentifiedExtractor,
		"dummy_4": mockedIdentifiedExtractor,
		"dummy_5": suite.mockedExtractor,
	}

	engine := NewEngine(
		dummyInstallationId,
		suite.mockedPublisher,
		registry,
		suite.mockedPremiumDetection,
	)

	engine.Start(ctx)

	suite.mockedPublisher.AssertNumberOfCalls(suite.T(), "Publish", 5)
	suite.mockedExtractor.AssertNumberOfCalls(suite.T(), "Extract", 3)
	mockedIdentifiedExtractor.AssertNumberOfCalls(suite.T(), "Extract", 2)
	mockedIdentifiedExtractor.AssertNumberOfCalls(suite.T(), "WithInstallationID", 2)
}

// Test_ExtractsAndPublishesAlsoWithSomeErrors tests the scenario of handling errors during extraction.
// That means that the engine will not publish if an error occurs in the extraction.
func (suite *EngineTestSuite) Test_ExtractsAndPublishesAlsoWithSomeErrors() {
	dummyInstallationId := suite.dummyInstallationId
	ctx, cancel := context.WithCancel(context.Background())
	callcount := 0

	mockederroringExtractor := new(MockExtractor)

	suite.mockedPublisher.On("Publish", "dummy_1", dummyInstallationId, mock.Anything).Return(nil)
	suite.mockedPublisher.On("Publish", "dummy_2", dummyInstallationId, mock.Anything).Return(nil)
	suite.mockedPublisher.On("Publish", "dummy_3", dummyInstallationId, mock.Anything).Run(func(args mock.Arguments) {
		callcount++
		if callcount == 1 {
			cancel()
		}
	}).Return(nil)

	suite.mockedExtractor.On("Extract").Return(&struct{}{}, nil)
	mockederroringExtractor.On("Extract").Return(nil, errors.New("dummy error"))

	registry := &TelemetryRegistry{
		"dummy_1": suite.mockedExtractor,
		"dummy_2": mockederroringExtractor,
		"dummy_3": suite.mockedExtractor,
	}

	engine := NewEngine(
		dummyInstallationId,
		suite.mockedPublisher,
		registry,
		suite.mockedPremiumDetection,
	)

	engine.Start(ctx)

	suite.mockedPublisher.AssertNumberOfCalls(suite.T(), "Publish", 2)
	suite.mockedExtractor.AssertNumberOfCalls(suite.T(), "Extract", 2)
	mockederroringExtractor.AssertNumberOfCalls(suite.T(), "Extract", 1)
}

func (suite *EngineTestSuite) Test_SkippingTelemetry() {
	ctx, cancel := context.WithCancel(context.Background())

	mockedPremiumDetection := new(services.MockPremiumDetection)
	mockedPremiumDetection.On("CanPublishTelemetry").Return(false, nil)

	registry := &TelemetryRegistry{
		"dummy_1": suite.mockedExtractor,
	}

	engine := NewEngine(
		suite.dummyInstallationId,
		suite.mockedPublisher,
		registry,
		mockedPremiumDetection,
	)

	ch := make(chan struct{})
	go func() {
		engine.Start(ctx)
		ch <- struct{}{}
	}()
	go cancel()
	<-ch

	suite.mockedPublisher.AssertNumberOfCalls(suite.T(), "Publish", 0)
	suite.mockedExtractor.AssertNumberOfCalls(suite.T(), "Extract", 0)
}
