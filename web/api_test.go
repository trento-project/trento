package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiPingTest(t *testing.T) {
	app, err := NewAppWithDeps(setupTestConfig(), setupTestDependencies())
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/ping", nil)
	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "pong", resp.Body.String())
}

func TestApiDocsRouteTest(t *testing.T) {
	app, err := NewAppWithDeps(setupTestConfig(), setupTestDependencies())
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/docs/index.html", nil)
	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	resp = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/docs/doc.json", nil)
	app.webEngine.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
}
