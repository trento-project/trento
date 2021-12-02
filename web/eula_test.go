package web

import (
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/web/services"
)

func TestEulaHandler(t *testing.T) {
	deps := setupTestDependencies()
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/eula", nil)

	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "License agreement")
}

func TestEulaAcceptHandler(t *testing.T) {
	mockedSettingsService := new(services.MockSettingsService)
	mockedSettingsService.On("AcceptEula").Return(nil)
	mockedSettingsService.On("IsEulaAccepted").Return(false, nil)
	mockedSettingsService.On("InitializeIdentifier").Return(uuid.MustParse("59fd8017-b7fd-477b-9ebe-b658c558f3e9"), nil)
	deps := setupTestDependencies()
	deps.settingsService = mockedSettingsService
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/accept-eula", nil)

	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 302, resp.Code)
}
