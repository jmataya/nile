package nile

// Match is a structure that contains the data for when a request matches an
// endpoint.
type Match struct {
	Endpoint      Endpoint
	Context       *Context
	RequestURI    string
	RequestMethod string
}

// NewMatch creates a new Match object.
func NewMatch(endPt Endpoint, path string, method string) *Match {
	return &Match{
		Endpoint:      endPt,
		Context:       &Context{params: map[string]string{}},
		RequestURI:    path,
		RequestMethod: method,
	}
}

// AddParam adds a parameter value to the Match.
func (m *Match) AddParam(key, value string) {
	m.Context.addParam(key, value)
}

// Param gets the value of a param, if it exists.
func (m *Match) Param(key string) (string, bool) {
	return m.Context.Param(key)
}
