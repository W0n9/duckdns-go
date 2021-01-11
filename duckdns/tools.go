package duckdns

import (
	"fmt"
	"net/http"
)

// An ErrorResponse represents an API response that generated an error.
type ErrorResponse struct {
	HTTPResponse *http.Response
}

// Error implements the error interface.
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %v", r.HTTPResponse.Request.Method, r.HTTPResponse.Request.URL, r.HTTPResponse.StatusCode)
}

//CheckResponse function
func CheckResponse(resp *http.Response) error {

	statusCode := resp.StatusCode
	if statusCode == 200 {
		return nil
	}

	errorResponse := &ErrorResponse{}
	errorResponse.HTTPResponse = resp

	return errorResponse
}
