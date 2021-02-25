package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomeRoute(t *testing.T) {
	engine := makeEngine()

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	engine.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
}
