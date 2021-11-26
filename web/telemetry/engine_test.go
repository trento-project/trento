package telemetry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type EngineTestSuite struct {
	suite.Suite
	dummyInstallationId uuid.UUID
}

func TestEngineTestSuite(t *testing.T) {
	suite.Run(t, new(EngineTestSuite))
}

func (suite *EngineTestSuite) SetupSuite() {
	suite.dummyInstallationId = uuid.New()
	telemetryCollectionInterval = 50 * time.Millisecond
}

// Test_ExtractsAndPublishesSingleTelemetry tests simple scenario of extracting and publishing single telemetry.
func (suite *EngineTestSuite) Test_ExtractsAndPublishesSingleTelemetry() {
	dummyInstallationId := suite.dummyInstallationId
	ctx, cancel := context.WithCancel(context.Background())
	callcount := 0

	mockedPublisher := new(MockPublisher)
	mockedExtractor := new(MockExtractor)

	mockedPublisher.On("Publish", "dummy_1", dummyInstallationId, mock.Anything).Run(func(args mock.Arguments) {
		callcount++
		if callcount == 2 {
			cancel()
		}
	}).Return(nil)
	mockedExtractor.On("Extract").Return(&struct{}{}, nil)

	registry := &TelemetryRegistry{
		"dummy_1": mockedExtractor,
	}

	engine := NewEngine(
		dummyInstallationId,
		mockedPublisher,
		registry,
	)

	engine.Start(ctx)

	mockedPublisher.AssertNumberOfCalls(suite.T(), "Publish", 2)
	mockedExtractor.AssertNumberOfCalls(suite.T(), "Extract", 2)
}

// Test_ExtractsAndPublishesMultipleTelemetry tests the scenarion of extracting and publishing multiple telemetry.
func (suite *EngineTestSuite) Test_ExtractsAndPublishesMultipleTelemetry() {
	dummyInstallationId := suite.dummyInstallationId
	ctx, cancel := context.WithCancel(context.Background())
	callcount := 0

	mockedPublisher := new(MockPublisher)
	mockedExtractor := new(MockExtractor)

	mockedPublisher.On("Publish", "dummy_1", dummyInstallationId, mock.Anything).Return(nil)
	mockedPublisher.On("Publish", "dummy_2", dummyInstallationId, mock.Anything).Return(nil)
	mockedPublisher.On("Publish", "dummy_3", dummyInstallationId, mock.Anything).Return(nil)
	mockedPublisher.On("Publish", "dummy_4", dummyInstallationId, mock.Anything).Return(nil)
	mockedPublisher.On("Publish", "dummy_5", dummyInstallationId, mock.Anything).Run(func(args mock.Arguments) {
		callcount++
		if callcount == 1 {
			cancel()
		}
	}).Return(nil)
	mockedExtractor.On("Extract").Return(&struct{}{}, nil)

	registry := &TelemetryRegistry{
		"dummy_1": mockedExtractor,
		"dummy_2": mockedExtractor,
		"dummy_3": mockedExtractor,
		"dummy_4": mockedExtractor,
		"dummy_5": mockedExtractor,
	}

	engine := NewEngine(
		dummyInstallationId,
		mockedPublisher,
		registry,
	)

	engine.Start(ctx)

	mockedPublisher.AssertNumberOfCalls(suite.T(), "Publish", 5)
	mockedExtractor.AssertNumberOfCalls(suite.T(), "Extract", 5)
}

// Test_ExtractsAndPublishesAlsoIdentifiedTelemetries tests the scenario of supporting identified extractors.
func (suite *EngineTestSuite) Test_ExtractsAndPublishesAlsoIdentifiedExtractors() {
	dummyInstallationId := suite.dummyInstallationId
	ctx, cancel := context.WithCancel(context.Background())
	callcount := 0

	mockedPublisher := new(MockPublisher)
	mockedExtractor := new(MockExtractor)
	mockedIdentifiedExtractor := new(MockIdentifiedExtractor)

	mockedPublisher.On("Publish", "dummy_1", dummyInstallationId, mock.Anything).Return(nil)
	mockedPublisher.On("Publish", "dummy_2", dummyInstallationId, mock.Anything).Return(nil)
	mockedPublisher.On("Publish", "dummy_3", dummyInstallationId, mock.Anything).Return(nil)
	mockedPublisher.On("Publish", "dummy_4", dummyInstallationId, mock.Anything).Return(nil)
	mockedPublisher.On("Publish", "dummy_5", dummyInstallationId, mock.Anything).Run(func(args mock.Arguments) {
		callcount++
		if callcount == 1 {
			cancel()
		}
	}).Return(nil)

	mockedExtractor.On("Extract").Return(&struct{}{}, nil)
	mockedIdentifiedExtractor.On("Extract").Return(&struct{}{}, nil)
	mockedIdentifiedExtractor.On("WithInstallationID", dummyInstallationId)

	registry := &TelemetryRegistry{
		"dummy_1": mockedExtractor,
		"dummy_2": mockedExtractor,
		"dummy_3": mockedIdentifiedExtractor,
		"dummy_4": mockedIdentifiedExtractor,
		"dummy_5": mockedExtractor,
	}

	engine := NewEngine(
		dummyInstallationId,
		mockedPublisher,
		registry,
	)

	engine.Start(ctx)

	mockedPublisher.AssertNumberOfCalls(suite.T(), "Publish", 5)
	mockedExtractor.AssertNumberOfCalls(suite.T(), "Extract", 3)
	mockedIdentifiedExtractor.AssertNumberOfCalls(suite.T(), "Extract", 2)
	mockedIdentifiedExtractor.AssertNumberOfCalls(suite.T(), "WithInstallationID", 2)
}

// Test_ExtractsAndPublishesAlsoWithSomeErrors tests the scenario of handling errors during extraction.
// That means that the engine will not publish if an error occurs in the extraction.
func (suite *EngineTestSuite) Test_ExtractsAndPublishesAlsoWithSomeErrors() {
	dummyInstallationId := suite.dummyInstallationId
	ctx, cancel := context.WithCancel(context.Background())
	callcount := 0

	mockedPublisher := new(MockPublisher)
	mockedExtractor := new(MockExtractor)
	mockederroringExtractor := new(MockExtractor)

	mockedPublisher.On("Publish", "dummy_1", dummyInstallationId, mock.Anything).Return(nil)
	mockedPublisher.On("Publish", "dummy_2", dummyInstallationId, mock.Anything).Return(nil)
	mockedPublisher.On("Publish", "dummy_3", dummyInstallationId, mock.Anything).Run(func(args mock.Arguments) {
		callcount++
		if callcount == 1 {
			cancel()
		}
	}).Return(nil)

	mockedExtractor.On("Extract").Return(&struct{}{}, nil)
	mockederroringExtractor.On("Extract").Return(nil, errors.New("dummy error"))

	registry := &TelemetryRegistry{
		"dummy_1": mockedExtractor,
		"dummy_2": mockederroringExtractor,
		"dummy_3": mockedExtractor,
	}

	engine := NewEngine(
		dummyInstallationId,
		mockedPublisher,
		registry,
	)

	engine.Start(ctx)

	mockedPublisher.AssertNumberOfCalls(suite.T(), "Publish", 2)
	mockedExtractor.AssertNumberOfCalls(suite.T(), "Extract", 2)
	mockederroringExtractor.AssertNumberOfCalls(suite.T(), "Extract", 1)
}
