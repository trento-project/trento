package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	webApi "github.com/trento-project/trento/web"
)

func (t *trentoApiService) GetChecksSettingsById(id string) (*webApi.JSONChecksSettings, error) {
	body, statusCode, err := t.getJson(fmt.Sprintf("checks/%s/settings", id))
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("error during the request with status code %d", statusCode)
	}

	var checksSettings *webApi.JSONChecksSettings

	err = json.Unmarshal(body, &checksSettings)
	if err != nil {
		return nil, err
	}

	return checksSettings, nil
}
