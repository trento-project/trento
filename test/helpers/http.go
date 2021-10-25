package helpers

import "net/http"

// Usage

// &http.Client{Transport: helpers.RoundTripFunc(func(req *http.Request) *http.Response {
// 	 bodyBytes, _ := ioutil.ReadAll(req.Body)
//
// 	 suite.EqualValues(expectedBody, string(bodyBytes))
//
// 	 suite.Equal(req.URL.String(), expectedUrl)
// 	 return &http.Response{
// 		StatusCode: returnedStatusCode,
// 	 }
// }}

// RoundTripFunc Needed to Mock Http client requests and responses
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip implements the RoundTripper interface
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}
