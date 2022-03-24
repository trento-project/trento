package web

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/web/services"
)

func TestEulaMiddleware(t *testing.T) {
	mockedPremiumDetection := new(services.MockPremiumDetectionService)
	mockedPremiumDetection.On("RequiresEulaAcceptance").Return(true, nil)
	deps := setupTestDependencies()
	deps.premiumDetectionService = mockedPremiumDetection
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

func TestEulaMiddlewareLettingThrough(t *testing.T) {
	mockedPremiumDetection := new(services.MockPremiumDetectionService)
	mockedPremiumDetection.On("RequiresEulaAcceptance").Return(false, nil)
	deps := setupTestDependencies()
	deps.premiumDetectionService = mockedPremiumDetection
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Body.String(), "homepage-component")
}

func TestEulaMiddlewareError(t *testing.T) {
	mockedPremiumDetection := new(services.MockPremiumDetectionService)
	mockedPremiumDetection.On("RequiresEulaAcceptance").Return(false, errors.New("EULA error"))
	deps := setupTestDependencies()
	deps.premiumDetectionService = mockedPremiumDetection
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
