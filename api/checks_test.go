package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/trento-project/trento/api/mocks"
)

func TestGetSelectedChecksById(t *testing.T) {
	mockGetJson := new(mocks.GetJson)

	returnJson := []byte(`{"selected_checks":["ABCDEF","123456"]}`)

	mockGetJson.On("Execute", "http://192.168.1.10:8000/api/checks/group1/selected").Return(
		returnJson, 200, nil,
	)

	mockGetJson.On("Execute", "http://192.168.1.10:8000/api/checks/otherId/selected").Return(
		nil, 404, fmt.Errorf("not found"),
	)

	getJson = mockGetJson.Execute

	trentoApi := NewTrentoApiService("http://192.168.1.10:8000")
	selectedChecks, err := trentoApi.GetSelectedChecksById("group1")

	assert.NoError(t, err)
	assert.Equal(t, []string{"ABCDEF", "123456"}, selectedChecks.SelectedChecks)

	selectedChecks, err = trentoApi.GetSelectedChecksById("otherId")

	assert.EqualError(t, err, "not found")
}
