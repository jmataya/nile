package nile

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Context represents the information needed to interpret and interact with the
// contents of an HTTP request within a specific HTTP endpoint or middleware
// handler function.
//
// It internally manages error states so that controller and handler code is
// more concise and easy to read. This means that the Context is _only_ valid
// during the execution of a single handler and is discarded after that.
type Context interface {
	// BindJSON attempts to unmarshal and validate JSON  data from an HTTP request
	// body into a Payload object. Failure to either unmarshal or validate will
	// result the Context's internal error state being set and defaulting in an
	// HTTP Bad Request Error. That error is returned here for convenience.
	BindJSON(payload Payload) error

	// Error returns any error that may be associated with the Context.
	Error() error

	// Param gets the value of a URL parameter. If the value is not found, an
	// error is set on the context.
	Param(name string) string

	// TryParam gets the value of a URL parameter based on its name. It returns a
	// tuple with the value of the parameter as a string and a bool indicating
	// whether that parameter exists in the URL.
	TryParam(name string) (string, bool)

	// Request gets the reference to the original HTTP request made by the client.
	Request() *http.Request
}

type context struct {
	err     *ErrorResponse
	params  map[string]string
	request *http.Request
}

func (c *context) BindJSON(payload Payload) error {
	if c.err != nil {
		return nil
	}

	decoder := json.NewDecoder(c.request.Body)
	if err := decoder.Decode(payload); err != nil {
		return c.setError(NewJSONMalformedError(err))
	}

	if err := payload.Validate(); err != nil {
		return c.setError(err)
	}

	return nil
}

func (c *context) Error() error {
	return c.err
}

func (c *context) setError(er *ErrorResponse) error {
	c.err = er
	return er
}

func (c context) Request() *http.Request {
	return c.request
}

func (c *context) Param(name string) string {
	param, exists := c.params[name]
	if !exists {
		err := fmt.Errorf("Unable to parse parameter %s", name)
		c.setError(NewInternalServiceError(err))
		return ""
	}
	return param
}

func (c context) TryParam(name string) (string, bool) {
	param, exists := c.params[name]
	return param, exists
}

func (c *context) addParam(name, value string) {
	c.params[name] = value
}

func (c *context) setRequest(req *http.Request) {
	c.request = req
}
