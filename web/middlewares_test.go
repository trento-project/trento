package web

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/web/services"
)

func TestEulaMiddleware(t *testing.T) {
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
	req := httptest.NewRequest("GET", "/", nil)

	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 302, resp.Code)
	assert.Contains(t, resp.Body.String(), "License agreement")
}

func TestEulaMiddlewareNotPremium(t *testing.T) {
	mockedSettingsService := new(services.MockSettingsService)
	mockedSettingsService.On("AcceptEula").Return(nil)
	mockedSettingsService.On("IsEulaAccepted").Return(false, nil)
	mockedSettingsService.On("InitializeIdentifier").Return(uuid.MustParse("59fd8017-b7fd-477b-9ebe-b658c558f3e9"), nil)

	mockedSubscriptionsService := new(services.MockSubscriptionsService)
	mockedSubscriptionsService.On("IsTrentoPremium").Return(false, nil)

	deps := setupTestDependencies()

	deps.settingsService = mockedSettingsService
	deps.subscriptionsService = mockedSubscriptionsService

	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "Blue Horizon for SAP")
}

func TestEulaMiddlewareLettingThrough(t *testing.T) {
	mockedSettingsService := new(services.MockSettingsService)
	mockedSettingsService.On("AcceptEula").Return(nil)
	mockedSettingsService.On("IsEulaAccepted").Return(true, nil)
	mockedSettingsService.On("InitializeIdentifier").Return(uuid.MustParse("59fd8017-b7fd-477b-9ebe-b658c558f3e9"), nil)
	deps := setupTestDependencies()
	deps.settingsService = mockedSettingsService
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "Blue Horizon for SAP")
}

func TestEulaMiddlewareError(t *testing.T) {
	mockedSettingsService := new(services.MockSettingsService)
	mockedSettingsService.On("AcceptEula").Return(nil)
	mockedSettingsService.On("IsEulaAccepted").Return(false, errors.New("EULA error"))
	mockedSettingsService.On("InitializeIdentifier").Return(uuid.MustParse("59fd8017-b7fd-477b-9ebe-b658c558f3e9"), nil)
	deps := setupTestDependencies()
	deps.settingsService = mockedSettingsService
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)
}
