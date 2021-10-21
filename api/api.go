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
	GetSelectedChecksById(clusterId string) (*webApi.JSONSelectedChecks, error)
}

type trentoApiService struct {
	webServer string
}

func NewTrentoApiService(webServer string) *trentoApiService {
	return &trentoApiService{webServer: webServer}
}

func (t *trentoApiService) composeQuery(resource string) string {
	return fmt.Sprintf("%s/api/%s", t.webServer, resource)
}

//go:generate mockery --name=GetJson

type GetJson func(query string) ([]byte, int, error)

var getJson GetJson = func(query string) ([]byte, int, error) {
	var err error
	resp, err := http.Get(query)
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

//go:generate mockery --name=GetHttp

type GetHttp func(url string) (*http.Response, error)

var getHttp GetHttp = func(url string) (*http.Response, error) {
	return http.Get(url)
}

func (t *trentoApiService) IsWebServerUp() bool {
	host := t.composeQuery("ping")
	log.Debugf("Looking for the Trento server state at: %s", host)

	resp, err := getHttp(host)

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
