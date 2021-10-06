package web

import (
	"fmt"
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

	deps := testDependencies()
	deps.checksService = checksMocks

	checks := models.GroupedCheckList{
		&models.GroupedChecks{
			Group: "group 1",
			Checks: models.CheckList{
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
				},
			},
		},
		&models.GroupedChecks{
			Group: "group 2",
			Checks: models.CheckList{
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

	assert.Regexp(t, regexp.MustCompile("<h4.*>group 1</h4>.*<td.*>ABCDEF</td><td.*>description 1<div.*id=info-ABCDEF.*><p>remediation 1</p></div><div.*implementation 1.*</div>.*</td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<h4.*>group 1</h4>.*<td.*>123456</td><td.*>description 2<div.*id=info-123456.*><p>remediation 2</p></div><div.*implementation 2.*</div>.*</td>"), responseBody)
	assert.Regexp(t, regexp.MustCompile("<h4.*>group 2</h4>.*<td.*>123ABC</td><td.*>description 3<div.*id=info-123ABC.*><p>remediation 3</p></div><div.*implementation 3.*</div>.*</td>"), responseBody)
	assert.Equal(t, 2, strings.Count(responseBody, "<h4"))
	assert.Equal(t, 5, strings.Count(responseBody, "<tr>"))

	checksMocks.AssertExpectations(t)
}

func TestChecksCatalogHandlerError(t *testing.T) {

	checksMocks := new(mocks.ChecksService)

	deps := testDependencies()
	deps.checksService = checksMocks

	checksMocks.On("GetChecksCatalogByGroup").Return(
		nil, fmt.Errorf("Error during GetChecksCatalogByGroup"),
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
	req.Header.Set("Accept", "text/html")

	app.ServeHTTP(resp, req)

	responseBody := minifyHtml(resp.Body.String())

	assert.Equal(t, 500, resp.Code)
	assert.NoError(t, err)
	assert.Contains(t, responseBody, "<h1>Ooops</h1>")

	tipMsg := "Checks catalog couldn't be retrieved. Check if the ARA service is running" +
		" and the --ara-addr flag is pointing corretly to the service"
	assert.Regexp(t, regexp.MustCompile("Error during GetChecksCatalogByGroup</br>"), responseBody)
	assert.Regexp(t, regexp.MustCompile(fmt.Sprintf("%s</br>", tipMsg)), responseBody)

	checksMocks.AssertExpectations(t)

}
