package web

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

func TestGetHttpSDTargets(t *testing.T) {
	targets := models.PrometheusTargetsList{
		&models.PrometheusTargets{
			Targets: []string{"192.168.1.1:9100"},
			Labels:  map[string]string{"hostname": "host1"},
		},
		&models.PrometheusTargets{
			Targets: []string{"192.168.1.2:9100"},
			Labels:  map[string]string{"hostname": "host2"},
		},
		&models.PrometheusTargets{
			Targets: []string{"192.168.1.3:9100"},
			Labels:  map[string]string{"hostname": "host3"},
		},
	}
	mockPrometheusService := new(services.MockPrometheusService)
	mockPrometheusService.On("GetHttpSDTargets").Return(targets, nil)

	deps := setupTestDependencies()
	deps.prometheusService = mockPrometheusService

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/prometheus/targets", nil)
	app.webEngine.ServeHTTP(resp, req)

	expectedBody, _ := json.Marshal([]map[string]interface{}{
		map[string]interface{}{
			"targets": []string{"192.168.1.1:9100"},
			"labels":  map[string]string{"hostname": "host1"}},
		map[string]interface{}{
			"targets": []string{"192.168.1.2:9100"},
			"labels":  map[string]string{"hostname": "host2"}},
		map[string]interface{}{
			"targets": []string{"192.168.1.3:9100"},
			"labels":  map[string]string{"hostname": "host3"}},
	})
	assert.JSONEq(t, string(expectedBody), resp.Body.String())
	assert.Equal(t, 200, resp.Code)
}
