package collector

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	_ "github.com/trento-project/trento/test"
	"github.com/trento-project/trento/test/helpers"
)

const (
	DummyMachineID = "dummy-machine-id"
	DummyAgentID   = "779cdd70-e9e2-58ca-b18a-bf3eb3f71244"
)

type CollectorClientTestSuite struct {
	suite.Suite
}

func TestCollectorClientTestSuite(t *testing.T) {
	suite.Run(t, new(CollectorClientTestSuite))
}

func (suite *CollectorClientTestSuite) SetupSuite() {
	fileSystem = afero.NewMemMapFs()

	afero.WriteFile(fileSystem, machineIdPath, []byte(DummyMachineID), 0644)

	// this is read by an env variable called TRENTO_COLLECTOR_ENABLED
	viper.Set("collector-enabled", true)
}

func (suite *CollectorClientTestSuite) TestCollectorClient_NewClientWithTLS() {
	collectorClient, err := NewCollectorClient(&Config{
		EnablemTLS:    true,
		CollectorHost: "localhost",
		CollectorPort: 8081,
		Cert:          "./test/certs/client-cert.pem",
		Key:           "./test/certs/client-key.pem",
		CA:            "./test/certs/ca-cert.pem",
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
		Cert:          "./test/certs/client-cert.pem",
		Key:           "./test/certs/client-key.pem",
		CA:            "./test/certs/ca-cert.pem",
	})

	suite.NoError(err)

	discoveredDataPayload := struct {
		FieldA string
	}{
		FieldA: "some discovered field",
	}

	discoveryType := "the_discovery_type"

	collectorClient.httpClient.Transport = helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		requestBody, _ := json.Marshal(map[string]interface{}{
			"agent_id":       DummyAgentID,
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
