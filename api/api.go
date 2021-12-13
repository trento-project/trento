package api

import (
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	webApi "github.com/trento-project/trento/web"
)

//go:generate mockery --all

type TrentoApiService interface {
	IsWebServerUp() bool
	GetClustersSettings() (webApi.ClustersSettingsResponse, error)
}

type trentoApiService struct {
	apiHost    string
	apiPort    int
	httpClient *http.Client
}

func NewTrentoApiService(apiHost string, apiPort int) *trentoApiService {
	client := &http.Client{}
	return &trentoApiService{apiHost: apiHost, apiPort: apiPort, httpClient: client}
}

func (t *trentoApiService) composeQuery(resource string) string {
	return fmt.Sprintf("http://%s:%d/api/%s", t.apiHost, t.apiPort, resource)
}

func (t *trentoApiService) getJson(query string) ([]byte, int, error) {
	var err error

	resp, err := t.httpClient.Get(t.composeQuery(query))
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}

func (t *trentoApiService) IsWebServerUp() bool {
	host := t.composeQuery("ping")
	log.Debugf("Looking for the Trento server state at: %s", host)

	resp, err := t.httpClient.Get(host)

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Debugf("Error requesting Trento server api: %s", err)
		return false
	}

	log.Debugf("Trento server response code: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}
