package collector

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/agent/discovery/mocks"
	"github.com/trento-project/trento/test/helpers"
)

type CollectorClientTestSuite struct {
	suite.Suite
	configuredClient *collectorClient
}

func TestCollectorClientTestSuite(t *testing.T) {
	suite.Run(t, new(CollectorClientTestSuite))
}

func (suite *CollectorClientTestSuite) SetupSuite() {
	fileSystem = afero.NewMemMapFs()

	afero.WriteFile(fileSystem, machineIdPath, []byte("the-machine-id"), 0644)

	// this is read by an env variable called TRENTO_COLLECTOR_ENABLED
	viper.Set("collector-enabled", true)

	collectorClient, err := NewCollectorClient(&Config{
		EnablemTLS:    true,
		CollectorHost: "localhost",
		CollectorPort: 8443,
		Cert:          "../../test/certs/client-cert.pem",
		Key:           "../../test/certs/client-key.pem",
		CA:            "../../test/certs/ca-cert.pem",
	})
	suite.NoError(err)

	suite.configuredClient = collectorClient
}

func (suite *CollectorClientTestSuite) TestCollectorClient_NewClientWithTLS() {
	collectorClient, err := NewCollectorClient(&Config{
		EnablemTLS:    true,
		CollectorHost: "localhost",
		CollectorPort: 8081,
		Cert:          "../../test/certs/client-cert.pem",
		Key:           "../../test/certs/client-key.pem",
		CA:            "../../test/certs/ca-cert.pem",
	})

	suite.NoError(err)

	transport, _ := (collectorClient.httpClient.Transport).(*http.Transport)

	suite.Equal(1, len(transport.TLSClientConfig.Certificates))
}

func (suite *CollectorClientTestSuite) TestCollectorClient_NewClientWithoutTLS() {
	collectorClient, err := NewCollectorClient(&Config{
		EnablemTLS:    false,
		CollectorHost: "localhost",
		CollectorPort: 8081,
		Cert:          "",
		Key:           "",
		CA:            "",
	})

	suite.NoError(err)

	transport, _ := (collectorClient.httpClient.Transport).(*http.Transport)

	suite.Equal((*tls.Config)(nil), transport.TLSClientConfig)
}

func (suite *CollectorClientTestSuite) TestCollectorClient_PublishingSuccess() {
	collectorClient, err := NewCollectorClient(&Config{
		EnablemTLS:    true,
		CollectorHost: "localhost",
		CollectorPort: 8081,
		Cert:          "../../test/certs/client-cert.pem",
		Key:           "../../test/certs/client-key.pem",
		CA:            "../../test/certs/ca-cert.pem",
	})

	suite.NoError(err)

	discoveredDataPayload := struct {
		FieldA string
	}{
		FieldA: "some discovered field",
	}

	discoveryType := "the_discovery_type"
	agentID := "the-machine-id"

	collectorClient.httpClient.Transport = helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		requestBody, _ := json.Marshal(map[string]interface{}{
			"agent_id":       agentID,
			"discovery_type": discoveryType,
			"payload":        discoveredDataPayload,
		})

		bodyBytes, _ := ioutil.ReadAll(req.Body)

		suite.EqualValues(requestBody, bodyBytes)

		suite.Equal(req.URL.String(), "https://localhost:8081/api/collect")
		return &http.Response{
			StatusCode: 202,
		}
	})

	err = collectorClient.Publish(discoveryType, discoveredDataPayload)

	suite.NoError(err)
}

func (suite *CollectorClientTestSuite) TestCollectorClient_PublishingFailure() {
	collectorClient, err := NewCollectorClient(&Config{
		EnablemTLS:    false,
		CollectorHost: "localhost",
		CollectorPort: 8081,
	})

	suite.NoError(err)

	collectorClient.httpClient.Transport = helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		suite.Equal(req.URL.String(), "http://localhost:8081/api/collect")
		return &http.Response{
			StatusCode: 500,
		}
	})

	err = collectorClient.Publish("some_discovery_type", struct{}{})

	suite.Error(err)
}

