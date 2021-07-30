package web

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/web/models"
	"github.com/trento-project/trento/web/services/mocks"
)

func TestChecksCatalogHandler(t *testing.T) {
	checksMocks := new(mocks.ChecksService)

	deps := DefaultDependencies()
	deps.checksService = checksMocks

	checks := map[string]map[string]*models.Check{
		"group 1": {
			"1.1.1": &models.Check{
				ID:             "1.1.1",
				Name:           "check 1",
				Group:          "group 1",
				Description:    "description 1",
				Remediation:    "remediation 1",
				Implementation: "implementation 1",
				Labels:         "labels 1",
			},
			"1.1.2": &models.Check{
				ID:             "1.1.2",
				Name:           "check 2",
				Group:          "group 1",
				Description:    "description 2",
				Remediation:    "remediation 2",
				Implementation: "implementation 2",
				Labels:         "labels 2",
			},
		},
		"group 2": {
			"1.2.3": &models.Check{
				ID:             "1.2.3",
				Name:           "check 3",
				Group:          "group 2",
				Description:    "description 3",
				Remediation:    "remediation 3",
				Implementation: "implementation 3",
				Labels:         "labels 3",
			},
		},
	}

	checksMocks.On("GetChecksCatalogByGroup").Return(
		checks, nil,
	)

	var err error
	app, err := NewAppWithDeps("", 80, deps)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/catalog", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ServeHTTP(resp, req)

	responseBody := minifyHtml(resp.Body.String())

	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, responseBody, "Checks catalog")

	assert.Regexp(t, regexp.MustCompile("<h4.*>group 1</h4>.*<td.*>1.1.1</td><td.*>description 1<div.*id=info-1-1-1.*><p>remediation 1</p></div><div.*implementation 1.*</div>.*</td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<h4.*>group 1</h4>.*<td.*>1.1.2</td><td.*>description 2<div.*id=info-1-1-2.*><p>remediation 2</p></div><div.*implementation 2.*</div>.*</td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<h4.*>group 2</h4>.*<td.*>1.2.3</td><td.*>description 3<div.*id=info-1-2-3.*><p>remediation 3</p></div><div.*implementation 3.*</div>.*</td>"), responseBody)
	assert.Equal(t, 2, strings.Count(responseBody, "<h4"))
	assert.Equal(t, 5, strings.Count(responseBody, "<tr>"))

	checksMocks.AssertExpectations(t)
}
