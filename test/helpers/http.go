package helpers

import "net/http"

// RoundTripFunc Needed to Mock Http client requests and responses
//
// Usage
// &http.Client{
// 	 Transport: helpers.RoundTripFunc(func(req *http.Request) *http.Response {
// 	    bodyBytes, _ := ioutil.ReadAll(req.Body)
//
// 	 	suite.EqualValues(expectedBody, string(bodyBytes))
//
// 	 	suite.Equal(req.URL.String(), expectedUrl)
// 	 	return &http.Response{
// 			StatusCode: returnedStatusCode,
// 	 	}
//     }),
// }
type RoundTripFunc func(req *http.Request) *http.Response

// ErroringRoundTripFunc Needed to Mock Http client that returns an error even before the request is made
//
// Usage
// &http.Client{
// 	 Transport: helpers.ErroringRoundTripFunc(func() error {
// 	 	return fmt.Errorf("some error")
// 	 }),
// }
type ErroringRoundTripFunc func() error

// RoundTripFunc implements the RoundTripper interface
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// ErroringRoundTripFunc implements the RoundTripper interface
func (f ErroringRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, f()
}
