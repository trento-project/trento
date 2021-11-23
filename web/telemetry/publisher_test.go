package telemetry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
)

type PublisherTestSuite struct {
	suite.Suite
}

func TestPublisherTestSuite(t *testing.T) {
	suite.Run(t, new(PublisherTestSuite))
}

func (suite *PublisherTestSuite) SetupSuite() {
	apiHost = "https://httpbin.org/anything"
}

// Test_PublishesExtractedTelemetry tests whether a DummyExtractedTelemetry is correctly published to the telemetry collection service.
func (suite *PublisherTestSuite) Test_PublishesExtractedTelemetry() {
	publisher, _ := NewTelemetryPublisher().(*TelemetryPublisher)
	extractedTelemetry := dummyExtractedTelemetry()

	publisher.httpClient.Transport = helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		expectedRequestBody, _ := json.Marshal(extractedTelemetry)
		suite.Equal(req.URL.String(), fmt.Sprintf("%s/api/collect/hosts", apiHost))

		outgoingRequestBody, _ := ioutil.ReadAll(req.Body)
		suite.EqualValues(expectedRequestBody, outgoingRequestBody)

		return &http.Response{
			StatusCode: 202,
		}
	})

	publisher.Publish("dummy_telemetry", uuid.New(), extractedTelemetry)
}

// Test_PublishesExtractedHostTelemetry tests whether extracted HostTelemetries is correctly published to the telemetry collection service.
func (suite *PublisherTestSuite) Test_PublishesExtractedHostTelemetry() {
	publisher, _ := NewTelemetryPublisher().(*TelemetryPublisher)
	extractedHostTelemetry := dummyExtractedHostTelemetry()

	publisher.httpClient.Transport = helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		expectedRequestBody, _ := json.Marshal(extractedHostTelemetry)
		suite.Equal(req.URL.String(), fmt.Sprintf("%s/api/collect/hosts", apiHost))

		outgoingRequestBody, _ := ioutil.ReadAll(req.Body)
		suite.EqualValues(expectedRequestBody, outgoingRequestBody)

		return &http.Response{
			StatusCode: 202,
		}
	})

	publisher.Publish("host_telemetry", uuid.New(), extractedHostTelemetry)
}

// Test_PublishingFailsOnMarshalingError tests whether an error is returned when marshaling the telemetry to JSON fails.
func (suite *PublisherTestSuite) Test_PublishingFailsOnMarshalingError() {
	publisher, _ := NewTelemetryPublisher().(*TelemetryPublisher)
	unmarshable := make(chan int)

	err := publisher.Publish("dummy_telemetry", uuid.New(), unmarshable)

	suite.Error(err)
	suite.Contains(err.Error(), "Failed to marshal telemetry dummy_telemetry")
}

// Test_PublishingFailsOnError tests whether an error is returned when publishing the telemetry fails at net/http level.
func (suite *PublisherTestSuite) Test_PublishingFailsOnError() {
	publisher, _ := NewTelemetryPublisher().(*TelemetryPublisher)
	publisher.httpClient.Transport = helpers.ErroringRoundTripFunc(func() error {
		return fmt.Errorf("some error")
	})

	err := publisher.Publish("dummy_telemetry", uuid.New(), "some")

	suite.Error(err)
	suite.Contains(err.Error(), "An error occurred while publishing telemetry dummy_telemetry")
}

// Test_PublishingFailsOnUnexpectedResponse tests whether an error is returned when the telemetry collection service responds with an unexpected status code.
func (suite *PublisherTestSuite) Test_PublishingFailsOnUnexpectedResponse() {
	publisher, _ := NewTelemetryPublisher().(*TelemetryPublisher)

	publisher.httpClient.Transport = helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
		}
	})

	err := publisher.Publish("dummy_telemetry", uuid.New(), "some")

	suite.Error(err)
	suite.Contains(err.Error(), "Unexpected response code 200 while publishing telemetry dummy_telemetry")
}

type DummyExtractedtelemetry struct {
	Dummy1 string `json:"dummy_1"`
	Dummy2 string `json:"dummy_2"`
}

func dummyExtractedTelemetry() DummyExtractedtelemetry {
	return DummyExtractedtelemetry{
		Dummy1: "dummy_1",
		Dummy2: "dummy_2",
	}
}

func dummyExtractedHostTelemetry() HostTelemetries {
	installationId := uuid.NewString()
	agentId := uuid.NewString()
	anotherAgentId := uuid.NewString()

	return HostTelemetries{
		&HostTelemetry{
			InstallationID: installationId,
			AgentID:        agentId,
			SLESVersion:    "15-sp2",
			CPUCount:       2,
			SocketCount:    8,
			TotalMemoryMB:  4096,
			CloudProvider:  "azure",
			Time:           time.Now(),
		},
		&HostTelemetry{
			InstallationID: installationId,
			AgentID:        anotherAgentId,
			SLESVersion:    "15-sp2",
			CPUCount:       2,
			SocketCount:    8,
			TotalMemoryMB:  4096,
			CloudProvider:  "azure",
			Time:           time.Now(),
		},
	}
}
