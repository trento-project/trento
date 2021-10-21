package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	webApi "github.com/trento-project/trento/web"
)

func (t *trentoApiService) GetSelectedChecksById(clusterId string) (*webApi.JSONSelectedChecks, error) {
	body, statusCode, err := getJson(t.composeQuery(fmt.Sprintf("checks/%s/selected", clusterId)))
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("error during the request with status code %d", statusCode)
	}

	var selectedChecks webApi.JSONSelectedChecks

	err = json.Unmarshal(body, &selectedChecks)
	if err != nil {
		return nil, err
	}

	return &selectedChecks, nil
}
