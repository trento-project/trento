package web

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_homeHandler(t *testing.T) {

	templs := NewTemplateRender(templatesFS, "templates/*.tmpl")

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(IndexHandler(templs.templates))
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// test body with console expected text
	expected := "SUSE Console for SAP Applications"
	if b, err := ioutil.ReadAll(rr.Body); err != nil {
		t.Fail()
	} else {
		if strings.Contains(string(b), "Error") {
			t.Errorf("header response shouldn't return error: %s", b)
		} else if !strings.Contains(string(b), expected) {
			t.Errorf("header response doesn't match:\n%s", b)
		}
	}

}
