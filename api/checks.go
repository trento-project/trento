package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	webApi "github.com/trento-project/trento/web"
	"github.com/trento-project/trento/web/models"
)

func (t *trentoApiService) GetSelectedChecksById(clusterId string) (*webApi.JSONSelectedChecks, error) {
	body, statusCode, err := t.getJson(fmt.Sprintf("checks/%s/selected", clusterId))
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

func (t *trentoApiService) GetConnectionDataById(id string) (map[string]*models.ConnectionData, error) {
	body, statusCode, err := t.getJson(fmt.Sprintf("checks/%s/connection_data", id))
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("error during the request with status code %d", statusCode)
	}

	connData := make(map[string]*models.ConnectionData)

	err = json.Unmarshal(body, &connData)
	if err != nil {
		return nil, err
	}

	return connData, nil
}
