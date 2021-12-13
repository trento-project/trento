package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	webApi "github.com/trento-project/trento/web"
)

func (t *trentoApiService) GetClustersSettings() (webApi.ClustersSettingsResponse, error) {
	body, statusCode, err := t.getJson("internal/clusters/settings")
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("error during the request with status code %d", statusCode)
	}

	var clustersSettings webApi.ClustersSettingsResponse

	err = json.Unmarshal(body, &clustersSettings)
	if err != nil {
		return nil, err
	}

	return clustersSettings, nil
}
