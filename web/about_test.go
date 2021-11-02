package web

import (
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"

	"github.com/trento-project/trento/web/services"
)

func TestAboutHandlerPremium(t *testing.T) {
	subscriptionsMocks := new(services.MockSubscriptionsService)

	subscriptionsMocks.On("GetSubscriptionData").Return(
		&services.SubscriptionData{Type: services.Premium, SubscribedCount: 2}, nil)

	deps := setupTestDependencies()
	deps.subscriptionsService = subscriptionsMocks

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/about", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.webEngine.ServeHTTP(resp, req)

	subscriptionsMocks.AssertExpectations(t)

	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepEndTags:         true,
	})
	minified, err := m.String("text/html", resp.Body.String())
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, minified, "About")
	assert.Regexp(t, regexp.MustCompile("<dt.*>Subscription</dt><dd.*>Premium.*</dd>"), minified)
	assert.Regexp(t, regexp.MustCompile("<dt.*>SLES_SAP machines</dt><dd.*>2.*</dd>"), minified)
}

func TestAboutHandlerFree(t *testing.T) {
	subscriptionsMocks := new(services.MockSubscriptionsService)

	subscriptionsMocks.On("GetSubscriptionData").Return(
		&services.SubscriptionData{Type: services.Free, SubscribedCount: 0}, nil)

	deps := setupTestDependencies()
	deps.subscriptionsService = subscriptionsMocks

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/about", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.webEngine.ServeHTTP(resp, req)

	subscriptionsMocks.AssertExpectations(t)

	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepEndTags:         true,
	})
	minified, err := m.String("text/html", resp.Body.String())
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, minified, "About")
	assert.Regexp(t, regexp.MustCompile("<dt.*>Subscription</dt><dd.*>Free.*</dd>"), minified)
	assert.NotContains(t, minified, "SLES_SAP machine")
}
