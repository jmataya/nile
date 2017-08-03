package nile

import "net/http"

// Context represents the information needed to interpret and interact with the
// contents of an HTTP request within a specific HTTP endpoint or middleware
// handler function.
//
// It internally manages error states so that controller and handler code is
// more concise and easy to read. This means that the Context is _only_ valid
// during the execution of a single handler and is discarded after that.
type Context interface {
	// Param gets the value of a URL parameter based on its name. It returns a
	// tuple with the value of the parameter as a string and a bool indicating
	// whether that parameter exists in the URL.
	Param(name string) (string, bool)

	// Request gets the reference to the original HTTP request made by the client.
	Request() *http.Request
}

type context struct {
	params  map[string]string
	request *http.Request
}

func (c context) Request() *http.Request {
	return c.request
}

func (c context) Param(name string) (string, bool) {
	param, exists := c.params[name]
	return param, exists
}

func (c *context) addParam(name, value string) {
	c.params[name] = value
}

func (c *context) setRequest(req *http.Request) {
	c.request = req
}
