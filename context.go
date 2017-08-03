package nile

import "net/http"

// Context represents the information needed interpret and interact with the
// contents of an HTTP request.
type Context struct {
	params  map[string]string
	Request *http.Request
}

// Param gets the value of a URL parameter based on its name. It returns a
// tuple with the value of the parameter as a string and a bool indicating
// whether that parameter exists in the URL.
func (c *Context) Param(name string) (string, bool) {
	param, exists := c.params[name]
	return param, exists
}

func (c *Context) addParam(name, value string) {
	c.params[name] = value
}
