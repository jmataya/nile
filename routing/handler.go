package routing

// HandlerFunc is the method signature for accepting an HTTP request and
// delivering a response.
type HandlerFunc func(*Context) Response
