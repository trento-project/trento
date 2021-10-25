package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/trento-project/trento/web/datapipeline"
	"github.com/trento-project/trento/web/services"
)

func TestApiCollectDataHandler(t *testing.T) {
	collectorService := new(services.MockCollectorService)
	collectorService.On("StoreEvent", mock.Anything).Return(nil)

	deps := setupTestDependencies()
	deps.collectorService = collectorService

	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	body, _ := json.Marshal(&datapipeline.DataCollectedEvent{
		AgentID:       "agent_id",
		DiscoveryType: "discovery",
		Payload:       []byte("{}"),
	})
	req, err := http.NewRequest("POST", "/api/collect", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	app.collectorEngine.ServeHTTP(resp, req)

	assert.Equal(t, 202, resp.Code)
}
