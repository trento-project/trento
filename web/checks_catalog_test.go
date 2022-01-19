package web

import (
	"fmt"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services"
)

func TestChecksCatalogHandler(t *testing.T) {
	checksService := new(services.MockChecksService)

	deps := setupTestDependencies()
	deps.checksService = checksService

	checks := models.GroupedCheckList{
		&models.GroupedChecks{
			Group: "group 1",
			Checks: models.ChecksCatalog{
				&models.Check{
					ID:             "ABCDEF",
					Name:           "1.1.1",
					Group:          "group 1",
					Description:    "description 1",
					Remediation:    "remediation 1",
					Implementation: "implementation 1",
					Labels:         "labels 1",
				},
				&models.Check{
					ID:             "123456",
					Name:           "1.1.2",
					Group:          "group 1",
					Description:    "description 2",
					Remediation:    "remediation 2",
					Implementation: "implementation 2",
					Labels:         "labels 2",
					Premium:        true,
				},
			},
		},
		&models.GroupedChecks{
			Group: "group 2",
			Checks: models.ChecksCatalog{
				&models.Check{
					ID:             "123ABC",
					Name:           "1.2.1",
					Group:          "group 2",
					Description:    "description 3",
					Remediation:    "remediation 3",
					Implementation: "implementation 3",
					Labels:         "labels 3",
				},
			},
		},
	}

	checksService.On("GetChecksCatalogByGroup").Return(
		checks, nil,
	)

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/catalog", nil)

	app.webEngine.ServeHTTP(resp, req)

	responseBody := minifyHtml(resp.Body.String())

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, responseBody, "Checks catalog")
	assert.Contains(t, responseBody, "<div class=check-group id=check-group-0>")
	assert.Regexp(t, regexp.MustCompile("<h4.*?>group 1</h4>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<td.*?>ABCDEF</td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<div class=check-description>.*description 1.*</div>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<div class=\"check-remediation collapse\" id=collapse-ABCDEF.*?>.*remediation 1.*<pre>implementation 1</pre>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<div class=check-description>.*<p>description 2 <span class=\"badge badge-trento-premium\">Premium</span></p>.*</div>"), responseBody)
	assert.Equal(t, 2, strings.Count(responseBody, "<div class=check-group"))
	assert.Equal(t, 3, strings.Count(responseBody, "<tr class=check-row"))

	checksService.AssertExpectations(t)
}

func TestChecksCatalogHandlerError(t *testing.T) {
	checksService := new(services.MockChecksService)

	deps := setupTestDependencies()
	deps.checksService = checksService

	checksService.On("GetChecksCatalogByGroup").Return(
		nil, fmt.Errorf("Error during GetChecksCatalogByGroup"),
	)

	var err error
	config := setupTestConfig()
	app, err := NewAppWithDeps(config, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/catalog", nil)
	req.Header.Set("Accept", "text/html")

	app.webEngine.ServeHTTP(resp, req)

	responseBody := minifyHtml(resp.Body.String())

	assert.Equal(t, 500, resp.Code)
	assert.NoError(t, err)
	assert.Contains(t, responseBody, "<h1>Ooops</h1>")

	tipMsg := "Checks catalog couldn't be retrieved"
	assert.Regexp(t, regexp.MustCompile("Error during GetChecksCatalogByGroup</br>"), responseBody)
	assert.Regexp(t, regexp.MustCompile(fmt.Sprintf("%s</br>", tipMsg)), responseBody)

	checksService.AssertExpectations(t)

}
