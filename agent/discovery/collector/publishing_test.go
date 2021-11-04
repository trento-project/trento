package collector

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/agent/discovery/mocks"
	_ "github.com/trento-project/trento/test"
	"github.com/trento-project/trento/test/helpers"
)

type PublishingTestSuite struct {
	suite.Suite
	configuredClient *client
}

func TestPublishingTestSuite(t *testing.T) {
	suite.Run(t, new(PublishingTestSuite))
}

func (suite *PublishingTestSuite) SetupSuite() {
	fileSystem = afero.NewMemMapFs()

	afero.WriteFile(fileSystem, machineIdPath, []byte("the-machine-id"), 0644)

	// this is read by an env variable called TRENTO_COLLECTOR_ENABLED
	viper.Set("collector-enabled", true)

	collectorClient, err := NewCollectorClient(&Config{
		EnablemTLS:    true,
		CollectorHost: "localhost",
		CollectorPort: 8443,
		Cert:          "./test/certs/client-cert.pem",
		Key:           "./test/certs/client-key.pem",
		CA:            "./test/certs/ca-cert.pem",
	})
	suite.NoError(err)

	suite.configuredClient = collectorClient
}

// Following test cover publishing data from the discovery loops

func (suite *PublishingTestSuite) TestCollectorClient_PublishingClusterDiscovery() {
	discoveryType := "ha_cluster_discovery"
	discoveredCluster := mocks.NewDiscoveredClusterMock()

	suite.runDiscoveryScenario(discoveryType, discoveredCluster, func(requestBodyAgainstCollector string) {
		suite.assertJsonMatchesJsonFileContent("./test/fixtures/discovery/cluster/expected_published_cluster_discovery.json", requestBodyAgainstCollector)
	})
}

func (suite *PublishingTestSuite) TestCollectorClient_PublishingCloudDiscovery() {
	discoveryType := "cloud_discovery"
	discoveredCloudInstance := mocks.NewDiscoveredCloudMock()

	suite.runDiscoveryScenario(discoveryType, discoveredCloudInstance, func(requestBodyAgainstCollector string) {
		suite.assertJsonMatchesJsonFileContent("./test/fixtures/discovery/azure/expected_published_cloud_discovery.json", requestBodyAgainstCollector)
	})
}

func (suite *PublishingTestSuite) TestCollectorClient_PublishingHostDiscovery() {
	discoveryType := "host_discovery"
	discoveredHost := mocks.NewDiscoveredHostMock()

	suite.runDiscoveryScenario(discoveryType, discoveredHost, func(requestBodyAgainstCollector string) {
		suite.assertJsonMatchesJsonFileContent("./test/fixtures/discovery/host/expected_published_host_discovery.json", requestBodyAgainstCollector)
	})
}

func (suite *PublishingTestSuite) TestCollectorClient_PublishingSubscriptionDiscovery() {
	discoveryType := "subscription_discovery"
	discoveredSubscriptions := mocks.NewDiscoveredSubscriptionsMock()

	suite.runDiscoveryScenario(discoveryType, discoveredSubscriptions, func(requestBodyAgainstCollector string) {
		suite.assertJsonMatchesJsonFileContent("./test/fixtures/discovery/subscriptions/expected_published_subscriptions_discovery.json", requestBodyAgainstCollector)
	})
}

type AssertionFunc func(requestBodyAgainstCollector string)

func (suite *PublishingTestSuite) runDiscoveryScenario(discoveryType string, payload interface{}, assertion AssertionFunc) {
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

func (suite *PublishingTestSuite) assertJsonMatchesJsonFileContent(expectedJsonContentPath string, actualJson string) {
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
