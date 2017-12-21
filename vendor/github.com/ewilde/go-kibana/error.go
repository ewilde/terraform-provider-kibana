package kibana

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
)

// Error represents an error response from the PagerDuty API.
type HttpError struct {
	ErrorResponse gorequest.Response
	Code          int
	Message       string
	Body          string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%s API call to %s failed %v. Code: %d, Body: %s, Message: %s",
		e.ErrorResponse.Request.Method, e.ErrorResponse.Request.URL.String(), e.ErrorResponse.Status, e.Code, e.Body, e.Message)
}

func NewError(response gorequest.Response, body string, message string) *HttpError {
	return &HttpError{
		Code:          response.StatusCode,
		ErrorResponse: response,
		Body:          body,
		Message:       message,
	}
}
