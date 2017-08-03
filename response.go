package nile

import (
	"bytes"
	"net/http"
)

// Response is the representation of an HTTP response.
type Response interface {
	Status() string
	StatusCode() int
}

// ErrorResponse is the representation of an HTTP error response.
type ErrorResponse interface {
	Error() string
	Status() string
	StatusCode() int
}

// NewBadRequest creates an ErrorResponse for when an HTTP Bad Request occurs.
func NewBadRequest(err error) ErrorResponse {
	return makeErrorResponse(err.Error(), http.StatusBadRequest)
}

// NewInternalServiceError creates an ErrorResponse for when an HTTP Internal
// Service Error occurs.
func NewInternalServiceError(err error) ErrorResponse {
	return makeErrorResponse(err.Error(), http.StatusInternalServerError)
}

type errorResponse struct {
	Errors     []string `json:"errors"`
	statusCode int
}

func makeErrorResponse(err string, statusCode int) errorResponse {
	return errorResponse{
		Errors:     []string{err},
		statusCode: statusCode,
	}
}

func (e errorResponse) Error() string {
	var buffer bytes.Buffer

	for idx, errString := range e.Errors {
		if idx > 0 {
			buffer.WriteString(", ")
		}

		buffer.WriteString(errString)
	}

	return buffer.String()
}

func (e errorResponse) StatusCode() int {
	return e.statusCode
}

func (e errorResponse) Status() string {
	return http.StatusText(e.statusCode)
}

var resourceNotFound = makeErrorResponse("Resource not found", http.StatusNotFound)
var methodNotAllowed = makeErrorResponse("Method not allowed", http.StatusMethodNotAllowed)
var internalServiceError = makeErrorResponse("Internal service error", http.StatusInternalServerError)
