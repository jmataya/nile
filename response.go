package nile

import "net/http"

// Response is the representation of an HTTP response.
type Response interface {
	Status() string
	StatusCode() int
}

type errorResponse struct {
	Error      string `json:"error"`
	statusCode int
}

func (e errorResponse) StatusCode() int {
	return e.statusCode
}

func (e errorResponse) Status() string {
	return http.StatusText(e.statusCode)
}

var resourceNotFound = errorResponse{
	Error:      "Resource not found",
	statusCode: http.StatusNotFound,
}

var methodNotAllowed = errorResponse{
	Error:      "Method not allowed",
	statusCode: http.StatusMethodNotAllowed,
}

var internalServiceError = errorResponse{
	Error:      "Internal service error",
	statusCode: http.StatusInternalServerError,
}
