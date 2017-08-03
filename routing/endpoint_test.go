package routing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jmataya/nile/utils"
)

func TestEndpointMethods(t *testing.T) {
	var tests = []struct {
		method string
		want   error
	}{
		{http.MethodConnect, fmt.Errorf(errUnsupportedMethod, http.MethodConnect)},
		{http.MethodDelete, nil},
		{http.MethodGet, nil},
		{http.MethodHead, fmt.Errorf(errUnsupportedMethod, http.MethodHead)},
		{http.MethodOptions, fmt.Errorf(errUnsupportedMethod, http.MethodOptions)},
		{http.MethodPatch, nil},
		{http.MethodPost, nil},
		{http.MethodPut, nil},
		{http.MethodTrace, fmt.Errorf(errUnsupportedMethod, http.MethodTrace)},
		{"get", fmt.Errorf(errInvalidMethod, "get")},
		{"", fmt.Errorf(errInvalidMethod, "")},
	}

	for _, test := range tests {
		method := test.method

		dummyHandler := func(Context) Response {
			return internalServiceError
		}

		endpoint, err := NewEndpoint(method, dummyHandler)
		if !utils.CheckError(test.want, err) {
			t.Errorf("NewEndpoint(%s) error, want %v, got %v", method, test.want, err)
			continue
		}

		if err != nil {
			continue
		}

		if endpoint.Method() != method {
			t.Errorf("Endpoint.Method(), expected %s, got %s", method, endpoint.Method())
		}
	}
}
