package api

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/trento-project/trento/test/helpers"

	"github.com/trento-project/trento/web/models"
)

func TestGetSelectedChecksById(t *testing.T) {
	trentoApi := NewTrentoApiService("http://192.168.1.10:8000")

	trentoApi.httpClient = &http.Client{Transport: helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), "http://192.168.1.10:8000/api/checks/group1/selected")

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"selected_checks":["ABCDEF","123456"]}`)),
		}
	})}

	selectedChecks, err := trentoApi.GetSelectedChecksById("group1")

	assert.NoError(t, err)
	assert.Equal(t, []string{"ABCDEF", "123456"}, selectedChecks.SelectedChecks)

}

func TestGetSelectedChecksByIdNotFound(t *testing.T) {
	trentoApi := NewTrentoApiService("http://192.168.1.10:8000")

	trentoApi.httpClient = &http.Client{Transport: helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), "http://192.168.1.10:8000/api/checks/otherId/selected")

		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(strings.NewReader("not found")),
		}
	})}

	_, err := trentoApi.GetSelectedChecksById("otherId")

	assert.EqualError(t, err, "error during the request with status code 404")
}

func TestGetConnectionDataById(t *testing.T) {
	trentoApi := NewTrentoApiService("http://192.168.1.10:8000")

	trentoApi.httpClient = &http.Client{Transport: helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), "http://192.168.1.10:8000/api/checks/group1/connection_data")

		return &http.Response{
			StatusCode: 200,
			Body: io.NopCloser(strings.NewReader(`
				{"node1":{"id":"group1","node":"node1","user":"user1"},
				"node2":{"id":"group1","node":"node2","user":"user2"}}`)),
		}
	})}

	connData, err := trentoApi.GetConnectionDataById("group1")

	expectedData := map[string]*models.ConnectionData{
		"node1": &models.ConnectionData{ID: "group1", Node: "node1", User: "user1"},
		"node2": &models.ConnectionData{ID: "group1", Node: "node2", User: "user2"},
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedData, connData)

}

func TestGetConnectionDataByIdNotFound(t *testing.T) {
	trentoApi := NewTrentoApiService("http://192.168.1.10:8000")

	trentoApi.httpClient = &http.Client{Transport: helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), "http://192.168.1.10:8000/api/checks/otherId/connection_data")

		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(strings.NewReader("not found")),
		}
	})}

	_, err := trentoApi.GetConnectionDataById("otherId")

	assert.EqualError(t, err, "error during the request with status code 404")
}
