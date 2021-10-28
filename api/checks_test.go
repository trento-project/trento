package api

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/trento-project/trento/test/helpers"

	"github.com/trento-project/trento/web"
)

func TestGetChecksSettingsById(t *testing.T) {
	trentoApi := NewTrentoApiService("http://192.168.1.10:8000")

	trentoApi.httpClient = &http.Client{Transport: helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), "http://192.168.1.10:8000/api/checks/group1/settings")

		return &http.Response{
			StatusCode: 200,
			Body: io.NopCloser(strings.NewReader(`
				{
					"selected_checks": ["ABCDEF","123456"],
					"connection_settings": {
						"node1":"user1",
						"node2":"user2"
					}
				}`)),
		}
	})}

	connData, err := trentoApi.GetChecksSettingsById("group1")

	expectedData := &web.JSONChecksSettings{
		SelectedChecks: []string{"ABCDEF", "123456"},
		ConnectionSettings: map[string]string{
			"node1": "user1",
			"node2": "user2",
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedData, connData)
}

func TestGetChecksSettingsByIdNotFound(t *testing.T) {
	trentoApi := NewTrentoApiService("http://192.168.1.10:8000")

	trentoApi.httpClient = &http.Client{Transport: helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), "http://192.168.1.10:8000/api/checks/otherId/settings")

		return &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(strings.NewReader("not found")),
		}
	})}

	_, err := trentoApi.GetChecksSettingsById("otherId")

	assert.EqualError(t, err, "error during the request with status code 404")
}
