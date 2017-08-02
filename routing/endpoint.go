package routing

import (
	"fmt"
	"net/http"
)

const (
	errInvalidMethod     = "Invalid HTTP method for endpoint %s"
	errUnsupportedMethod = "HTTP method %s is not currently supported as an HTTP endpoint"
)

// Endpoint the leaf node on Segment tree that corresponds to an actual
// endpoint that can be hit via an HTTP request. It contains an HTTP method
// and handler for when the request matches.
type Endpoint interface {
	// Method gets the HTTP method to which Endpoint will respond.
	Method() string
}

// NewEndpoint creates a new, valid Endpoint based on an HTTP method.
func NewEndpoint(method string) (Endpoint, error) {
	// Validate that method is a currently supported HTTP method.
	isSupported, ok := supportedMethods[method]
	if !ok {
		return nil, fmt.Errorf(errInvalidMethod, method)
	} else if !isSupported {
		return nil, fmt.Errorf(errUnsupportedMethod, method)
	}

	return &httpEndpoint{method: method}, nil
}

var supportedMethods = map[string]bool{
	http.MethodConnect: false,
	http.MethodDelete:  true,
	http.MethodGet:     true,
	http.MethodHead:    false,
	http.MethodOptions: false,
	http.MethodPatch:   true,
	http.MethodPost:    true,
	http.MethodPut:     true,
	http.MethodTrace:   false,
}

type httpEndpoint struct {
	method string
}

func (h *httpEndpoint) Method() string {
	return h.method
}
