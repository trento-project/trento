package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trento-project/trento/test/helpers"
	"github.com/trento-project/trento/web/models"
)

const clustersSettingsResponseMock = `
[
	{
		"id": "c115e7ecf5901d9a768d06a28ab0ba26",
		"selected_checks": [],
		"hosts": [
			{
				"name": "vmnetweaver01",
				"address": "10.1.2.10",
				"user": "root"
			},
			{
				"name": "vmnetweaver02",
				"address": "10.1.2.11",
				"user": "root"
			}
		]
	},
	{
		"id": "5dfbd28f35cbfb38969f9b99243ae8d4",
		"selected_checks": [
			"check1",
			"check2"
		],
		"hosts": [
			{
				"name": "vmhana01",
				"address": "10.1.2.12",
				"user": "cloudadmin"
			},
			{
				"name": "vmhana02",
				"address": "10.1.2.13",
				"user": "cloudadmin"
			}
		]
	}
]`

type ClusterSettingsApiTestCase struct {
	suite.Suite
	trentoApi *trentoApiService
}

func TestClusterSettingsApiTestCase(t *testing.T) {
	suite.Run(t, new(ClusterSettingsApiTestCase))
}

func (suite *ClusterSettingsApiTestCase) SetupSuite() {
	suite.trentoApi = NewTrentoApiService("192.168.1.10", 8000)
}

func (suite *ClusterSettingsApiTestCase) Test_AnErrorOccursInCommunication() {
	suite.trentoApi.httpClient.Transport = helpers.ErroringRoundTripFunc(func() error {
		return fmt.Errorf("some error")
	})

	_, err := suite.trentoApi.GetClustersSettings()

	suite.Error(err, "some error")
}

func (suite *ClusterSettingsApiTestCase) Test_ResponseIsNotSuccessful() {
	suite.trentoApi.httpClient.Transport = helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 401,
		}
	})

	_, err := suite.trentoApi.GetClustersSettings()

	suite.Error(err, "error during the request with status code 401")
}

func (suite *ClusterSettingsApiTestCase) Test_ResponseUnmarshalingFailure() {
	suite.trentoApi.httpClient.Transport = helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("unmatching")),
		}
	})

	_, err := suite.trentoApi.GetClustersSettings()

	suite.Contains(err.Error(), "invalid character")
}

func (suite *ClusterSettingsApiTestCase) Test_ClustersSettingsAreSuccessfullyRetrieved() {
	suite.trentoApi.httpClient.Transport = helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(clustersSettingsResponseMock)),
		}
	})

	clustersSettings, _ := suite.trentoApi.GetClustersSettings()
	suite.Len(clustersSettings, 2)

	cluster0 := clustersSettings[0]
	suite.Equal("c115e7ecf5901d9a768d06a28ab0ba26", cluster0.ID)
	suite.Empty(cluster0.SelectedChecks)
	assertMatchingHostsOnCluster0(suite, cluster0.Hosts)

	cluster1 := clustersSettings[1]
	suite.Equal("5dfbd28f35cbfb38969f9b99243ae8d4", cluster1.ID)
	suite.NotEmpty(cluster1.SelectedChecks)
	suite.Equal([]string{"check1", "check2"}, cluster1.SelectedChecks)
	assertMatchingHostsOnCluster1(suite, cluster1.Hosts)

}

func assertMatchingHostsOnCluster0(suite *ClusterSettingsApiTestCase, hosts []*models.HostConnection) {
	suite.Len(hosts, 2)

	suite.Equal("vmnetweaver01", hosts[0].Name)
	suite.Equal("10.1.2.10", hosts[0].Address)
	suite.Equal("root", hosts[0].User)

	suite.Equal("vmnetweaver02", hosts[1].Name)
	suite.Equal("10.1.2.11", hosts[1].Address)
	suite.Equal("root", hosts[1].User)
}

func assertMatchingHostsOnCluster1(suite *ClusterSettingsApiTestCase, hosts []*models.HostConnection) {
	suite.Len(hosts, 2)

	suite.Equal("vmhana01", hosts[0].Name)
	suite.Equal("10.1.2.12", hosts[0].Address)
	suite.Equal("cloudadmin", hosts[0].User)

	suite.Equal("vmhana02", hosts[1].Name)
	suite.Equal("10.1.2.13", hosts[1].Address)
	suite.Equal("cloudadmin", hosts[1].User)
}
