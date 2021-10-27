package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/trento-project/trento/test/helpers"
)

func TestIsWebServerUp(t *testing.T) {
	trentoApi := NewTrentoApiService("192.168.1.10", 8000)

	trentoApi.httpClient = &http.Client{Transport: helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), "http://192.168.1.10:8000/api/ping")

		return &http.Response{
			StatusCode: 200,
		}
	})}

	result := trentoApi.IsWebServerUp()

	assert.Equal(t, true, result)

	trentoApi.httpClient = &http.Client{Transport: helpers.RoundTripFunc(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), "http://192.168.1.10:8000/api/ping")

		return &http.Response{
			StatusCode: 500,
		}
	})}

	result = trentoApi.IsWebServerUp()

	assert.Equal(t, false, result)
}
