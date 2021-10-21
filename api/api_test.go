package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/trento-project/trento/api/mocks"
)

func TestIsWebServerUp(t *testing.T) {
	mockGetHttp := new(mocks.GetHttp)

	found := &http.Response{StatusCode: http.StatusOK}
	mockGetHttp.On("Execute", "http://192.168.1.10:8000/api/ping").Return(
		found, nil,
	).Times(1)

	notFound := &http.Response{StatusCode: http.StatusBadRequest}
	mockGetHttp.On("Execute", "http://192.168.1.10:8000/api/ping").Return(
		notFound, nil,
	)

	getHttp = mockGetHttp.Execute

	trentoApi := NewTrentoApiService("http://192.168.1.10:8000")
	result := trentoApi.IsWebServerUp()

	assert.Equal(t, true, result)

	result = trentoApi.IsWebServerUp()

	assert.Equal(t, false, result)
}
