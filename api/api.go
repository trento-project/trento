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
	GetChecksSettingsById(id string) (*webApi.JSONChecksSettings, error)
}

type trentoApiService struct {
	webServer  string
	httpClient *http.Client
}

func NewTrentoApiService(webServer string) *trentoApiService {
	client := &http.Client{}
	return &trentoApiService{webServer: webServer, httpClient: client}
}

func (t *trentoApiService) composeQuery(resource string) string {
	return fmt.Sprintf("%s/api/%s", t.webServer, resource)
}

func (t *trentoApiService) getJson(query string) ([]byte, int, error) {
	var err error

	resp, err := t.httpClient.Get(t.composeQuery(query))
	if err != nil {
		return nil, resp.StatusCode, err
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
