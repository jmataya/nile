package nile

import "net/http"

// ErrorResponse is an opinionated structure for how errors should be
// represented in an API. At their bare minimum, they should contain
// a status, human-readable error message, code, and a link to documentation
// (if documentation exists).
//
// ErrorResponse conforms to the Response interface and can be embedded in other
// response constructs, if desired.
//
// In addition, ErrorResponse implements the error interface, so it can be
// returned like an error in services.
type ErrorResponse struct {
	Status          int
	Code            string
	Message         string
	InternalMessage string
	MoreInfo        string
}

// Error writes the contents of the response as a string.
func (er ErrorResponse) Error() string {
	return er.InternalMessage
}

// Body response with the contents of the Response that should be returned in an
// HTTP response.
func (er ErrorResponse) Body() interface{} {
	body := map[string]interface{}{
		"status":  er.Status,
		"code":    er.Code,
		"message": er.Message,
	}

	if er.MoreInfo != "" {
		body["more_info"] = er.MoreInfo
	}

	return body
}

// StatusCode gives the status code that should be used in an HTTP response.
func (er ErrorResponse) StatusCode() int {
	return er.Status
}

// NewInternalServiceError returns an error response that can be used when an
// unexpected error occurs.
func NewInternalServiceError(err error) *ErrorResponse {
	return &ErrorResponse{
		Status:          http.StatusInternalServerError,
		Code:            "00001",
		Message:         "An unknown error occurred",
		InternalMessage: err.Error(),
	}
}

// NewResourceNotFound returns an error when a 404 occurs because a route is not
// found.
func NewResourceNotFound() *ErrorResponse {
	const msg = "Requested resource is not found"

	return &ErrorResponse{
		Status:          http.StatusNotFound,
		Code:            "00002",
		Message:         msg,
		InternalMessage: msg,
	}
}

// NewMethodNotAllowed returns an error that occurs when a route matches an
// HTTP request path, but does not have a matching HTTP method.
func NewMethodNotAllowed() *ErrorResponse {
	const msg = "Method not allowed"

	return &ErrorResponse{
		Status:          http.StatusMethodNotAllowed,
		Code:            "00003",
		Message:         msg,
		InternalMessage: msg,
	}
}

// NewBadRequest returns an error when a 400 Bad Request should occur. The
// general guidance is to use this error when a request is malformed for some
// reason.
func NewBadRequest(code string, err error) *ErrorResponse {
	return &ErrorResponse{
		Status:          http.StatusBadRequest,
		Code:            code,
		Message:         err.Error(),
		InternalMessage: err.Error(),
	}
}

// NewJSONMalformedError returns an error that occurs when parsing a JSON
// payload fails.
func NewJSONMalformedError(err error) *ErrorResponse {
	return NewBadRequest("00004", err)
}

// NewNotFoundError returns an error that is appropriate to use when an entity
// is not found during the processing of a request and you want to signify the
// result using a 404.
func NewNotFoundError(code string, err error) *ErrorResponse {
	return &ErrorResponse{
		Status:          http.StatusNotFound,
		Code:            code,
		Message:         err.Error(),
		InternalMessage: err.Error(),
	}
}
