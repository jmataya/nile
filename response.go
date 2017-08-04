package nile

// Response is the representation of an HTTP response.
type Response interface {
	// Body response with the contents of the Response that should be returned in
	// an HTTP response.
	Body() interface{}

	// StatusCode gives the status code that should be used in an HTTP response.
	StatusCode() int
}

// GenericResponse is the structure for a simple HTTP response.
type GenericResponse struct {
	body   interface{}
	status int
}

// NewGenericResponse creates a new GenericResponse object.
func NewGenericResponse(status int, body interface{}) *GenericResponse {
	return &GenericResponse{
		body:   body,
		status: status,
	}
}

// Body response with the contents of the Response that should be returned in
// an HTTP response.
func (gr GenericResponse) Body() interface{} {
	return gr.body
}

// StatusCode gives the status code that should be used in an HTTP response.
func (gr GenericResponse) StatusCode() int {
	return gr.status
}
