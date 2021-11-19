package web

import (
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"

	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

func TestAboutHandlerPremium(t *testing.T) {
	subscriptionsMocks := new(services.MockSubscriptionsService)
	premiumData := &models.PremiumData{
		IsPremium:     true,
		Sles4SapCount: 2,
	}
	subscriptionsMocks.On("GetPremiumData").Return(premiumData, nil)

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

func TestAboutHandlerCommunity(t *testing.T) {
	subscriptionsMocks := new(services.MockSubscriptionsService)
	premiumData := &models.PremiumData{
		IsPremium:     false,
		Sles4SapCount: 0,
	}
	subscriptionsMocks.On("GetPremiumData").Return(premiumData, nil)

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
	assert.Regexp(t, regexp.MustCompile("<dt.*>Subscription</dt><dd.*>Community.*</dd>"), minified)
	assert.NotContains(t, minified, "SLES_SAP machine")
}
