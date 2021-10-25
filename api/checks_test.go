package api

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/trento-project/trento/test/helpers"
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