// Following test cover publishing data from the discovery loops

func (suite *CollectorClientTestSuite) TestCollectorClient_PublishingClusterDiscovery() {
	discoveryType := "ha_cluster_discovery"
	discoveredCluster := mocks.NewDiscoveredClusterMock()

	suite.runDiscoveryScenario(discoveryType, discoveredCluster, func(requestBodyAgainstCollector string) {
		suite.assertJsonMatchesJsonFileContent("../../test/fixtures/discovery/cluster/expected_published_cluster_discovery.json", requestBodyAgainstCollector)
	})
}

func (suite *CollectorClientTestSuite) TestCollectorClient_PublishingCloudDiscovery() {
	discoveryType := "cloud_discovery"
	discoveredCloudInstance := mocks.NewDiscoveredCloudMock()

	suite.runDiscoveryScenario(discoveryType, discoveredCloudInstance, func(requestBodyAgainstCollector string) {
		suite.assertJsonMatchesJsonFileContent("../../test/fixtures/discovery/azure/expected_published_cloud_discovery.json", requestBodyAgainstCollector)
	})
}

func (suite *CollectorClientTestSuite) TestCollectorClient_PublishingHostDiscovery() {
	discoveryType := "host_discovery"
	discoveredHost := mocks.NewDiscoveredHostMock()

	suite.runDiscoveryScenario(discoveryType, discoveredHost, func(requestBodyAgainstCollector string) {
		suite.assertJsonMatchesJsonFileContent("../../test/fixtures/discovery/host/expected_published_host_discovery.json", requestBodyAgainstCollector)
	})
}

func (suite *CollectorClientTestSuite) TestCollectorClient_PublishingSubscriptionDiscovery() {
	discoveryType := "subscription_discovery"
	discoveredSubscriptions := mocks.NewDiscoveredSubscriptionsMock()

	suite.runDiscoveryScenario(discoveryType, discoveredSubscriptions, func(requestBodyAgainstCollector string) {
		suite.assertJsonMatchesJsonFileContent("../../test/fixtures/discovery/subscriptions/expected_published_subscriptions_discovery.json", requestBodyAgainstCollector)
	})
}

type AssertionFunc func(requestBodyAgainstCollector string)

func (suite *CollectorClientTestSuite) runDiscoveryScenario(discoveryType string, payload interface{}, assertion AssertionFunc) {
	agentID := "the-machine-id"

	collectorClient := suite.configuredClient

	collectorClient.httpClient.Transport = helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		requestBody, _ := json.Marshal(map[string]interface{}{
			"agent_id":       agentID,
			"discovery_type": discoveryType,
			"payload":        payload,
		})

		outgoingRequestBody, _ := ioutil.ReadAll(req.Body)

		suite.EqualValues(requestBody, outgoingRequestBody)

		assertion(string(outgoingRequestBody))

		suite.Equal(req.URL.String(), "https://localhost:8443/api/collect")
		return &http.Response{
			StatusCode: 202,
		}
	})

	err := collectorClient.Publish(discoveryType, payload)

	suite.NoError(err)
}

func (suite *CollectorClientTestSuite) assertJsonMatchesJsonFileContent(expectedJsonContentPath string, actualJson string) {
	expectedJsonContent, err := os.Open(expectedJsonContentPath)
	if err != nil {
		panic(err)
	}

	var expectedJsonContentBytesBuffer *bytes.Buffer = new(bytes.Buffer)

	expectedJsonContentByte, _ := ioutil.ReadAll(expectedJsonContent)

	json.Compact(expectedJsonContentBytesBuffer, expectedJsonContentByte)

	b, _ := ioutil.ReadAll(expectedJsonContentBytesBuffer)

	suite.EqualValues(actualJson, b)
}
